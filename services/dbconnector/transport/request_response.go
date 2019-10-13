package transport

import "github.com/angrymuskrat/event-monitoring-system/services/dbconnector"

type PushRequest struct {
	Posts []dbconnector.Post
}

type PushResponse struct {
	Err error
}

type SelectRequest struct {
	Err error
	Interval dbconnector.SpatialTemporalInterval
}

type SelectResponse struct {
	Posts []dbconnector.Post
	Err error
}




