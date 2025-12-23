package server

import (
	"context"
	"net"

	"log"

	pb "github.com/fox-md/fox-dummy-keda-scaler/proto"
	"github.com/fox-md/fox-dummy-keda-scaler/scaler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGRPC(ctx context.Context, address string) {

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(scaler.LoggingInterceptor),
	)
	pb.RegisterExternalScalerServer(grpcServer, &scaler.FoxScaler{})

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	log.Printf("KEDA external fox scaler listening on %s", address)
	log.Fatal(grpcServer.Serve(lis))

	go func() {
		<-ctx.Done()
		log.Print("stopping grpc server")
		grpcServer.GracefulStop()
	}()

	log.Printf("grpc listening on %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("grpc serve: %v", err)
	}
}
