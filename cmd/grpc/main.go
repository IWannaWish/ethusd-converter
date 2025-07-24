package main

import (
	"context"
	"log"
	"net"

	apigrpc "github.com/IWannaWish/ethusd-converter/internal/api/grpc"
	"github.com/IWannaWish/ethusd-converter/internal/applog"
	"github.com/IWannaWish/ethusd-converter/internal/config"
	"github.com/IWannaWish/ethusd-converter/internal/eth/abi"
	"github.com/IWannaWish/ethusd-converter/internal/eth/source"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load("config.env")
	cfg := config.Load()

	logger := applog.NewLogger(cfg)
	ctx := applog.WithRequestID(context.Background(), uuid.NewString())

	logger.Info(ctx, "Starting gRPC server...", applog.String("rpc_url", cfg.RPCURL))

	// Подключаемся к Ethereum
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

	sources, err := source.BuildAssetSources(ctx, logger, cfg.Tokens, client, erc20ABI, feedABI)
	if err != nil {
		logger.Error(ctx, "Ошибка построения источников", applog.WithStack(err)...)
		return
	}

	assetService := source.NewAssetService(sources, logger)
	grpcServer := grpc.NewServer()

	server := apigrpc.NewEthusdGRPCServer(assetService)

	// Регистрация gRPC сервера
	apigrpc.RegisterEthusdConverterServer(grpcServer, server)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Error(ctx, "Ошибка запуска TCP-листенера", applog.WithStack(err)...)
		return
	}

	logger.Info(ctx, "gRPC сервер слушает на порту :50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
}
