package service

import (
	"context"
	"errors"
	"github.com/angrymuskrat/event-monitoring-system/services/data-storage/storage"
	"github.com/angrymuskrat/event-monitoring-system/services/proto"
	"time"
)

var (
	ErrEmptyGridId = errors.New("empty grid id")
	ErrEmptyGrid   = errors.New("empty grid")
)

const (
	TimeWaitingClient = 30 * time.Second // in seconds
	MaxMsgSize        = 1000000000       // in bytes
)

type Service interface {
	// input array of Posts and write every Post to database
	// return array of statuses of adding posts
	// 		and error if one or more Post wasn't pushed
	PushPosts(ctx context.Context, cityId string, posts []data.Post) ([]int32, error)

	// input SpatioTemporalInterval
	// return array of post, every of which satisfy the interval conditions
	// 		and error if storage can't return posts due to other reasons
	SelectPosts(ctx context.Context, cityId string, interval data.SpatioTemporalInterval) ([]data.Post, error)

	SelectAggrPosts(ctx context.Context, cityId string, interval data.SpatioHourInterval) ([]data.AggregatedPost, error)

	PullTimeline(ctx context.Context, cityId string, start, finish int64) ([]data.Timestamp, error)

	// input not empty id and not empty array of bytes
	// if blob successfully added to database, return nil
	// else return error
	PushGrid(ctx context.Context, cityId string, id string, blob []byte) error

	// input not empty id
	// if there are exist some blob with the id in database return the blob
	// else return error
	PullGrid(ctx context.Context, cityId string, id string) ([]byte, error)

	PushEvents(ctx context.Context, cityId string, events []data.Event) error

	PullEvents(ctx context.Context, cityId string, interval data.SpatioHourInterval) ([]data.Event, error)

	PushLocations(ctx context.Context, cityId string, locations []data.Location) error

	PullLocations(ctx context.Context, cityId string) ([]data.Location, error)
}

type basicService struct {
	db *storage.Storage
}

func (s basicService) PushPosts(ctx context.Context, cityId string, posts []data.Post) ([]int32, error) {
	return s.db.PushPosts(ctx, posts)
}

func (s basicService) SelectPosts(ctx context.Context, cityId string, interval data.SpatioTemporalInterval) ([]data.Post, error) {
	return s.db.SelectPosts(ctx, interval)
}

func (s basicService) SelectAggrPosts(ctx context.Context, cityId string, interval data.SpatioHourInterval) ([]data.AggregatedPost, error) {
	return s.db.SelectAggrPosts(ctx, interval)
}

func (s basicService) PullTimeline(ctx context.Context, cityId string, start, finish int64) ([]data.Timestamp, error) {
	return s.db.PullTimeline(ctx, cityId, start, finish)
}

func (s basicService) PushGrid(ctx context.Context, cityId string, id string, blob []byte) error {
	if id == "" {
		return ErrEmptyGridId
	}
	if blob == nil || len(blob) == 0 {
		return ErrEmptyGrid
	}
	return s.db.PushGrid(ctx, id, blob)
}

func (s basicService) PullGrid(ctx context.Context, cityId string, id string) ([]byte, error) {
	if id == "" {
		return nil, ErrEmptyGridId
	}
	return s.db.PullGrid(ctx, id)
}

func (s basicService) PushEvents(ctx context.Context, cityId string, events []data.Event) error {
	return s.db.PushEvents(ctx, events)
}

func (s basicService) PullEvents(ctx context.Context, cityId string, interval data.SpatioHourInterval) ([]data.Event, error) {
	return s.db.PullEvents(ctx, interval)
}

func (s basicService) PushLocations(ctx context.Context, cityId string, locations []data.Location) error {
	return s.db.PushLocations(ctx, cityId, locations)
}

func (s basicService) PullLocations(ctx context.Context, cityId string) ([]data.Location, error) {
	return s.db.PullLocations(ctx, cityId)
}
