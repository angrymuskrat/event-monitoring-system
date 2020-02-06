package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/angrymuskrat/event-monitoring-system/services/coordinator/service/status"
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler"
	crawlerdata "github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/data"
	crservice "github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/service"
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	"github.com/google/uuid"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Session struct {
	ID        string
	Params    SessionParameters
	Status    status.Status
	Endpoints SessionEndpoints
}

type SessionParameters struct {
	CityID          string
	CityName        string
	TopLeft         data.Point
	BottomRight     data.Point
	Locations       []string
	FinishTimestamp int64
}

type SessionEndpoints struct {
	Crawler string
}

func NewSession(p SessionParameters, e SessionEndpoints) Session {
	id := uuid.New().String()
	return Session{
		ID:        id,
		Params:    p,
		Endpoints: e,
		Status:    status.HistoricBuilding{},
	}
}

func (s *Session) Run() {

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
		FinishTimestamp: s.Params.FinishTimestamp,
	}
	d, err := json.Marshal(p)
	if err != nil {
		unilog.Logger().Error("unable to marshal crawler parameters", zap.Error(err))
		return err
	}
	url := fmt.Sprintf("%s/new", s.Endpoints.Crawler)
	buf := bytes.NewBuffer(d)
	resp, err := http.Post(url, "application/json", buf)
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
	url := fmt.Sprintf("%s/status/%s", s.Endpoints.Crawler, cs.SessionID)
	resp, err := http.Get(url)
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
	url := fmt.Sprintf("%s/stop", s.Endpoints.Crawler)
	resp, err := http.Post(url, "application/json", buf)
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

}

func (s *Session) monitoring() error {

}
