package service

import (
	"context"
	"errors"
	"github.com/angrymuskrat/event-monitoring-system/services/data-storage/connector"
	"github.com/angrymuskrat/event-monitoring-system/services/proto"
)

var (
	ErrSelectInterval = errors.New("incorrect interval")
)

type Service interface {
	// input array of Posts and write every Post to database
	// return array of statuses of adding posts
	// 		and error if one or more Post wasn't pushed
	PushPosts(ctx context.Context, posts []data.Post) ([]int32, error)

	// input SpatioTemporalInterval
	// return array of post, every of which satisfy the interval conditions
	// 		and error if interval is incorrect or storage can't return posts due to other reasons
	SelectPosts(ctx context.Context, interval data.SpatioTemporalInterval) ([]data.Post, error)
}

type basicService struct{
	db *connector.Storage
}

func (s basicService) PushPosts(_ context.Context, posts []data.Post) ([]int32, error) {
	return s.db.PushPosts(posts)
}

func (s basicService) SelectPosts(_ context.Context, interval data.SpatioTemporalInterval) ([]data.Post, error) {
	if interval.MaxLon < interval.MinLon || interval.MaxLat < interval.MinLat || interval.MaxTime < interval.MinTime {
		return nil, ErrSelectInterval
	}
	return s.db.SelectPosts(interval)
}
