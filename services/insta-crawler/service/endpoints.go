package service

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler"
)

func makeNewEndpoint(svc CrawlerService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		p := request.(crawler.Parameters)
		id, err := svc.New(p)
		if err != nil {
			return newEpResponse{"", err.Error()}, nil
		}
		return newEpResponse{id, ""}, nil
	}
}

func makeStatusEndpoint(svc CrawlerService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(idEpRequest)
		s, err := svc.Status(req.ID)
		if err != nil {
			return statusEpResponse{crawler.OutStatus{}, err.Error()}, nil
		}
		return statusEpResponse{s, ""}, nil
	}
}

func makeStopEndpoint(svc CrawlerService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(idEpRequest)
		ok, err := svc.Stop(req.ID)
		if err != nil {
			return stopEpResponse{false, err.Error()}, nil
		}
		return stopEpResponse{ok, ""}, nil
	}
}

func makeEntitiesEndpoint(svc CrawlerService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(idEpRequest)
		ents, err := svc.Entities(req.ID)
		if err != nil {
			return entitiesEpResponse{nil, err.Error()}, nil
		}
		return entitiesEpResponse{ents, ""}, nil
	}
}

func makePostsEndpoint(svc CrawlerService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(postsEpRequest)
		posts, cursor, err := svc.Posts(req.ID, req.Offset, req.Num)
		if err != nil {
			return postsEpResponse{nil, "", err.Error()}, nil
		}
		return postsEpResponse{posts, cursor, ""}, nil
	}
}
