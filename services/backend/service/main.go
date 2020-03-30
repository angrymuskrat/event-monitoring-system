package service

import (
	"encoding/json"
	"errors"
	"github.com/angrymuskrat/event-monitoring-system/services/proto"
	"github.com/gorilla/mux"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

var svc *backendService

func Start(confPath string) {
	conf, err := readConfig(confPath)
	if err != nil {
		return
	}
	conn, err := setConnector(conf.Connector, conf.ConnectorParams)
	if err != nil {
		unilog.Logger().Error("unable to create storage connector", zap.Error(err))
		return
	}
	svc = &backendService{
		storageConn: conn,
	}
	sm, err := newAuthManager(conf.SessionKey, conf.AuthLogPath)
	if err != nil {
		panic(err)
	}
	tm, err := newTimerManager(conf.TimerLogPath)
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/heatmap/{city}/{topLeft}/{botRight}/{hour}", heatmap).Methods("GET")
	r.HandleFunc("/timeline/{city}/{start}/{finish}", timeline).Methods("GET")
	r.HandleFunc("/events/{city}/{topLeft}/{botRight}/{hour}", events).Methods("GET")
	r.HandleFunc("/search/{city}/{tags}/{start}/{finish}", search).Methods("GET")
	r.HandleFunc("/login", sm.login).Methods("POST")
	r.Use(sm.Handler)
	r.Use(tm.Handler)

	http.Handle("/", accessControl(r))
	err = http.ListenAndServe(conf.Address, nil)
	if err != nil {
		unilog.Logger().Error("error in service handler", zap.Error(err))
	}
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Access-Token, Uid")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

func heatmap(w http.ResponseWriter, r *http.Request) {
	req, err := decodeHeatmapRequest(r)
	if err != nil {
		unilog.Logger().Error("unable to decode request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	d, err := svc.HeatmapPosts(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(d)
	if err != nil {
		unilog.Logger().Error("unable to encode result to JSON", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func timeline(w http.ResponseWriter, r *http.Request) {
	req, err := decodeTimelineRequest(r)
	if err != nil {
		unilog.Logger().Error("unable to decode request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	d, err := svc.Timeline(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(d)
	if err != nil {
		unilog.Logger().Error("unable to encode result to JSON", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func events(w http.ResponseWriter, r *http.Request) {
	req, err := decodeEventsRequest(r)
	if err != nil {
		unilog.Logger().Error("unable to decode request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	d, err := svc.Events(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(d)
	if err != nil {
		unilog.Logger().Error("unable to encode result to JSON", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func search(w http.ResponseWriter, r *http.Request) {
	req, err := decodeSearchRequest(r)
	if err != nil {
		unilog.Logger().Error("unable to decode request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	d, err := svc.SearchEvents(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(d)
	if err != nil {
		unilog.Logger().Error("unable to encode result to JSON", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func decodeHeatmapRequest(r *http.Request) (HeatmapRequest, error) {
	vars := mux.Vars(r)
	req := HeatmapRequest{}
	city, ok := vars["city"]
	if !ok {
		return HeatmapRequest{}, errors.New("unable to get city name")
	}
	req.City = city
	topLeftRaw, ok := vars["topLeft"]
	if !ok {
		return HeatmapRequest{}, errors.New("unable to get top left coordinates")
	}
	coords := strings.Split(topLeftRaw, ",")
	if len(coords) != 2 {
		return HeatmapRequest{}, errors.New("incorrect format of top left coordinates")
	}
	lat, err := strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return HeatmapRequest{}, errors.New("unable to parse latitude of top left")
	}
	lon, err := strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return HeatmapRequest{}, errors.New("unable to parse longitude of top left")
	}
	req.TopLeft = data.Point{
		Lat: lat,
		Lon: lon,
	}
	botRightRaw, ok := vars["botRight"]
	if !ok {
		return HeatmapRequest{}, errors.New("unable to get bottom right coordinates")
	}
	coords = strings.Split(botRightRaw, ",")
	if len(coords) != 2 {
		return HeatmapRequest{}, errors.New("incorrect format of bottom right coordinates")
	}
	lat, err = strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return HeatmapRequest{}, errors.New("unable to parse latitude of bottom right")
	}
	lon, err = strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return HeatmapRequest{}, errors.New("unable to parse longitude of bottom right")
	}
	req.BottomRight = data.Point{
		Lat: lat,
		Lon: lon,
	}
	hourRaw, ok := vars["hour"]
	if !ok {
		return HeatmapRequest{}, errors.New("unable to get hour")
	}
	hour, err := strconv.ParseInt(hourRaw, 10, 64)
	if err != nil {
		return HeatmapRequest{}, errors.New("incorrect format of hour")
	}
	req.Hour = hour
	return req, nil
}

func decodeTimelineRequest(r *http.Request) (TimelineRequest, error) {
	vars := mux.Vars(r)
	req := TimelineRequest{}
	city, ok := vars["city"]
	if !ok {
		return TimelineRequest{}, errors.New("unable to get city name")
	}
	req.City = city
	startRaw, ok := vars["start"]
	if !ok {
		return TimelineRequest{}, errors.New("unable to get start")
	}
	start, err := strconv.ParseInt(startRaw, 10, 64)
	if err != nil {
		return TimelineRequest{}, errors.New("incorrect format of start")
	}
	req.Start = start
	finishRaw, ok := vars["finish"]
	if !ok {
		return TimelineRequest{}, errors.New("unable to get finish")
	}
	finish, err := strconv.ParseInt(finishRaw, 10, 64)
	if err != nil {
		return TimelineRequest{}, errors.New("incorrect format of finish")
	}
	req.Finish = finish
	return req, nil
}

func decodeEventsRequest(r *http.Request) (EventsRequest, error) {
	vars := mux.Vars(r)
	req := EventsRequest{}
	city, ok := vars["city"]
	if !ok {
		return EventsRequest{}, errors.New("unable to get city name")
	}
	req.City = city
	topLeftRaw, ok := vars["topLeft"]
	if !ok {
		return EventsRequest{}, errors.New("unable to get top left coordinates")
	}
	coords := strings.Split(topLeftRaw, ",")
	if len(coords) != 2 {
		return EventsRequest{}, errors.New("incorrect format of top left coordinates")
	}
	lat, err := strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return EventsRequest{}, errors.New("unable to parse latitude of top left")
	}
	lon, err := strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return EventsRequest{}, errors.New("unable to parse longitude of top left")
	}
	req.TopLeft = data.Point{
		Lat: lat,
		Lon: lon,
	}
	botRightRaw, ok := vars["botRight"]
	if !ok {
		return EventsRequest{}, errors.New("unable to get bottom right coordinates")
	}
	coords = strings.Split(botRightRaw, ",")
	if len(coords) != 2 {
		return EventsRequest{}, errors.New("incorrect format of bottom right coordinates")
	}
	lat, err = strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return EventsRequest{}, errors.New("unable to parse latitude of bottom right")
	}
	lon, err = strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return EventsRequest{}, errors.New("unable to parse longitude of bottom right")
	}
	req.BottomRight = data.Point{
		Lat: lat,
		Lon: lon,
	}
	hourRaw, ok := vars["hour"]
	if !ok {
		return EventsRequest{}, errors.New("unable to get hour")
	}
	hour, err := strconv.ParseInt(hourRaw, 10, 64)
	if err != nil {
		return EventsRequest{}, errors.New("incorrect format of hour")
	}
	req.Hour = hour
	return req, nil
}

func decodeSearchRequest(r *http.Request) (SearchRequest, error) {
	vars := mux.Vars(r)
	req := SearchRequest{}
	city, ok := vars["city"]
	if !ok {
		return SearchRequest{}, errors.New("unable to get city name")
	}
	req.City = city
	startRaw, ok := vars["start"]
	if !ok {
		return SearchRequest{}, errors.New("unable to get start")
	}
	start, err := strconv.ParseInt(startRaw, 10, 64)
	if err != nil {
		return SearchRequest{}, errors.New("incorrect format of start")
	}
	req.Start = start
	finishRaw, ok := vars["finish"]
	if !ok {
		return SearchRequest{}, errors.New("unable to get finish")
	}
	finish, err := strconv.ParseInt(finishRaw, 10, 64)
	if err != nil {
		return SearchRequest{}, errors.New("incorrect format of finish")
	}
	req.Finish = finish
	tagsRaw, ok := vars["tags"]
	tags := strings.Split(tagsRaw, ",")
	req.Keytags = tags
	return req, nil
}
