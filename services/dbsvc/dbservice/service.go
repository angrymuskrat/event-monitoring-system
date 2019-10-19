package dbservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc/pb"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

var (
	ErrCmdRepository   = errors.New("unable to command repository")
	ErrQueryRepository = errors.New("unable to query repository")
	ErrDbNotCreated = errors.New("db not created") //tmp error
)

type Service interface {
	Push (ctx context.Context, posts []pb.Post) error
	Select (ctx context.Context, interval pb.SpatialTemporalInterval) ([]pb.Post, error)
}

// New returns a basic Service with all of the expected middlewares wired in.
func New(logger log.Logger, pushed, selected metrics.Counter) Service {
	var svc Service
	{
		svc = NewBasicService()
		svc = LoggingMiddleware(logger)(svc)
		svc = InstrumentingMiddleware(pushed, selected)(svc)
	}
	return svc
}

// NewBasicService returns a naïve, stateless implementation of Service.
func NewBasicService() Service {
	return basicService{}
}

type basicService struct{}

func (s basicService) Push (ctx context.Context, posts []pb.Post) error {
	for _, post := range posts {
		fmt.Print(post.ID)
	}
	return nil
}

// Concat implements Service.
func (s basicService) Select(_ context.Context, interval pb.SpatialTemporalInterval) ([]pb.Post, error) {
	fmt.Print(interval.ForLog())
	return nil, ErrDbNotCreated
}