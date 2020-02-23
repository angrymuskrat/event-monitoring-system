package service

import (
	"context"
	"github.com/angrymuskrat/event-monitoring-system/services/data-storage/proto"
	"github.com/go-kit/kit/endpoint"
)

func makeInsertCityEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.InsertCityRequest)
		err = s.InsertCity(ctx, req.City, req.UpdateIfExists)
		var msg string
		if err != nil {
			msg = err.Error()
		}
		return proto.InsertCityReply{Err: msg}, nil
	}
}

func makeGetAllCitiesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (response interface{}, err error) {
		cities, err := s.GetAllCities(ctx)
		var msg string
		if err != nil {
			msg = err.Error()
		}
		return proto.GetAllCitiesReply{Cities: cities, Err: msg}, nil
	}
}

func makeGetCityEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.GetCityRequest)
		city, err := s.GetCity(ctx, req.CityId)
		var msg string
		if err != nil {
			msg = err.Error()
		}
		return proto.GetCityReply{City: city, Err: msg}, nil
	}
}

func makePushPostsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.PushPostsRequest)
		err = s.PushPosts(ctx, req.CityId, req.Posts)
		var msg string
		if err != nil {
			msg = err.Error()
		}
		return proto.PushPostsReply{Err: msg}, nil
	}
}

func makeSelectPostsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.SelectPostsRequest)
		posts, area, err := s.SelectPosts(ctx, req.CityId, req.StartTime, req.FinishTime)
		var msg string
		if err != nil {
			msg = err.Error()
		}
		return proto.SelectPostsReply{Posts: posts, Area: area, Err: msg}, nil
	}
}

func makeSelectAggrPostsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.SelectAggrPostsRequest)
		posts, err := s.SelectAggrPosts(ctx, req.CityId, req.Interval)
		var msg string
		if err != nil {
			msg = err.Error()
		}
		return proto.SelectAggrPostsReply{Posts: posts, Err: msg}, nil
	}
}

func makePullTimelineEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.PullTimelineRequest)
		timeline, err := s.PullTimeline(ctx, req.CityId, req.Start, req.Finish)
		var msg string
		if err != nil {
			msg = err.Error()
		}
		return proto.PullTimelineReply{Timeline: timeline, Err: msg}, nil
	}
}

func makePushGridEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.PushGridRequest)
		err = s.PushGrid(ctx, req.CityId, req.Grids)
		var msg string
		if err != nil {
			msg = err.Error()
		}
		return proto.PushGridReply{Err: msg}, nil
	}
}

func makePullGridEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.PullGridRequest)
		grids, err := s.PullGrid(ctx, req.CityId, req.Ids)
		var msg string
		if err != nil {
			msg = err.Error()
		}
		return proto.PullGridReply{Grids: grids, Err: msg}, nil
	}
}

func makePushEventsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.PushEventsRequest)
		err = s.PushEvents(ctx, req.CityId, req.Events)
		var msg string
		if err != nil {
			msg = err.Error()
		}
		return proto.PushEventsReply{Err: msg}, nil
	}
}

func makePullEventsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.PullEventsRequest)
		events, err := s.PullEvents(ctx, req.CityId, req.Interval)
		var msg string
		if err != nil {
			msg = err.Error()
		}
		return proto.PullEventsReply{Events: events, Err: msg}, nil
	}
}

func makePullEventsTagsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.PullEventsTagsRequest)
		events, err := s.PullEventsTags(ctx, req.CityId, req.Tags, req.StartTime, req.FinishTime)
		var msg string
		if err != nil {
			msg = err.Error()
		}
		return proto.PullEventsTagsReply{Events: events, Err: msg}, nil
	}
}

func makePushLocationsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.PushLocationsRequest)
		err = s.PushLocations(ctx, req.CityId, req.Locations)
		var msg string
		if err != nil {
			msg = err.Error()
		}
		return proto.PushLocationsReply{Err: msg}, nil
	}
}

func makePullLocationsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(proto.PullLocationsRequest)
		locations, err := s.PullLocations(ctx, req.CityId)
		var msg string
		if err != nil {
			msg = err.Error()
		}
		return proto.PullLocationsReply{Locations: locations, Err: msg}, nil
	}
}
