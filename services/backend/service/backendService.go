package service

import data "github.com/angrymuskrat/event-monitoring-system/services/proto"

type backendService struct {
	storageConn StorageConnector
}

func (s *backendService) HeatmapPosts(req HeatmapRequest) ([]data.AggregatedPost, error) {
	return s.storageConn.HeatmapPosts(req.City, req.TopLeft, req.BottomRight, req.Hour)
}

func (s *backendService) Timeline(req TimelineRequest) (Timeline, error) {
	return s.storageConn.Timeline(req.City, req.Start, req.Finish)
}

func (s *backendService) Events(req EventsRequest) ([]data.Event, error) {
	return s.storageConn.Events(req.City, req.TopLeft, req.BottomRight, req.Hour)
}

func (s *backendService) SearchEvents(req SearchRequest) ([]data.Event, error) {
	return s.storageConn.EventsByTags(req.City, req.Keytags, req.Start, req.Finish)
}

func (s *backendService) ShortPostsInInterval(req ShortPostsRequests) ([]data.ShortPost, error) {
	return s.storageConn.ShortPostsInInterval(req.City, req.Shortcodes, req.Start, req.End)
}
