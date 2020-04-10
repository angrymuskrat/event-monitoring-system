package service

import (
	"context"
	"time"

	"github.com/angrymuskrat/event-monitoring-system/services/data-storage/storage"
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
)

const (
	// Max time of waiting of execution of request for client; time.Duration
	TimeWaitingClient = 30 * time.Second

	// Max size for income messages of grpcService; in bytes
	// For client is needed also set max income message side, but this is done by the client during initialization of grpcClient
	MaxMsgSize = 15000000000
)

type Service interface {
	// input: context, a city object
	// output: error
	// result: if the city was successfully added, will return nil, otherwise will return an error
	InsertCity(ctx context.Context, city data.City, updateIfExists bool) error

	// input: context
	// output: array of city objects, error
	// result: if request was successfully finished, will return cities and nil error, otherwise will return empty array and error
	GetAllCities(ctx context.Context) ([]data.City, error)

	// input: context, city.Code - id of the city
	// output: city and error
	// result: if request was successfully finished, will return the city and nil error, otherwise return nil and some error
	GetCity(ctx context.Context, cityId string) (*data.City, error)

	// input: context, cityId string, array of posts
	// output: array of commit statuses (see data-storage/data/typePushResponse) and error
	// if all posts were successfully added to the city's db, will return statuses and nil error, otherwise statuses and some error
	// Either all posts will be added or not a single one.
	PushPosts(ctx context.Context, cityId string, posts []data.Post) error

	// input: context, id of the city, start and finish UTC-time in second - time interval, for which posts of this city will be returned
	// output: array of posts, area object, error
	// result: if request was successfully finished, will return posts and Border of city - TopLeft and BotRight Points and nil error,
	// 		otherwise, empty array, nil area and some error
	SelectPosts(ctx context.Context, cityId string, startTime, finishTime int64) ([]data.Post, *data.Area, error)

	// input: context, id of the city, interval, which contains UTC-time in second - start time of needed hour and
	// 		area - TopLeft and BotRight Points of needed space
	// output: array of aggregated posts, each aggr post has coordinate of its aggregated cell, and amount of posts in this hour and this cell
	// result: if request was successfully finished will return aggregated posts and nil error, otherwise empty array and some error
	SelectAggrPosts(ctx context.Context, cityId string, interval data.SpatioHourInterval) ([]data.AggregatedPost, error)

	// input: context, id of the city, start hour and finish hour UTC-time in seconds, both are beginning of needed hours
	// output: array of timestamps and error
	// result: if request was successfully finished, will return timeline - amount of posts and events in this city between start and finish hour
	PullTimeline(ctx context.Context, cityId string, start, finish int64) ([]data.Timestamp, error)

	// input: context, id of the city, map of grids, keys of this map are ids and value is byte array - historic grid.
	// 		ids description: first two digit: month; 3th: 0 - work day,  1 - holiday; 4th and 5th: hour.
	//		Example: 02013 - 02 March, 0 work day, 13 o'clock
	// output: error
	// if all grids were successfully added to the city's db, will return nil error, otherwise statuses and some error
	// Either all grids will be added or not a single one.
	PushGrid(ctx context.Context, cityId string, grids map[int64][]byte) error

	// input: context, id of the city, start and finish ids
	// 		ids description: first two digit: month; 3th: 0 - work day,  1 - holiday; 4th and 5th: hour.
	//		Example: 02013 - 02 March, 0 work day, 13 o'clock
	// output: map of map of grids, keys of this map are ids and value is byte array
	// result: if request was successfully finished, will return grids and nil error otherwise return nil map and some error
	PullGrid(ctx context.Context, cityId string, ids []int64) (map[int64][]byte, error)

	// input: context, id of the city, array od events
	// output: error
	// if all events were successfully added to the city's db, will return nil error, otherwise statuses and some error
	// Either all events will be added or not a single one.
	PushEvents(ctx context.Context, cityId string, events []data.Event) error

	// input: context, id of the city, array od events
	// output: error
	// if all events were successfully updated to the city's db, will return nil error, otherwise statuses and some error
	// Either all events will be added or not a single one.
	UpdateEvents(ctx context.Context, cityId string, events []data.Event) error

	// input: context, id of the city, interval, which contains UTC-time in second - start time of needed hour and
	// 		area - TopLeft and BotRight Points of needed space
	// output: events and error
	// result: if request was successfully finished, will return events and nil error, otherwise, empty array and some error
	PullEvents(ctx context.Context, cityId string, interval data.SpatioHourInterval) ([]data.Event, error)

	// input: context, id of the city, array of tags, start hour and finish hour UTC-time in seconds, both are beginning of needed hours
	// output: events and error
	// result: if request was successfully finished, will return events,
	//		which include all input tags and have time between startTime and finishTime, and nil error, otherwise, empty array and some error
	PullEventsTags(ctx context.Context, cityId string, tags []string, startTime, finishTime int64) ([]data.Event, error)

	// input: context, id of the city, start hour and finish hour UTC-time in seconds, both are beginning of needed hours
	// output: events, posts and error
	// result: if request was successfully finished, will return events with ids,
	//		which have time between startTime and finishTime, their posts and nil error, otherwise, empty array and some error
	PullEventsWithIDs(ctx context.Context, cityId string, startTime, finishTime int64) ([]data.Event, error)

	// input: context, id of the city, array of ids of events
	// output: error
	// result: if events deleted successfully, will return nil error
	DeleteEvents(ctx context.Context, cityId string, ids []int64) error

	// input: context, id of the city, array of locations - Instagram's locations for this city
	// output: error
	// if all locations were successfully added to the city's db, will return nil error, otherwise statuses and some error
	// Either all locations will be added or not a single one.
	PushLocations(ctx context.Context, cityId string, locations []data.Location) error

	// input: context, id of the city
	// output: array of locations
	// result: if request was successfully finished, return all locations of this city and nil error, otherwise return some error
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

func (s basicService) PushPosts(ctx context.Context, cityId string, posts []data.Post) error {
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

func (s basicService) PullGrid(ctx context.Context, cityId string, ids []int64) (map[int64][]byte, error) {
	return s.db.PullGrid(ctx, cityId, ids)
}

func (s basicService) PushEvents(ctx context.Context, cityId string, events []data.Event) error {
	return s.db.PushEvents(ctx, cityId, events)
}

func (s basicService) UpdateEvents(ctx context.Context, cityId string, events []data.Event) error {
	return s.db.UpdateEvents(ctx, cityId, events)
}

func (s basicService) PullEvents(ctx context.Context, cityId string, interval data.SpatioHourInterval) ([]data.Event, error) {
	return s.db.PullEvents(ctx, cityId, interval)
}

func (s basicService) PullEventsTags(ctx context.Context, cityId string, tags []string, startTime, finishTime int64) ([]data.Event, error) {
	return s.db.PullEventsTags(ctx, cityId, tags, startTime, finishTime)
}

func (s basicService) PullEventsWithIDs(ctx context.Context, cityId string, startTime, finishTime int64) ([]data.Event, error) {
	return s.db.PullEventsWithIDs(ctx, cityId, startTime, finishTime)
}

func (s basicService) DeleteEvents(ctx context.Context, cityId string, ids []int64) error {
	return s.db.DeleteEvents(ctx, cityId, ids)
}

func (s basicService) PushLocations(ctx context.Context, cityId string, locations []data.Location) error {
	return s.db.PushLocations(ctx, cityId, locations)
}

func (s basicService) PullLocations(ctx context.Context, cityId string) ([]data.Location, error) {
	return s.db.PullLocations(ctx, cityId)
}
