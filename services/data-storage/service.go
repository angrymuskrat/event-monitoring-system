package service

import (
	"context"
	"errors"
	"github.com/angrymuskrat/event-monitoring-system/services/data-storage/connector"
	"github.com/angrymuskrat/event-monitoring-system/services/proto"
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

type basicService struct{
	db *connector.Storage
}

func (s basicService) Push(_ context.Context, posts []data.Post) ([]int32, error) {
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
