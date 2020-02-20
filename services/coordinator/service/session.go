package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/angrymuskrat/event-monitoring-system/services/coordinator/service/status"
	"github.com/angrymuskrat/event-monitoring-system/services/event-detection/proto"
	"github.com/angrymuskrat/event-monitoring-system/services/event-detection/service"
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler"
	crawlerdata "github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/data"
	crservice "github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/service"
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	"github.com/google/uuid"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Session struct {
	ID        string
	Params    SessionParameters
	Status    status.Status
	Endpoints ServiceEndpoints
	edClient  service.Client
}

type SessionParameters struct {
	CityID         string
	CityName       string
	Timezone       string
	TopLeft        data.Point
	BottomRight    data.Point
	Locations      []string
	CrawlerFinish  int64
	HistoricStart  int64
	HistoricFinish int64
	GridSize       float64
	AuthCookie     string
	SkipCrawling   bool
	SkipHistoric   bool
}

func NewSession(p SessionParameters, e ServiceEndpoints) (*Session, error) {
	id := uuid.New().String()
	conn, err := grpc.Dial(e.EventDetection.Address, grpc.WithInsecure())
	if err != nil {
		unilog.Logger().Error("unable to connect to data strorage", zap.Error(err))
		return nil, err
	}
	client := service.NewClient(conn)
	return &Session{
		ID:        id,
		Params:    p,
		Endpoints: e,
		Status:    status.HistoricBuilding{SessionID: id, Status: "session is starting"},
		edClient:  client,
	}, nil
}

func (s *Session) Run() {
	var err error
	if !s.Params.SkipCrawling {
		area := data.Area{TopLeft: &s.Params.TopLeft, BotRight: &s.Params.BottomRight}
		city := data.City{Title: s.Params.CityName, Code: s.Params.CityID, Area: area}
		err = Storage.InsertCity(context.Background(), city, true)
		if err != nil {
			unilog.Logger().Error("unable to insert city", zap.Any("city", city), zap.Error(err))
			return
		}
		err = s.historicCollect()
		if err != nil {
			return
		}
	}
	if !s.Params.SkipHistoric {
		err = s.historicBuild()
		if err != nil {
			return
		}
	}
	s.monitoring()
}

func (s *Session) historicCollect() error {
	err := s.startCollect()
	if err != nil {
		st := status.Failed{Error: err}
		s.Status = st
		return err
	}
	run := true
	for run {
		num, err := s.checkCollect()
		if err != nil {
			st := status.Failed{Error: err}
			s.Status = st
			return err
		}
		run = num > 0
		time.Sleep(5 * time.Minute)
	}
	err = s.deleteCollect()
	if err != nil {
		st := status.Failed{Error: err}
		s.Status = st
		return err
	}
	return nil
}

func (s *Session) startCollect() error {
	p := crawler.Parameters{
		CityID:          s.Params.CityID,
		Type:            crawlerdata.LocationsType,
		Entities:        s.Params.Locations,
		FinishTimestamp: s.Params.CrawlerFinish,
		AuthCookie:      s.Params.AuthCookie,
	}
	d, err := json.Marshal(p)
	if err != nil {
		unilog.Logger().Error("unable to marshal crawler parameters", zap.Error(err))
		return err
	}
	url := fmt.Sprintf("http://%s/new", s.Endpoints.Crawler)
	buf := bytes.NewBuffer(d)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		unilog.Logger().Error("unable to make request to crawler", zap.Error(err))
		return err
	}
	req.SetBasicAuth(s.Endpoints.Crawler.User, s.Endpoints.Crawler.Password)
	resp, err := client.Do(req)
	if err != nil {
		unilog.Logger().Error("unable to make request to crawler", zap.Error(err))
		return err
	}
	if resp.StatusCode != 200 {
		unilog.Logger().Error("error status code", zap.Int("code", resp.StatusCode), zap.String("status", resp.Status))
		return err
	}
	defer resp.Body.Close()
	var ep crservice.NewEpResponse
	err = json.NewDecoder(resp.Body).Decode(&ep)
	if err != nil {
		unilog.Logger().Error("unable to read response", zap.Error(err))
		return nil
	}
	if ep.Error != "" {
		unilog.Logger().Error("error in crawler", zap.String("msg", ep.Error))
		return errors.New(ep.Error)
	}
	st := status.HistoricCollection{
		SessionID:      ep.ID,
		PostsCollected: 0,
		LocationsLeft:  len(p.Entities),
	}
	s.Status = st
	return nil
}

func (s *Session) checkCollect() (int, error) {
	cs, ok := s.Status.(status.HistoricCollection)
	if !ok {
		unilog.Logger().Error("incorrect session status", zap.String("status", s.Status.String()))
	}
	url := fmt.Sprintf("http://%s/status/%s", s.Endpoints.Crawler, cs.SessionID)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		unilog.Logger().Error("unable to make request to crawler", zap.Error(err))
		return -1, err
	}
	req.SetBasicAuth(s.Endpoints.Crawler.User, s.Endpoints.Crawler.Password)
	resp, err := client.Do(req)
	//resp, err := http.Get(url)
	if err != nil {
		unilog.Logger().Error("unable to make request to crawler", zap.Error(err))
		return -1, err
	}
	if resp.StatusCode != 200 {
		unilog.Logger().Error("error status code", zap.Int("code", resp.StatusCode), zap.String("status", resp.Status))
		return -1, err
	}
	defer resp.Body.Close()
	var ep crservice.StatusEpResponse
	err = json.NewDecoder(resp.Body).Decode(&ep)
	if err != nil {
		unilog.Logger().Error("unable to read response", zap.Error(err))
		return -1, nil
	}
	if ep.Error != "" {
		unilog.Logger().Error("error in crawler", zap.String("msg", ep.Error))
		return -1, errors.New(ep.Error)
	}
	st := status.HistoricCollection{
		SessionID:      cs.SessionID,
		PostsCollected: ep.Status.PostsCollected,
		LocationsLeft:  ep.Status.EntitiesLeft,
	}
	s.Status = st
	return ep.Status.EntitiesLeft, nil
}

func (s *Session) deleteCollect() error {
	cs, ok := s.Status.(status.HistoricCollection)
	if !ok {
		unilog.Logger().Error("incorrect session status", zap.String("status", s.Status.String()))
	}
	p := crservice.IDEpRequest{ID: cs.SessionID}
	d, err := json.Marshal(p)
	if err != nil {
		unilog.Logger().Error("unable to marshal crawler parameters", zap.Error(err))
		return err
	}
	buf := bytes.NewBuffer(d)
	url := fmt.Sprintf("http://%s/stop", s.Endpoints.Crawler)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		unilog.Logger().Error("unable to make request to crawler", zap.Error(err))
		return err
	}
	req.SetBasicAuth(s.Endpoints.Crawler.User, s.Endpoints.Crawler.Password)
	resp, err := client.Do(req)
	//resp, err := http.Post(url, "application/json", buf)
	if err != nil {
		unilog.Logger().Error("unable to make request to crawler", zap.Error(err))
		return err
	}
	if resp.StatusCode != 200 {
		unilog.Logger().Error("error status code", zap.Int("code", resp.StatusCode), zap.String("status", resp.Status))
		return err
	}
	defer resp.Body.Close()
	var ep crservice.StopEpResponse
	err = json.NewDecoder(resp.Body).Decode(&ep)
	if err != nil {
		unilog.Logger().Error("unable to read response", zap.Error(err))
		return nil
	}
	if ep.Error != "" {
		unilog.Logger().Error("error in crawler", zap.String("msg", ep.Error))
		return errors.New(ep.Error)
	}
	if !ep.Ok {
		msg := "unable to delete session in crawler"
		unilog.Logger().Error(msg)
		return errors.New(msg)
	}
	return nil
}

func (s *Session) historicBuild() error {
	err := s.historicStart()
	if err != nil {
		s.Status = status.Failed{Error: err}
		return err
	}
	finished := false
	for !finished {
		finished, err = s.historicStatus()
		if err != nil {
			s.Status = status.Failed{Error: err}
			return err
		}
		time.Sleep(1 * time.Minute)
	}
	return nil
}

func (s *Session) historicStart() error {
	req := proto.HistoricRequest{
		Timezone:   s.Params.Timezone,
		CityId:     s.Params.CityID,
		StartTime:  s.Params.HistoricStart,
		FinishTime: s.Params.HistoricFinish,
		Area:       &data.Area{TopLeft: &s.Params.TopLeft, BotRight: &s.Params.BottomRight},
		GridSize:   s.Params.GridSize,
	}
	unilog.Logger().Info("historicStart was started")
	respRaw, err := s.edClient.HistoricGrids(context.Background(), req)
	if err != nil {
		unilog.Logger().Error("error during historic building initiation", zap.Error(err))
		return err
	}
	resp := respRaw.(proto.HistoricResponse)
	if resp.Err != "" {
		err := errors.New(resp.Err)
		unilog.Logger().Error("server error", zap.Error(err))
		return err
	}
	s.Status = status.HistoricBuilding{
		SessionID: resp.Id,
	}
	return nil
}

func (s *Session) historicStatus() (bool, error) {
	st := s.Status.Get().(status.HistoricBuilding)
	req := proto.StatusRequest{
		Id: st.SessionID,
	}
	respRaw, err := s.edClient.HistoricStatus(context.Background(), req)
	if err != nil {
		unilog.Logger().Error("error during historic status checking", zap.Error(err))
		return false, err
	}
	resp := respRaw.(proto.StatusResponse)
	if resp.Err != "" {
		err := errors.New(resp.Err)
		unilog.Logger().Error("server error", zap.Error(err))
		return false, err
	}
	s.Status = status.HistoricBuilding{
		SessionID: st.SessionID,
		Status:    resp.Status,
	}
	return resp.Finished, nil
}

func (s *Session) monitoring() error {
	start, finish := s.Params.HistoricFinish, time.Now().Unix()
	var err error
	for {
		err = s.eventsStart(start, finish)
		if err != nil {
			s.Status = status.Failed{Error: err}
			return err
		}
		finished := false
		for !finished {
			finished, err = s.eventsStatus()
			if err != nil {
				s.Status = status.Failed{Error: err}
				return err
			}
			time.Sleep(20 * time.Second)
		}
		time.Sleep(5 * time.Minute)
		st := s.Status.Get().(status.Monitoring)
		start, finish = st.CurrentTimestamp, time.Now().Unix()
	}
}

func (s *Session) eventsStart(start, finish int64) error {
	req := proto.EventRequest{
		Timezone:   s.Params.Timezone,
		CityId:     s.Params.CityID,
		StartTime:  start,
		FinishDate: finish,
	}
	respRaw, err := s.edClient.FindEvents(context.Background(), req)
	if err != nil {
		unilog.Logger().Error("error during events search initiation", zap.Error(err))
		return err
	}
	resp := respRaw.(proto.EventResponse)
	if resp.Err != "" {
		err := errors.New(resp.Err)
		unilog.Logger().Error("server error", zap.Error(err))
		return err
	}
	s.Status = status.Monitoring{
		SessionID:        resp.Id,
		CurrentTimestamp: finish,
	}
	return nil
}

func (s *Session) eventsStatus() (bool, error) {
	st := s.Status.Get().(status.Monitoring)
	req := proto.StatusRequest{
		Id: st.SessionID,
	}
	respRaw, err := s.edClient.EventsStatus(context.Background(), req)
	if err != nil {
		unilog.Logger().Error("error during events status checking", zap.Error(err))
		return false, err
	}
	resp := respRaw.(proto.StatusResponse)
	if resp.Err != "" {
		err := errors.New(resp.Err)
		unilog.Logger().Error("server error", zap.Error(err))
		return false, err
	}
	s.Status = status.Monitoring{
		SessionID:        st.SessionID,
		CurrentTimestamp: st.CurrentTimestamp,
		Status:           resp.Status,
	}
	return resp.Finished, nil
}
