package service

import (
	"context"
	"errors"

	"github.com/angrymuskrat/event-monitoring-system/services/data-storage/proto"
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
)

type GrpcService struct {
	insertCity      endpoint.Endpoint
	getAllCities    endpoint.Endpoint
	getCity         endpoint.Endpoint
	pushPosts       endpoint.Endpoint
	selectPosts     endpoint.Endpoint
	selectAggrPosts endpoint.Endpoint
	pullTimeline    endpoint.Endpoint
	pushGrid        endpoint.Endpoint
	pullGrid        endpoint.Endpoint
	pushEvents      endpoint.Endpoint
	pullEvents      endpoint.Endpoint
	pushLocations   endpoint.Endpoint
	pullLocations   endpoint.Endpoint
}

func (svc GrpcService) InsertCity(ctx context.Context, city data.City, updateIfExists bool) error {
	resp, err := svc.insertCity(ctx, proto.InsertCityRequest{City: city, UpdateIfExists: updateIfExists})
	if err != nil {
		return err
	}
	response := resp.(proto.InsertCityReply)
	if response.Err != "" {
		return errors.New(response.Err)
	}
	return nil
}

func (svc GrpcService) GetAllCities(ctx context.Context) ([]data.City, error) {
	resp, err := svc.getAllCities(ctx, proto.GetAllCitiesRequest{})
	if err != nil {
		return nil, err
	}
	response := resp.(proto.GetAllCitiesReply)
	if response.Err != "" {
		return nil, errors.New(response.Err)
	}
	return response.Cities, nil
}

func (svc GrpcService) GetCity(ctx context.Context, cityId string) (*data.City, error) {
	resp, err := svc.getCity(ctx, proto.GetCityRequest{CityId: cityId})
	if err != nil {
		return nil, err
	}
	response := resp.(proto.GetCityReply)
	if response.Err != "" {
		return nil, errors.New(response.Err)
	}
	return response.City, nil
}

func (svc GrpcService) PushPosts(ctx context.Context, cityId string, posts []data.Post) ([]int32, error) {
	resp, err := svc.pushPosts(ctx, proto.PushPostsRequest{CityId: cityId, Posts: posts})
	if err != nil {
		return nil, err
	}
	response := resp.(proto.PushPostsReply)
	if response.Err != "" {
		return nil, errors.New(response.Err)
	}
	return response.Ids, nil
}

func (svc GrpcService) SelectPosts(ctx context.Context, cityId string, startTime, finishTime int64) ([]data.Post, *data.Area, error) {
	resp, err := svc.selectPosts(ctx, proto.SelectPostsRequest{CityId: cityId, StartTime: startTime, FinishTime: finishTime})
	if err != nil {
		return nil, nil, err
	}
	response := resp.(proto.SelectPostsReply)
	if response.Err != "" {
		return nil, nil, errors.New(response.Err)
	}
	return response.Posts, response.Area, nil
}

func (svc GrpcService) SelectAggrPosts(ctx context.Context, cityId string, interval data.SpatioHourInterval) ([]data.AggregatedPost, error) {
	resp, err := svc.selectAggrPosts(ctx, proto.SelectAggrPostsRequest{CityId: cityId, Interval: interval})
	if err != nil {
		return nil, err
	}
	response := resp.(proto.SelectAggrPostsReply)
	if response.Err != "" {
		return nil, errors.New(response.Err)
	}
	return response.Posts, nil
}

func (svc GrpcService) PullTimeline(ctx context.Context, cityId string, start, finish int64) ([]data.Timestamp, error) {
	resp, err := svc.pullTimeline(ctx, proto.PullTimelineRequest{CityId: cityId, Start: start, Finish: finish})
	if err != nil {
		return nil, err
	}
	response := resp.(proto.PullTimelineReply)
	if response.Err != "" {
		return nil, errors.New(response.Err)
	}
	return response.Timeline, nil
}

func (svc GrpcService) PushGrid(ctx context.Context, cityId string, ids []int64, blobs [][]byte) error {
	resp, err := svc.pushGrid(ctx, proto.PushGridRequest{CityId: cityId, Ids: ids, Blobs: blobs})
	if err != nil {
		return err
	}
	response := resp.(proto.PushGridReply)
	if response.Err != "" {
		return errors.New(response.Err)
	}
	return nil
}

func (svc GrpcService) PullGrid(ctx context.Context, cityId string, startId, finishId int64) ([]int64, [][]byte, error) {
	resp, err := svc.pullGrid(ctx, proto.PullGridRequest{CityId: cityId, StartId: startId, FinishId: finishId})
	if err != nil {
		return nil, nil, err
	}
	response := resp.(proto.PullGridReply)
	if response.Err != "" {
		return nil, nil, errors.New(response.Err)
	}
	return response.Ids, response.Blobs, nil
}

func (svc GrpcService) PushEvents(ctx context.Context, cityId string, events []data.Event) error {
	resp, err := svc.pushEvents(ctx, proto.PushEventsRequest{CityId: cityId, Events: events})
	if err != nil {
		return err
	}
	response := resp.(proto.PushEventsReply)
	if response.Err != "" {
		return errors.New(response.Err)
	}
	return nil
}

func (svc GrpcService) PullEvents(ctx context.Context, cityId string, interval data.SpatioHourInterval) ([]data.Event, error) {
	resp, err := svc.pullEvents(ctx, proto.PullEventsRequest{CityId: cityId, Interval: interval})
	if err != nil {
		return nil, err
	}
	response := resp.(proto.PullEventsReply)
	if response.Err != "" {
		return nil, errors.New(response.Err)
	}
	return response.Events, nil
}

func (svc GrpcService) PushLocations(ctx context.Context, cityId string, locations []data.Location) error {
	resp, err := svc.pushLocations(ctx, proto.PushLocationsRequest{CityId: cityId, Locations: locations})
	if err != nil {
		return err
	}
	response := resp.(proto.PushLocationsReply)
	if response.Err != "" {
		return errors.New(response.Err)
	}
	return nil
}

func (svc GrpcService) PullLocations(ctx context.Context, cityId string) ([]data.Location, error) {
	resp, err := svc.pullLocations(ctx, proto.PullLocationsRequest{CityId: cityId})
	if err != nil {
		return nil, err
	}
	response := resp.(proto.PullLocationsReply)
	if response.Err != "" {
		return nil, errors.New(response.Err)
	}
	return response.Locations, nil
}

func NewGRPCClient(conn *grpc.ClientConn) GrpcService {
	svc := GrpcService{}
	insertCityEndpoint := grpctransport.NewClient(
		conn, "proto.DataStorage", "InsertCity",
		encodeGRPCInsertCityRequest,
		decodeGRPCInsertCityResponse,
		proto.InsertCityReply{},
	).Endpoint()
	svc.insertCity = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "InsertCity",
		Timeout: TimeWaitingClient,
	}))(insertCityEndpoint)

	getAllCitiesEndpoint := grpctransport.NewClient(
		conn, "proto.DataStorage", "GetAllCities",
		encodeGRPCGetAllCitiesRequest,
		decodeGRPCGetAllCitiesResponse,
		proto.GetAllCitiesReply{},
	).Endpoint()
	svc.getAllCities = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "GetAllCities",
		Timeout: TimeWaitingClient,
	}))(getAllCitiesEndpoint)

	getCityEndpoint := grpctransport.NewClient(
		conn, "proto.DataStorage", "GetCity",
		encodeGRPCGetCityRequest,
		decodeGRPCGetCityResponse,
		proto.GetCityReply{},
	).Endpoint()
	svc.getCity = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "GetCity",
		Timeout: TimeWaitingClient,
	}))(getCityEndpoint)

	pushPostsEndpoint := grpctransport.NewClient(
		conn, "proto.DataStorage", "PushPosts",
		encodeGRPCPushPostsRequest,
		decodeGRPCPushPostsResponse,
		proto.PushPostsReply{},
	).Endpoint()
	svc.pushPosts = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "PushPosts",
		Timeout: TimeWaitingClient,
	}))(pushPostsEndpoint)

	selectPostsEndpoint := grpctransport.NewClient(
		conn, "proto.DataStorage", "SelectPosts",
		encodeGRPCSelectPostsRequest,
		decodeGRPCSelectPostsResponse,
		proto.SelectPostsReply{},
	).Endpoint()
	svc.selectPosts = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "SelectPosts",
		Timeout: TimeWaitingClient,
	}))(selectPostsEndpoint)

	selectAggrPostsEndpoint := grpctransport.NewClient(
		conn, "proto.DataStorage", "SelectAggrPosts",
		encodeGRPCSelectAggrPostsRequest,
		decodeGRPCSelectAggrPostsResponse,
		proto.SelectAggrPostsReply{},
	).Endpoint()
	svc.selectAggrPosts = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "SelectAggrPosts",
		Timeout: TimeWaitingClient,
	}))(selectAggrPostsEndpoint)

	pullTimelineEndpoint := grpctransport.NewClient(
		conn, "proto.DataStorage", "PullTimeline",
		encodeGRPCPullTimelineRequest,
		decodeGRPCPullTimelineResponse,
		proto.PullTimelineReply{},
	).Endpoint()
	svc.pullTimeline = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "PullTimeline",
		Timeout: TimeWaitingClient,
	}))(pullTimelineEndpoint)

	pushGridEndpoint := grpctransport.NewClient(
		conn, "proto.DataStorage", "PushGrid",
		encodeGRPCPushGridRequest,
		decodeGRPCPushGridResponse,
		proto.PushGridReply{},
	).Endpoint()
	svc.pushGrid = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "PushGrid",
		Timeout: TimeWaitingClient,
	}))(pushGridEndpoint)

	pullGridEndpoint := grpctransport.NewClient(
		conn, "proto.DataStorage", "PullGrid",
		encodeGRPCPullGridRequest,
		decodeGRPCPullGridResponse,
		proto.PullGridReply{},
	).Endpoint()
	svc.pullGrid = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "PullGrid",
		Timeout: TimeWaitingClient,
	}))(pullGridEndpoint)

	pushEventsEndpoint := grpctransport.NewClient(
		conn, "proto.DataStorage", "PushEvents",
		encodeGRPCPushEventsRequest,
		decodeGRPCPushEventsResponse,
		proto.PushEventsReply{},
	).Endpoint()
	svc.pushEvents = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "PushEvents",
		Timeout: TimeWaitingClient,
	}))(pushEventsEndpoint)

	pullEventsEndpoint := grpctransport.NewClient(
		conn, "proto.DataStorage", "PullEvents",
		encodeGRPCPullEventsRequest,
		decodeGRPCPullEventsResponse,
		proto.PullEventsReply{},
	).Endpoint()
	svc.pullEvents = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "PullEvents",
		Timeout: TimeWaitingClient,
	}))(pullEventsEndpoint)

	pushLocationsEndpoint := grpctransport.NewClient(
		conn, "proto.DataStorage", "PushLocations",
		encodeGRPCPushLocationsRequest,
		decodeGRPCPushLocationsResponse,
		proto.PushLocationsReply{},
	).Endpoint()
	svc.pushLocations = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "PushLocations",
		Timeout: TimeWaitingClient,
	}))(pushLocationsEndpoint)

	pullLocationsEndpoint := grpctransport.NewClient(
		conn, "proto.DataStorage", "PullLocations",
		encodeGRPCPullLocationsRequest,
		decodeGRPCPullLocationsResponse,
		proto.PullLocationsReply{},
	).Endpoint()
	svc.pullLocations = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "PullLocations",
		Timeout: TimeWaitingClient,
	}))(pullLocationsEndpoint)

	return svc
}
