package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fox-md/fox-dummy-keda-scaler/scaler"
	"github.com/fox-md/fox-dummy-keda-scaler/server"
	"go.uber.org/zap"
)

func main() {

	gRPCport := flag.Int("grpc-port", 6000, "gRPC port")
	httpPort := flag.Int("http-port", 8080, "http port")
	debug := flag.Bool("debug", false, "enable debug logging")

	flag.Parse()

	level := zap.InfoLevel

	if *debug {
		level = zap.DebugLevel
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(level)

	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	scaler.Logger = l.Sugar()
	defer l.Sync()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	go server.StartGRPC(ctx, fmt.Sprintf(":%d", *gRPCport))
	go server.StartHTTP(ctx, fmt.Sprintf(":%d", *httpPort))

	<-ctx.Done()
	log.Print("shutting down")
}
