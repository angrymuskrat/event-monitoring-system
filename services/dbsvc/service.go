package dbsvc

import (
	"context"
	"errors"
)

var (
	ErrCmdRepository   = errors.New("unable to command repository")
	ErrQueryRepository = errors.New("unable to query repository")
)

type Service interface {
	Push (ctx context.Context, posts []Post) error
	Select (ctx context.Context, interval SpatialTemporalInterval) ([]Post, error)
}

