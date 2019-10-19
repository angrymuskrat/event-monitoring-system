package transport

import "github.com/angrymuskrat/event-monitoring-system/services/dbsvc"

type PushRequest struct {
	Posts []dbsvc.Post
}

type PushResponse struct {
	Err error
}

type SelectRequest struct {
	Err error
	Interval dbsvc.SpatialTemporalInterval
}

type SelectResponse struct {
	Posts []dbsvc.Post
	Err error
}




