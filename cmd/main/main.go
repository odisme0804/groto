package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/jessevdk/go-flags"
	"google.golang.org/grpc"

	"groto/internal/endpoint"
	"groto/internal/service"
	"groto/internal/transport"
	"groto/lib/appserve"
	"groto/pkg"
)

type Args struct {
	HTTPAddr string `long:"http.addr" env:"HTTP_ADDR" default:":8080"`
	GRPCAddr string `long:"grpc.addr" env:"GRPC_ADDR" default:":8081"`
	appserve.GracefulConfig
}

func main() {
	args := Args{}
	_, err := flags.NewParser(&args, flags.Default).Parse()
	if err != nil {
		panic(err)
	}

	logger := log.NewLogfmtLogger(os.Stderr)

	logger.Log("server start, args:", args)

	err = appserve.Graceful(genGraceStartFunc(logger, args), args.GracefulConfig)

	logger.Log("server shutdown, err:", err)
}

func genGraceStartFunc(logger log.Logger, args Args) appserve.GraceStartFunc {
	s := service.NewServer()

	e := endpoint.MakeEndpoints(s)

	t := transport.NewHTTPServer(e, logger)

	// setup http transport
	httpistener, err := net.Listen("tcp", args.HTTPAddr)
	if err != nil {
		panic(err)
	}

	httpServer := http.Server{
		Handler: t,
	}

	// setup grpc transport
	grpcListener, err := net.Listen("tcp", args.GRPCAddr)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pkg.RegisterGreeterServer(grpcServer, transport.NewGRPCServer(e, logger))

	return func(ctx context.Context) error {

		go func() {
			if err := httpServer.Serve(httpistener); err != nil {
				panic(err)
			}
		}()

		go func() {
			if err := grpcServer.Serve(grpcListener); err != nil {
				panic(err)
			}
		}()

		<-ctx.Done()

		var result error
		var wg sync.WaitGroup

		go func() {
			wg.Add(1)
			defer wg.Done()

			if err := httpServer.Shutdown(context.Background()); err != nil && !errors.Is(err, http.ErrServerClosed) {
				result = err
			}

			httpistener.Close()
		}()

		go func() {
			wg.Add(1)
			defer wg.Done()

			grpcServer.GracefulStop()
			grpcListener.Close()
		}()

		wg.Wait()

		return result
	}
}
