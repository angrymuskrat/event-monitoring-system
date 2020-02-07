package service

import (
	"context"
	"github.com/angrymuskrat/event-monitoring-system/services/data-storage/storage"
	"github.com/angrymuskrat/event-monitoring-system/services/proto"
	"time"
)

const (
	TimeWaitingClient = 30 * time.Second // in seconds
	MaxMsgSize        = 1000000000       // in bytes
)

type Service interface {
	InsertCity(ctx context.Context, city data.City, updateIfExists bool) error

	GetAllCities(ctx context.Context) ([]data.City, error)

	GetCity(ctx context.Context, cityId string) (*data.City, error)

	// input array of Posts and write every Post to database
	// return array of statuses of adding posts
	// 		and error if one or more Post wasn't pushed
	PushPosts(ctx context.Context, cityId string, posts []data.Post) ([]int32, error)

	// input SpatioTemporalInterval
	// return array of post, every of which satisfy the interval conditions
	// 		and error if storage can't return posts due to other reasons
	SelectPosts(ctx context.Context, cityId string, startTime, finishTime int64) ([]data.Post, *data.Area, error)

	SelectAggrPosts(ctx context.Context, cityId string, interval data.SpatioHourInterval) ([]data.AggregatedPost, error)

	PullTimeline(ctx context.Context, cityId string, start, finish int64) ([]data.Timestamp, error)

	// input not empty id and not empty array of bytes
	// if blob successfully added to database, return nil
	// else return error
	PushGrid(ctx context.Context, cityId string, grids map[int64][]byte) error

	// input not empty id
	// if there are exist some blob with the id in database return the blob
	// else return error
	PullGrid(ctx context.Context, cityId string, startId, finishId int64) (map[int64][]byte, error)

	PushEvents(ctx context.Context, cityId string, events []data.Event) error

	PullEvents(ctx context.Context, cityId string, interval data.SpatioHourInterval) ([]data.Event, error)

	PullEventsTags(ctx context.Context, cityId string, tags []string, startTime, finishTime int64) ([]data.Event, error)

	PushLocations(ctx context.Context, cityId string, locations []data.Location) error

	PullLocations(ctx context.Context, cityId string) ([]data.Location, error)
}

type basicService struct {
	db *storage.Storage
}

func (s basicService) InsertCity(ctx context.Context, city data.City, updateIfExists bool) error {
	return s.db.InsertCity(ctx, city, updateIfExists)
}

func (s basicService) GetAllCities(ctx context.Context) ([]data.City, error) {
	return s.db.GetCities(ctx)
}

func (s basicService) GetCity(ctx context.Context, cityId string) (*data.City, error) {
	return s.db.SelectCity(ctx, cityId)
}

func (s basicService) PushPosts(ctx context.Context, cityId string, posts []data.Post) ([]int32, error) {
	return s.db.PushPosts(ctx, cityId, posts)
}

func (s basicService) SelectPosts(ctx context.Context, cityId string, startTime, finishTime int64) ([]data.Post, *data.Area, error) {
	return s.db.SelectPosts(ctx, cityId, startTime, finishTime)
}

func (s basicService) SelectAggrPosts(ctx context.Context, cityId string, interval data.SpatioHourInterval) ([]data.AggregatedPost, error) {
	return s.db.SelectAggrPosts(ctx, cityId, interval)
}

func (s basicService) PullTimeline(ctx context.Context, cityId string, start, finish int64) ([]data.Timestamp, error) {
	return s.db.PullTimeline(ctx, cityId, start, finish)
}

func (s basicService) PushGrid(ctx context.Context, cityId string, grids map[int64][]byte) error {
	return s.db.PushGrid(ctx, cityId, grids)
}

func (s basicService) PullGrid(ctx context.Context, cityId string, startId, finishId int64) (map[int64][]byte, error) {
	return s.db.PullGrid(ctx, cityId, startId, finishId)
}

func (s basicService) PushEvents(ctx context.Context, cityId string, events []data.Event) error {
	return s.db.PushEvents(ctx, cityId, events)
}

func (s basicService) PullEvents(ctx context.Context, cityId string, interval data.SpatioHourInterval) ([]data.Event, error) {
	return s.db.PullEvents(ctx, cityId, interval)
}

func (s basicService) PullEventsTags(ctx context.Context, cityId string, tags []string, startTime, finishTime int64) ([]data.Event, error) {
	return s.db.PullEventsTags(ctx, cityId, tags, startTime, finishTime)
}

func (s basicService) PushLocations(ctx context.Context, cityId string, locations []data.Location) error {
	return s.db.PushLocations(ctx, cityId, locations)
}

func (s basicService) PullLocations(ctx context.Context, cityId string) ([]data.Location, error) {
	return s.db.PullLocations(ctx, cityId)
}
