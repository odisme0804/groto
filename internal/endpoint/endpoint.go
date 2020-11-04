package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"groto/pkg"
)

// Endpoints holds all Go kit endpoints for the Order service.
type Endpoints struct {
	SayHello endpoint.Endpoint
}

// MakeEndpoints initializes all Go kit endpoints for the Order service.
func MakeEndpoints(s pkg.GreeterServer) Endpoints {
	return Endpoints{
		SayHello: makeSayHello(s),
	}
}

func makeSayHello(s pkg.GreeterServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pkg.HelloRequest)
		reply, err := s.SayHello(ctx, req)
		if err == nil {
			return reply, nil
		}

		return nil, err
	}
}
