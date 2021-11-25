package service

import data "github.com/angrymuskrat/event-monitoring-system/services/proto"

// TODO: this is temporary solution of problem with timestamps!
func getDelta(city string) int64 {
	var delta int64
	switch city {
	case "nyc":
		delta = 3600 * 7
	case "london":
		delta = 3600 * 3
	case "spb":
		delta = 0
	case "moscow":
		delta = 0
	default:
		delta = 3600 * 3
	}
	return delta
}

func fixTimestamp(city string, timestamp int64) int64 {
	return timestamp + getDelta(city)
}

func fixReverseTimestamp(city string, timestamp int64) int64 {
	return timestamp - getDelta(city)
}

func fixTimeline(city string, timeline Timeline) {
	for i := range timeline {
		timeline[i].Time = fixReverseTimestamp(city, timeline[i].Time)
	}
}

func fixEvents(city string, events []data.Event) {
	for i := range events {
		events[i].Start = fixReverseTimestamp(city, events[i].Start)
		events[i].Finish = fixReverseTimestamp(city, events[i].Finish)
	}
}

func fixPost(city string, post *data.ShortPost) {
	post.Timestamp = fixReverseTimestamp(city, post.Timestamp)
}

func fixPosts(city string, posts []data.ShortPost) {
	for i := range posts {
		posts[i].Timestamp = fixReverseTimestamp(city, posts[i].Timestamp)
	}
}

type backendService struct {
	storageConn StorageConnector
}

func (s *backendService) HeatmapPosts(req HeatmapRequest) ([]data.AggregatedPost, error) {
	return s.storageConn.HeatmapPosts(req.City, req.TopLeft, req.BottomRight, fixTimestamp(req.City, req.Hour))
}

func (s *backendService) Timeline(req TimelineRequest) (Timeline, error) {
	timeline, err := s.storageConn.Timeline(req.City, fixTimestamp(req.City, req.Start), fixTimestamp(req.City, req.Finish))
	if err == nil {
		fixTimeline(req.City, timeline)
	}
	return timeline, err
}

func (s *backendService) Events(req EventsRequest) ([]data.Event, error) {
	events, err := s.storageConn.Events(req.City, req.TopLeft, req.BottomRight, fixTimestamp(req.City, req.Hour))
	if err == nil {
		fixEvents(req.City, events)
	}
	return events, err
}

func (s *backendService) SearchEvents(req SearchRequest) ([]data.Event, error) {
	events, err := s.storageConn.EventsByTags(req.City, req.Keytags, fixTimestamp(req.City, req.Start), fixTimestamp(req.City, req.Finish))
	if err == nil {
		fixEvents(req.City, events)
	}
	return events, err
}

func (s *backendService) ShortPostsInInterval(req ShortPostsRequest) ([]data.ShortPost, error) {
	posts, err := s.storageConn.ShortPostsInInterval(req.City, req.Shortcodes, fixTimestamp(req.City, req.Start), fixTimestamp(req.City, req.End))
	if err == nil {
		fixPosts(req.City, posts)
	}
	return posts, err
}

func (s *backendService) SingleShortPost(req SingleShortPostRequest) (*data.ShortPost, error) {
	post, err := s.storageConn.SingleShortPost(req.City, req.Shortcode)
	if err == nil {
		fixPost(req.City, post)
	}
	return post, err
}
