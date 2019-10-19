package implementation

import (
	"context"
	"database/sql"
	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc/dbservice"
	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc/pb"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc"
)

// service implements the dbsvc Service
type service struct {
	repository dbsvc.Repository
	logger     log.Logger
}

// NewService creates and returns a new dbsvc service instance
func NewService(rep dbsvc.Repository, logger log.Logger) dbservice.Service {
	return &service{
		repository: rep,
		logger:     logger,
	}
}

// Push save array of Posts to db
func (s *service) Push(ctx context.Context, posts [] pb.Post) error {
	logger := log.With(s.logger, "method", "Push")
	if err := s.repository.PushPosts(ctx, posts); err != nil {
		level.Error(logger).Log("err", err)
		return dbservice.ErrCmdRepository
	}
	return nil
}


// Select returns array of posts, which  SpatialTemporal coordinates are in given interval
func (s *service) Select(ctx context.Context, interval pb.SpatialTemporalInterval) ([]pb.Post, error) {
	logger := log.With(s.logger, "method", "Select")
	posts, err := s.repository.SelectPosts(ctx, interval)
	if err != nil && err != sql.ErrNoRows {
		level.Error(logger).Log("err", err)
		return posts, dbservice.ErrQueryRepository
	}
	return posts, nil
}



