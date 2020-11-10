package main

import (
	"context"
	"fmt"
	"groto/pkg"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jessevdk/go-flags"
	"google.golang.org/grpc"
)

type Args struct {
	GRPCAddr string `long:"grpc.addr" env:"GRPC_ADDR" default:":8081"`
}

func main() {

	args := Args{}
	_, err := flags.NewParser(&args, flags.Default).Parse()
	if err != nil {
		panic(err)
	}
	fmt.Println("args", args)

	logger := log.NewLogfmtLogger(os.Stderr)

	conn, err := grpc.Dial(args.GRPCAddr,
		grpc.WithInsecure(), // disable ssl
		grpc.WithBlock(),    // wait the dial
	)
	if err != nil {
		level.Error(logger).Log("grpc.Dial", err)
	}
	defer conn.Close()

	client := pkg.NewGreeterClient(conn)

	res, err := client.SayHello(context.Background(), &pkg.HelloRequest{})
	if err != nil {
		level.Error(logger).Log("client.SayHello", err)
	}

	logger.Log("res:", res)
}
