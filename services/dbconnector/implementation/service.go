package implementation

import (
	"context"
	"database/sql"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	dbconnectorsvc "github.com/angrymuskrat/event-monitoring-system/services/dbconnector"
)

// service implements the dbconnector Service
type service struct {
	repository dbconnectorsvc.Repository
	logger     log.Logger
}

// NewService creates and returns a new dbconnector service instance
func NewService(rep dbconnectorsvc.Repository, logger log.Logger) dbconnectorsvc.Service {
	return &service{
		repository: rep,
		logger:     logger,
	}
}

// Push save array of Posts to db
func (s *service) Push(ctx context.Context, posts [] dbconnectorsvc.Post) error {
	logger := log.With(s.logger, "method", "Push")
	if err := s.repository.PushPosts(ctx, posts); err != nil {
		level.Error(logger).Log("err", err)
		return dbconnectorsvc.ErrCmdRepository
	}
	return nil
}


// Select returns array of posts, which  SpatialTemporal coordinates are in given interval
func (s *service) Select(ctx context.Context, interval dbconnectorsvc.SpatialTemporalInterval) ([]dbconnectorsvc.Post, error) {
	logger := log.With(s.logger, "method", "Select")
	posts, err := s.repository.SelectPosts(ctx, interval)
	if err != nil && err != sql.ErrNoRows {
		level.Error(logger).Log("err", err)
		return posts, dbconnectorsvc.ErrQueryRepository
	}
	return posts, nil
}



