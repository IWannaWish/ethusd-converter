package main

import (
	"context"
	"github.com/IWannaWish/ethusd-converter/cmd/cli/display"
	"github.com/IWannaWish/ethusd-converter/internal/core/mapper"
	"github.com/IWannaWish/ethusd-converter/internal/eth/abi"
	"github.com/IWannaWish/ethusd-converter/internal/eth/source"
	"github.com/IWannaWish/ethusd-converter/internal/log"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"os"

	"github.com/IWannaWish/ethusd-converter/internal/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	_ = godotenv.Load("config.env")

	cfg := config.Load()
	logger := log.NewLogger(cfg)
	ctx := log.WithRequestID(context.Background(), uuid.NewString())

	logger.Info(ctx, "ethusd-converter started",
		log.String("log_level", cfg.LogLevel),
		log.String("module", "main"),
	)

	if len(os.Args) < 2 {
		logger.Error(ctx, "Usage: ./ethusd-converter <ethereum_address>")
		os.Exit(1)
	}

	rawAddr := os.Args[1]
	if !common.IsHexAddress(rawAddr) {
		logger.Error(ctx, "Неверный Ethereum адрес", log.String("address", rawAddr))
		os.Exit(1)
	}
	address := common.HexToAddress(rawAddr)

	client, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		logger.Error(ctx, "Ошибка подключения к Ethereum node", log.WithStack(err)...)
		os.Exit(1)
	}
	defer client.Close()

	erc20ABI, err := abi.LoadERC20ABI()
	if err != nil {
		logger.Error(ctx, "Ошибка загрузки ERC20 ABI", log.WithStack(err)...)
		os.Exit(1)
	}

	feedABI, err := abi.LoadAggregatorABI()
	if err != nil {
		logger.Error(ctx, "Ошибка загрузки Chainlink ABI", log.WithStack(err)...)
		os.Exit(1)
	}

	sources, err := source.BuildAssetSources(cfg.Tokens, client, erc20ABI, feedABI)
	if err != nil {
		logger.Error(ctx, "Ошибка построения токенов и фидов", log.WithStack(err)...)
		os.Exit(1)
	}

	service := source.NewAssetService(sources)

	assets, err := service.GetAssets(ctx, address)
	if err != nil {
		logger.Error(ctx, "Ошибка получения активов", log.WithStack(err)...)
		os.Exit(1)
	}

	mapper := mapper.NewDisplayAssetMapper()
	printer := display.NewTablePrinter()

	info, total, err := mapper.Map(assets)
	if err != nil {
		logger.Error(ctx, "Ошибка преобразования активов", log.WithStack(err)...)
		os.Exit(1)
	}

	printer.Print(address.Hex(), info, total)
}
