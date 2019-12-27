package service

import (
	"context"
	"errors"
	"fmt"
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

	PushGrid(ctx context.Context, id string, blob []byte) error

	PullGrid(ctx context.Context, id string) ([]byte, error)
}

type basicService struct {
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

func (s basicService) PushGrid(_ context.Context, id string, blob []byte) error {
	fmt.Printf("ID: %v, size: %v", id, len(blob))
	return nil
}

func (s basicService) PullGrid(_ context.Context, id string) ([]byte, error) {
	fmt.Printf("ID: %v", id)
	return []byte("test"), nil
}
