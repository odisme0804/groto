package appserve

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type GracefulConfig struct {
	Delay   time.Duration `long:"delay" description:"delay shutdown after get signal" env:"DELAY" default:"10s"`
	Timeout time.Duration `long:"timeout" description:"graceful timeout" env:"TIMEOUT" default:"5s"`
}

type GraceStartFunc = func(context.Context) error

// Graceful passes a context to fn.
// Once a system interrupt signal is detected, it will cancel the context.
// The fn should return error after the context is canceled
func Graceful(startFunc GraceStartFunc, config GracefulConfig) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error, 1)

	go func() {
		done <- startFunc(ctx)
	}()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	var err error
	select {
	case err = <-done:
		return err
	case <-c:
		// since there is a race window between k8s removing endpoint and
		// be known by loadbalancer, so we have no choice but wait for a time.
		time.Sleep(config.Delay)
		cancel()

		ctx, timeout := context.WithTimeout(ctx, config.Timeout)
		defer timeout()

		select {
		case err = <-done:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
