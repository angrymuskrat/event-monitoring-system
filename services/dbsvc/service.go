package dbsvc

import (
	"context"
	"errors"
	"fmt"
	"github.com/angrymuskrat/event-monitoring-system/services/proto"
	"github.com/go-kit/kit/log"
)

var (
	ErrCmdRepository   = errors.New("unable to command repository")
	ErrQueryRepository = errors.New("unable to query repository")
	ErrDbNotCreated    = errors.New("db not created") //tmp error
)

type Service interface {
	Push(ctx context.Context, posts []proto.Post) error
	Select(ctx context.Context, interval proto.SpatioTemporalInterval) ([]proto.Post, error)
}

// New returns a basic Service with all of the expected middlewares wired in.
func NewService(logger log.Logger) Service {
	var svc Service
	{
		svc = NewBasicService()
		svc = LoggingMiddleware(logger)(svc)
	}
	return svc
}

// NewBasicService returns a naïve, stateless implementation of Service.
func NewBasicService() Service {
	return basicService{}
}

type basicService struct{}

func (s basicService) Push(ctx context.Context, posts []proto.Post) error {
	for _, post := range posts {
		fmt.Print(post.ID)
	}
	return nil
}

// Concat implements Service.
func (s basicService) Select(_ context.Context, interval proto.SpatioTemporalInterval) ([]proto.Post, error) {
	return nil, ErrDbNotCreated
}
