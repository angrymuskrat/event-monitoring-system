package service

import (
	"context"
	"github.com/angrymuskrat/event-monitoring-system/services/data-storage/proto"
	"github.com/angrymuskrat/event-monitoring-system/services/proto"

	"github.com/go-kit/kit/endpoint"
)

// Set collects all of the endpoints that compose an add service
type Set struct {
	PushPostsEndpoint       endpoint.Endpoint
	SelectPostsEndpoint     endpoint.Endpoint
	SelectAggrPostsEndpoint endpoint.Endpoint
	PushGridEndpoint        endpoint.Endpoint
	PullGridEndpoint        endpoint.Endpoint
	PushEventsEndpoint      endpoint.Endpoint
	PullEventsEndpoint      endpoint.Endpoint
	PushLocationsEndpoint   endpoint.Endpoint
	PullLocationsEndpoint   endpoint.Endpoint
}

func NewEndpoint(svc Service) Set {
	return Set{
		PushPostsEndpoint:       makePushPostsEndpoint(svc),
		SelectPostsEndpoint:     makeSelectPostsEndpoint(svc),
		SelectAggrPostsEndpoint: makeSelectAggrPostsEndopoint(svc),
		PushGridEndpoint:        makePushGridEndpoint(svc),
		PullGridEndpoint:        makePullGridEndpoint(svc),
		PushEventsEndpoint:      makePushEventsEndpoint(svc),
		PullEventsEndpoint:      makePullEventsEndpoint(svc),
		PushLocationsEndpoint:   makePushLocationsEndpoint(svc),
		PullLocationsEndpoint:   makePullLocationsEndpoint(svc),
	}
}

func (s Set) PushPosts(ctx context.Context, posts []data.Post) ([]int32, error) {
	resp, err := s.PushPostsEndpoint(ctx, proto.PushPostsRequest{Posts: posts})
	if err != nil {
		return nil, err
	}
	response := resp.(PushPostsResponse)
	return response.Ids, response.Err
}

func (s Set) SelectPosts(ctx context.Context, interval data.SpatioTemporalInterval) ([]data.Post, error) {
	resp, err := s.SelectPostsEndpoint(ctx, proto.SelectPostsRequest{Interval: interval})
	if err != nil {
		return nil, err
	}
	response := resp.(SelectPostsResponse)
	return response.Posts, response.Err
}

func (s Set) SelectAggrPosts(ctx context.Context, interval data.SpatioHourInterval) ([]data.AggregatedPost, error) {
	resp, err := s.SelectAggrPostsEndpoint(ctx, proto.SelectAggrPostsRequest{Interval: interval})
	if err != nil {
		return nil, err
	}
	response := resp.(SelectAggrPostsResponse)
	return response.Posts, response.Err
}

func (s Set) PushGrid(ctx context.Context, id string, blob []byte) error {
	resp, err := s.PushGridEndpoint(ctx, proto.PushGridRequest{Id: id, Blob: blob})
	if err != nil {
		return err
	}
	response := resp.(PushGridResponse)
	return response.Err
}

func (s Set) PullGrid(ctx context.Context, id string) ([]byte, error) {
	resp, err := s.PullGridEndpoint(ctx, proto.PullGridRequest{Id: id})
	if err != nil {
		return nil, err
	}
	response := resp.(PullGridResponse)
	return response.Blob, response.Err
}

func (s Set) PushEvents(ctx context.Context, events []data.Event) error {
	resp, err := s.PushEventsEndpoint(ctx, proto.PushEventsRequest{Events: events})
	if err != nil {
		return err
	}
	response := resp.(PushEventsResponse)
	return response.Err
}

func (s Set) PullEvents(ctx context.Context, interval data.SpatioHourInterval) ([]data.Event, error) {
	resp, err := s.PullEventsEndpoint(ctx, proto.PullEventsRequest{Interval: interval})
	if err != nil {
		return nil, err
	}
	response := resp.(PullEventsResponse)
	return response.Events, response.Err
}

func (s Set) PushLocations(ctx context.Context, city data.City, locations []data.Location) error {
	resp, err := s.PushLocationsEndpoint(ctx, proto.PushLocationsRequest{City: city, Locations: locations})
	if err != nil {
		return err
	}
	response := resp.(PushLocationsResponse)
	return response.Err
}

func (s Set) PullLocations(ctx context.Context, cityId string) ([]data.Location, error) {
	resp, err := s.PullLocationsEndpoint(ctx, proto.PullLocationsRequest{CityId: cityId})
	if err != nil {
		return nil, err
	}
	response := resp.(PullLocationsResponse)
	return response.Locations, response.Err
}

// makeEndpoint functions
func makePushPostsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.PushPostsRequest)
		ids, err := s.PushPosts(ctx, req.Posts)
		return PushPostsResponse{Ids: ids, Err: err}, nil
	}
}

func makeSelectPostsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.SelectPostsRequest)
		posts, err := s.SelectPosts(ctx, req.Interval)
		return SelectPostsResponse{Posts: posts, Err: err}, nil
	}
}

func makeSelectAggrPostsEndopoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.SelectAggrPostsRequest)
		posts, err := s.SelectAggrPosts(ctx, req.Interval)
		return SelectAggrPostsResponse{Posts: posts, Err: err}, nil
	}
}

func makePushGridEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.PushGridRequest)
		err = s.PushGrid(ctx, req.Id, req.Blob)
		return PushGridResponse{Err: err}, nil
	}
}

func makePullGridEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.PullGridRequest)
		blob, err := s.PullGrid(ctx, req.Id)
		return PullGridResponse{Blob: blob, Err: err}, nil
	}
}

func makePushEventsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.PushEventsRequest)
		err = s.PushEvents(ctx, req.Events)
		return PushEventsResponse{Err: err}, nil
	}
}

func makePullEventsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.PullEventsRequest)
		events, err := s.PullEvents(ctx, req.Interval)
		return PullEventsResponse{Events: events, Err: err}, nil
	}
}

func makePushLocationsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.PushLocationsRequest)
		err = s.PushLocations(ctx, req.City, req.Locations)
		return PushLocationsResponse{Err: err}, nil
	}
}

func makePullLocationsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.PullLocationsRequest)
		locations, err := s.PullLocations(ctx, req.CityId)
		return PullLocationsResponse{Locations: locations, Err: err}, nil
	}
}

// compile time assertions for our response types implementing endpoint.Failer.
var (
	_ endpoint.Failer = PushPostsResponse{}
	_ endpoint.Failer = SelectPostsResponse{}

	_ endpoint.Failer = SelectAggrPostsResponse{}

	_ endpoint.Failer = PushGridResponse{}
	_ endpoint.Failer = PullGridResponse{}

	_ endpoint.Failer = PushEventsResponse{}
	_ endpoint.Failer = PullEventsResponse{}

	_ endpoint.Failer = PushLocationsResponse{}
	_ endpoint.Failer = PullLocationsResponse{}
)

type PushPostsResponse struct {
	Ids []int32 `json:"ids"`
	Err error   `json:"-"`
}
type SelectPostsResponse struct {
	Posts []data.Post `json:"posts"`
	Err   error       `json:"-"`
}

type SelectAggrPostsResponse struct {
	Posts []data.AggregatedPost `json:"posts"`
	Err   error                 `json:"-"`
}

type PushGridResponse struct {
	Err error `json:"-"`
}
type PullGridResponse struct {
	Blob []byte `json:"blob"`
	Err  error  `json:"-"`
}

type PushEventsResponse struct {
	Err error `json:"-"`
}
type PullEventsResponse struct {
	Events []data.Event `json:"events"`
	Err    error        `json:"-"`
}

type PushLocationsResponse struct {
	Err error `json:"-"`
}
type PullLocationsResponse struct {
	Locations []data.Location `json:"locations"`
	Err       error           `json:"-"`
}

func (r PushPostsResponse) Failed() error   { return r.Err }
func (r SelectPostsResponse) Failed() error { return r.Err }

func (r SelectAggrPostsResponse) Failed() error { return r.Err }

func (r PushGridResponse) Failed() error { return r.Err }
func (r PullGridResponse) Failed() error { return r.Err }

func (r PushEventsResponse) Failed() error { return r.Err }
func (r PullEventsResponse) Failed() error { return r.Err }

func (r PushLocationsResponse) Failed() error { return r.Err }
func (r PullLocationsResponse) Failed() error { return r.Err }
