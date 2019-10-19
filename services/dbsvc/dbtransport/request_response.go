package dbtransport

import (
	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc/pb"
)

type PushRequest struct {
	Posts []pb.Post
}

type PushResponse struct {
	Err error
}

type SelectRequest struct {
	Err error
	Interval pb.SpatialTemporalInterval
}

type SelectResponse struct {
	Posts []pb.Post
	Err error
}




