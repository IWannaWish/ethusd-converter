package main

import (
	"context"
	"google.golang.org/grpc/reflection"
	"log"
	"net"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	applog "github.com/IWannaWish/ethusd-converter/internal/applog"
	"github.com/IWannaWish/ethusd-converter/internal/config"
	"github.com/IWannaWish/ethusd-converter/internal/eth/abi"
	"github.com/IWannaWish/ethusd-converter/internal/eth/source"

	apigrpc "github.com/IWannaWish/ethusd-converter/internal/api/grpc"
	ethusd_pb "github.com/IWannaWish/ethusd-converter/proto/ethusd/gen"

	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load("config.env")

	cfg := config.Load()
	logger := applog.NewLogger(cfg)
	ctx := applog.WithRequestID(context.Background(), uuid.NewString())

	logger.Info(ctx, "gRPC-сервер запускается...",
		applog.String("module", "grpc"),
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

	sources, err := source.BuildAssetSources(ctx, logger, cfg.Tokens, client, erc20ABI, feedABI)
	if err != nil {
		logger.Error(ctx, "Ошибка построения источников активов", applog.WithStack(err)...)
		return
	}

	assetService := source.NewAssetService(sources, logger)

	grpcServer := grpc.NewServer()
	server := apigrpc.NewEthusdGRPCServer(assetService)
	ethusd_pb.RegisterEthusdConverterServer(grpcServer, server)
	// Включаем рефлексию — нужно для grpcurl
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Error(ctx, "Ошибка запуска TCP-листенера", applog.WithStack(err)...)
		return
	}

	logger.Info(ctx, "gRPC сервер слушает на :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка запуска gRPC сервера: %v", err)
	}
}
