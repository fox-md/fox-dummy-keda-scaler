package scaler

import (
	"context"
	"strconv"
	"time"

	pb "github.com/fox-md/fox-dummy-keda-scaler/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FoxScaler struct {
	pb.UnimplementedExternalScalerServer
}

// type Quote struct {
// 	CatchPhrase string `json:"catchPhrase"`
// 	Name        string `json:"name"`
// }

var DummyCounter int = 1

func parseScaledObject(scaledObject *pb.ScaledObjectRef) (int, error) {
	capacityStr, ok := scaledObject.ScalerMetadata["dummy_capacity"]
	if !ok {
		return 0, status.Error(codes.InvalidArgument, "dummy_capacity must be specified")
	}

	capacity, err := strconv.Atoi(capacityStr)

	if err != nil {
		return 0, status.Error(codes.InvalidArgument, "Error during conversion:"+err.Error())
	}

	return capacity, nil
}

func (s *FoxScaler) IsActive(ctx context.Context, scaledObject *pb.ScaledObjectRef) (*pb.IsActiveResponse, error) {

	Logger.Debugw("IsActive called",
		"name", scaledObject.Name,
		"namespace", scaledObject.Namespace,
		"metadata", scaledObject.ScalerMetadata,
	)

	capacity, err := parseScaledObject(scaledObject)
	if err != nil {
		return nil, err
	}

	active := DummyCounter > capacity

	Logger.Debugw("IsActive result",
		"active", active,
	)

	return &pb.IsActiveResponse{Result: active}, nil
}

func (e *FoxScaler) StreamIsActive(scaledObject *pb.ScaledObjectRef, stream pb.ExternalScaler_StreamIsActiveServer) error {

	Logger.Infow("StreamIsActive started",
		"name", scaledObject.Name,
		"namespace", scaledObject.Namespace,
	)

	capacity, err := parseScaledObject(scaledObject)
	if err != nil {
		return err
	}

	active := false
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stream.Context().Done():
			// call cancelled
			Logger.Info("StreamIsActive closed")
			return nil
		case <-ticker.C:
			newActive := DummyCounter > capacity
			if newActive != active {
				Logger.Debugw("StreamIsActive result changed",
					"active", active,
					"newActive", newActive,
				)
				active = newActive
				stream.Send(&pb.IsActiveResponse{
					Result: active,
				})
			}
		}
	}
}

func (s *FoxScaler) GetMetricSpec(ctx context.Context, scaledObject *pb.ScaledObjectRef) (*pb.GetMetricSpecResponse, error) {

	Logger.Debugw("GetMetricSpec called",
		"metadata", scaledObject.ScalerMetadata,
	)

	capacityStr, ok := scaledObject.ScalerMetadata["dummy_capacity"]
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "missing dummy_capacity")
	}

	capacity, err := strconv.ParseFloat(capacityStr, 64)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid dummy_capacity")
	}

	return &pb.GetMetricSpecResponse{
		MetricSpecs: []*pb.MetricSpec{
			{
				// When the metric equals this value, the desired replica count is 1. These number of quotes can be handled by 1 replica.
				MetricName:      "dummy_capacity",
				TargetSizeFloat: capacity,
			},
		},
	}, nil
}

func (s *FoxScaler) GetMetrics(ctx context.Context, req *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {

	Logger.Debugw("GetMetrics called",
		"metric", req.MetricName,
		"name", req.ScaledObjectRef.Name,
		"metadata", req.ScaledObjectRef.ScalerMetadata,
	)

	return &pb.GetMetricsResponse{
		MetricValues: []*pb.MetricValue{
			{
				MetricName:  "dummy_capacity",
				MetricValue: int64(DummyCounter),
			},
		},
	}, nil
}
