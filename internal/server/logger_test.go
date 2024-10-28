package server

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Mock request type
type mockRequest struct{}

// Mock handler to simulate a successful gRPC call
func successfulHandler(ctx context.Context, req any) (any, error) {
	return "success", nil
}

// Mock handler to simulate a gRPC call that returns an error
func errorHandler(ctx context.Context, req any) (any, error) {
	return nil, status.Error(codes.Unknown, "unknown error")
}

func TestFetchLoggingInterceptor_Success(t *testing.T) {
	// Set up the logger for testing
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.NewConsoleWriter())

	// Create the interceptor
	interceptor := FetchLoggingInterceptor()

	// Create a new context
	ctx := context.Background()

	// Create UnaryServerInfo
	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.MockService/MockMethod",
	}

	// Call the interceptor with a successful handler
	_, err := interceptor(ctx, &mockRequest{}, info, successfulHandler)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestFetchLoggingInterceptor_Error(t *testing.T) {
	// Set up the logger for testing
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.NewConsoleWriter())

	// Create the interceptor
	interceptor := FetchLoggingInterceptor()

	// Create a new context
	ctx := context.Background()

	// Create UnaryServerInfo
	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.MockService/MockMethod",
	}

	// Call the interceptor with an error handler
	_, err := interceptor(ctx, &mockRequest{}, info, errorHandler)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}

	// Check if the error message is as expected
	expectedErrorMessage := "unknown error"
	if status.Code(err) != codes.Unknown || status.Convert(err).Message() != expectedErrorMessage {
		t.Fatalf("Expected error message '%s', got '%s'", expectedErrorMessage, status.Convert(err).Message())
	}
}
