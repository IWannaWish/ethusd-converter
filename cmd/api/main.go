package main

import (
	"github.com/TimRutte/api/internal/config"
	"github.com/TimRutte/api/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func main() {
	// Print application version information
	Version()

	// Load the configuration
	cfg, err := config.New("config.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading configuration")
	}

	// Initialize zero log
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	// Initialize the gRPC server
	srv := server.New(cfg)

	// Start the gRPC server
	if cfg.GrpcServer.Enabled {
		lis, grpc := srv.InitGrpcServer()
		log.Info().Msgf("gRPC server starting on %s:%d", cfg.GrpcServer.Host, cfg.GrpcServer.Port)
		go srv.StartGrpcServer(lis, grpc)
	} else {
		log.Info().Msg("gRPC server is disabled")
	}

	// Start the HTTP server
	if cfg.HttpServer.Enabled {
		http := srv.InitHttpServer()
		log.Info().Msgf("HTTP server starting on %s:%d", cfg.HttpServer.Host, cfg.HttpServer.Port)
		srv.StartHttpServer(http)
		err = http.ListenAndServe()
		if err != nil {
			log.Fatal().Err(err).Msg("Error starting HTTP server")
		}
		log.Info().Msgf("HTTP server started on %s:%d", cfg.HttpServer.Host, cfg.HttpServer.Port)
	} else {
		log.Info().Msg("HTTP server is disabled")
	}
}
