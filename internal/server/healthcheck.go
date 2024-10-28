package server

import (
	"context"
	healthpb "github.com/TimRutte/api/proto/healthcheck/gen"
)

type HealthCheck struct {
	healthpb.UnimplementedHealthServer
}

func (grpcHealth HealthCheck) Check(ctx context.Context, in *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{
		Status: healthpb.HealthCheckResponse_SERVING,
	}, nil
}
