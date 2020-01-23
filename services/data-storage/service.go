package service

import (
	"context"
	"errors"
	"github.com/angrymuskrat/event-monitoring-system/services/data-storage/connector"
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
	PushPosts(ctx context.Context, posts []data.Post) ([]int32, error)

	// input SpatioTemporalInterval
	// return array of post, every of which satisfy the interval conditions
	// 		and error if storage can't return posts due to other reasons
	SelectPosts(ctx context.Context, interval data.SpatioTemporalInterval) ([]data.Post, error)

	SelectAggrPosts(ctx context.Context, interval data.SpatioHourInterval) ([]data.AggregatedPost, error)

	// input not empty id and not empty array of bytes
	// if blob successfully added to database, return nil
	// else return error
	PushGrid(ctx context.Context, id string, blob []byte) error

	// input not empty id
	// if there are exist some blob with the id in database return the blob
	// else return error
	PullGrid(ctx context.Context, id string) ([]byte, error)

	PushEvents(ctx context.Context, events []data.Event) error

	PullEvents(ctx context.Context, interval data.SpatioHourInterval) ([]data.Event, error)

	PushLocations(ctx context.Context, city data.City, locations []data.Location) error

	PullLocations(ctx context.Context, cityId string) ([]data.Location, error)
}

type basicService struct {
	db *connector.Storage
}

func (s basicService) PushPosts(_ context.Context, posts []data.Post) ([]int32, error) {
	return s.db.PushPosts(posts)
}

func (s basicService) SelectPosts(_ context.Context, interval data.SpatioTemporalInterval) ([]data.Post, error) {
	return s.db.SelectPosts(interval)
}

func (s basicService) SelectAggrPosts(_ context.Context, interval data.SpatioHourInterval) ([]data.AggregatedPost, error) {
	return s.db.SelectAggrPosts(interval)
}

func (s basicService) PushGrid(_ context.Context, id string, blob []byte) error {
	if id == "" {
		return ErrEmptyGridId
	}
	if blob == nil || len(blob) == 0 {
		return ErrEmptyGrid
	}
	return s.db.PushGrid(id, blob)
}

func (s basicService) PullGrid(_ context.Context, id string) ([]byte, error) {
	if id == "" {
		return nil, ErrEmptyGridId
	}
	return s.db.PullGrid(id)
}

func (s basicService) PushEvents(_ context.Context, events []data.Event) error {
	return s.db.PushEvents(events)
}

func (s basicService) PullEvents(_ context.Context, interval data.SpatioHourInterval) ([]data.Event, error) {
	return s.db.PullEvents(interval)
}

func (s basicService) PushLocations(_ context.Context, city data.City, locations []data.Location) error {
	return s.db.PushLocations(city, locations)
}

func (s basicService) PullLocations(_ context.Context, cityId string) ([]data.Location, error) {
	return s.db.PullLocations(cityId)
}
