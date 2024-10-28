package server

import (
	"testing"
	"time"

	"github.com/TimRutte/api/internal/config"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

// MockConfig simulates the configuration.
var mockConfig = config.Config{
	GrpcServer: config.GrpcServerConfig{
		Host: "localhost",
		Port: 50053,
		TLS:  false,
	},
	HttpServer: config.HttpServerConfig{
		Host: "localhost",
		Port: 8079,
	},
}

func TestInitGrpcServer(t *testing.T) {
	server := New(mockConfig)

	lis, grpcServer := server.InitGrpcServer()
	assert.NotNil(t, grpcServer)
	assert.NotNil(t, lis)

	// Clean up by closing the listener
	defer lis.Close()
}

func TestStartGrpcServer(t *testing.T) {
	server := New(mockConfig)

	lis, grpcServer := server.InitGrpcServer()
	assert.NotNil(t, grpcServer)

	// Start a goroutine to run the gRPC server
	go func() {
		err := grpcServer.Serve(lis)
		assert.NoError(t, err, "gRPC server should start without error")
	}()

	// Allow some time for the server to start
	time.Sleep(100 * time.Millisecond)

	// Test if we can dial the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	// Clean up
	conn.Close()
	grpcServer.Stop()
	lis.Close()
}
