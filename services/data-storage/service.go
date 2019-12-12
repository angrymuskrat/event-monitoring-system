package service

import (
	"context"
	"errors"
	"github.com/angrymuskrat/event-monitoring-system/services/proto"
	"github.com/go-kit/kit/log"
)

var (
	ErrDBConnecting = errors.New("do not be able to connect with db")
	ErrPushStatement = errors.New("one or more posts wasn't pushed")
	ErrSelectInterval = errors.New("incorrect interval")
	ErrSelectStatement = errors.New("don't be able to return posts")
)

type Service interface {
	// input array of Posts and write every Post to database
	// return array of statuses of adding posts
	// 		and error if one or more Post wasn't pushed
	Push(ctx context.Context, posts []data.Post) ([]int32, error)

	// input SpatioTemporalInterval
	// return array of post, every of which satisfy the interval conditions
	// 		and error if interval is incorrect or storage can't return posts due to other reasons
	Select(ctx context.Context, interval data.SpatioTemporalInterval) ([]data.Post, error)
}

func NewService(logger log.Logger, db *Storage) Service {
	var svc Service
	{
		svc = NewBasicService(db)
		svc = LoggingMiddleware(logger)(svc)
	}
	return svc
}

func NewBasicService(db *Storage) Service {
	return basicService{ db: db }
}

type basicService struct{
	db *Storage
}

func (s basicService) Push(ctx context.Context, posts []data.Post) ([]int32, error) {
	ids, err := s.db.Push(posts)
	if err != nil {
		err = ErrPushStatement
	}
	return ids, err
}

func (s basicService) Select(_ context.Context, interval data.SpatioTemporalInterval) ([]data.Post, error) {
	if interval.MaxLon < interval.MinLon || interval.MaxLat < interval.MinLat || interval.MaxTime < interval.MinTime {
		return nil, ErrSelectInterval
	}
	posts, err := s.db.Select(interval)
	if err != nil {
		return nil, ErrSelectStatement
	}
	return posts, nil
}
