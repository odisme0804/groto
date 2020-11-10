package transport

import (
	"context"
	"groto/internal/endpoint"
	"groto/pkg"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	"github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	SayHelloHandler grpc.Handler
}

func (s *grpcServer) SayHello(ctx context.Context, req *pkg.HelloRequest) (*pkg.HelloReply, error) {
	_, resp, err := s.SayHelloHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pkg.HelloReply), nil
}

func NewGRPCServer(e endpoint.Endpoints, logger log.Logger) pkg.GreeterServer {
	options := []grpc.ServerOption{
		grpc.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &grpcServer{
		SayHelloHandler: grpc.NewServer(
			e.SayHello,
			decodeSayHelloGRPCRequest,
			encodeSayHelloRPCResponse,
			options...,
		),
	}
}
func decodeSayHelloGRPCRequest(_ context.Context, in interface{}) (interface{}, error) {
	// already decode when grpc get the msg
	// this is just for the go-kit pattern
	// just pass through the args

	// no need
	// req := in.(*pkg.HelloRequest)
	// return req, nil

	return in, nil
}
func encodeSayHelloRPCResponse(_ context.Context, in interface{}) (interface{}, error) {
	return in, nil
}
