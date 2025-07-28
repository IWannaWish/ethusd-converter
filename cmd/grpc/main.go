package main

import (
	"context"
	"google.golang.org/grpc/reflection"
	"net"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"

	"github.com/IWannaWish/ethusd-converter/internal/applog"
	"github.com/IWannaWish/ethusd-converter/internal/config"
	"github.com/IWannaWish/ethusd-converter/internal/eth/abi"
	"github.com/IWannaWish/ethusd-converter/internal/eth/source"

	"github.com/IWannaWish/ethusd-converter/internal/api/grpc_server"
	ethusdpb "github.com/IWannaWish/ethusd-converter/proto/ethusd/gen"

	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load("config.env")

	cfg := config.Load()
	logger := applog.NewLogger(cfg)
	ctx := context.Background()

	logger.Info(ctx, "gRPC-сервер запускается...",
		applog.String("module", "grpc_server"),
		applog.String("rpc_url", cfg.RPCURL),
	)

	client, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		logger.Error(ctx, "Ошибка подключения к Ethereum", applog.WithStack(err)...)
		return
	}
	defer client.Close()

	erc20ABI, err := abi.LoadERC20ABI()
	if err != nil {
		logger.Error(ctx, "Ошибка загрузки ERC20 ABI", applog.WithStack(err)...)
		return
	}

	feedABI, err := abi.LoadAggregatorABI()
	if err != nil {
		logger.Error(ctx, "Ошибка загрузки Chainlink ABI", applog.WithStack(err)...)
		return
	}

	sources, err := source.BuildAssetSources(ctx, logger, cfg.Tokens, client, cfg, erc20ABI, feedABI)
	if err != nil {
		logger.Error(ctx, "Ошибка построения источников активов", applog.WithStack(err)...)
		return
	}

	assetService := source.NewAssetService(sources, logger)

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_server.RequestIDInterceptor(logger)))
	server := grpc_server.NewEthusdGRPCServer(assetService)
	ethusdpb.RegisterEthusdConverterServer(grpcServer, server)
	// Включаем рефлексию — нужно для grpcurl
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Error(ctx, "Ошибка запуска TCP-листенера", applog.WithStack(err)...)
		return
	}

	logger.Info(ctx, "gRPC сервер слушает на :50051")
	if err := grpcServer.Serve(listener); err != nil {
		logger.Error(ctx, "Ошибка запуска gRPC сервера", applog.WithStack(err)...)
		os.Exit(1)
	}
}
