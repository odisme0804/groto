package service

import (
	"context"
	"groto/pkg"
)

type server struct {
	// pkg.UnimplementedGreeterServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) SayHello(context.Context, *pkg.HelloRequest) (*pkg.HelloReply, error) {
	return &pkg.HelloReply{
		Message: "hello go-kit",
	}, nil
}

var _ pkg.GreeterServer = (*server)(nil)
