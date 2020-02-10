package service

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func makeNewSessionEndpoint(svc CoordinatorService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		p := request.(SessionParameters)
		res, err := svc.NewSession(p)
		return res, err
	}
}

func makeStatusEndpoint(svc CoordinatorService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		id := request.(string)
		res, err := svc.Status(id)
		return res, err
	}
}
