package scaler

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var Logger *zap.SugaredLogger

func LoggingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {

	start := time.Now()

	Logger.Debugw("grpc request",
		"method", info.FullMethod,
		"request", req,
	)

	resp, err := handler(ctx, req)

	if err != nil {
		Logger.Errorw("grpc error",
			"method", info.FullMethod,
			"error", err,
			"duration", time.Since(start),
		)
		return resp, err
	}

	Logger.Debugw("grpc response",
		"method", info.FullMethod,
		"response", resp,
		"duration", time.Since(start),
	)

	return resp, nil
}
