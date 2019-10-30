package dbsvc

import (
	"context"
	"errors"
	"github.com/angrymuskrat/event-monitoring-system/services/proto"
	"github.com/go-kit/kit/log"
)

var (
	ErrPushStatement = errors.New("do not be able to push posts")
	ErrSelectStatement = errors.New("do not return posts from DB")
	ErrDBConnector = errors.New("do not create db connection")
	ErrCmdRepository   = errors.New("unable to command repository")
	ErrQueryRepository = errors.New("unable to query repository")
	ErrDbNotCreated    = errors.New("db not created") //tmp error
)

type Service interface {
	Push(ctx context.Context, posts []data.Post) error
	Select(ctx context.Context, interval data.SpatioTemporalInterval) ([]data.Post, error)
}

// New returns a basic Service with all of the expected middlewares wired in.
func NewService(logger log.Logger, db *DBConnector) Service {
	var svc Service
	{
		svc = NewBasicService(db)
		svc = LoggingMiddleware(logger)(svc)
	}
	return svc
}

// NewBasicService returns a naïve, stateless implementation of Service.
func NewBasicService(db *DBConnector) Service {
	return basicService{ db: db }
}

type basicService struct{
	db *DBConnector
}

func (s basicService) Push(ctx context.Context, posts []data.Post) error {
	err := s.db.Push(posts);
	return err
}

// Concat implements Service.
func (s basicService) Select(_ context.Context, interval data.SpatioTemporalInterval) ([]data.Post, error) {
	posts, err := s.db.Select(interval)
	if err != nil {
		return nil, err
	}
	return posts, err
}
