package server

import (
	"context"
	"fmt"
	"github.com/TimRutte/api/internal/config"
	"github.com/TimRutte/api/proto/api"
	"github.com/TimRutte/api/proto/healthcheck"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
)

type Server struct {
	Config config.Config
}

func New(cfg config.Config) *Server {
	return &Server{
		Config: cfg,
	}
}

func (s *Server) InitGrpcServer(grpcOptions ...grpc.ServerOption) (net.Listener, *grpc.Server) {
	// Setup server but don't listen
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Config.GrpcServer.Host, s.Config.GrpcServer.Port))
	if grpcOptions == nil {
		grpcOptions = make([]grpc.ServerOption, 0)
	}

	if err != nil {
		log.Fatal().Err(err).Msgf("Could not open listener on `%s:%d`", s.Config.GrpcServer.Host, s.Config.GrpcServer.Port)
	}

	if s.Config.GrpcServer.TLS {
		cert, err := credentials.NewServerTLSFromFile(s.Config.GrpcServer.CertFile, s.Config.GrpcServer.KeyFile)
		if err != nil {
			log.Logger.Fatal().Err(err).Msg("Error setting up TLS")
		}
		grpcOptions = append(grpcOptions, grpc.Creds(cert))
	}

	grpcOptions = append(grpcOptions, grpc.ChainUnaryInterceptor(
		FetchLoggingInterceptor(),
	))

	grpcServer := grpc.NewServer(grpcOptions...)

	// Register services
	s.registerServices(grpcServer)

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	return lis, grpcServer
}

func (s *Server) StartGrpcServer(lis net.Listener, server *grpc.Server) {
	// Listen for incoming gRPC connections
	err := server.Serve(lis)
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("Failed to start gRPC server")
	}
}

func (s *Server) InitHttpServer() *http.Server {
	ctx := context.Background()
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := healthcheck.RegisterHealthHandlerFromEndpoint(ctx, mux, fmt.Sprintf("%s:%d", s.Config.GrpcServer.Host, s.Config.GrpcServer.Port), opts); err != nil {
		log.Fatal().Err(err).Msg("Failed to register HealthCheck HTTP gateway")
	}

	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.Config.HttpServer.Host, s.Config.HttpServer.Port),
		Handler: mux,
	}
}

func (s *Server) StartHttpServer(httpServer *http.Server) {
	// Listen for incoming HTTP connections
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start HTTP server")
	}
}

func (s *Server) registerServices(grpcServer *grpc.Server) {
	healthcheck.RegisterHealthServer(grpcServer, HealthCheck{})
	api.RegisterApiServer(grpcServer, Api{})
}
