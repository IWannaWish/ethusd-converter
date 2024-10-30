package api

import (
	"context"
	"testing"

	healthpb "github.com/TimRutte/api/proto/healthcheck/gen"
)

func TestHealthCheck_Check(t *testing.T) {
	// Initialize the HealthCheck server
	grpcHealth := HealthCheck{}

	// Create a context and an empty HealthCheckRequest
	ctx := context.Background()
	req := &healthpb.HealthCheckRequest{}

	// Call the Check method
	resp, err := grpcHealth.Check(ctx, req)

	// Check for errors
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the response status is SERVING
	if resp.Status != healthpb.HealthCheckResponse_SERVING {
		t.Errorf("Expected status %v, got %v", healthpb.HealthCheckResponse_SERVING, resp.Status)
	}
}
