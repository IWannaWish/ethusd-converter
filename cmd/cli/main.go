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
	_ = godotenv.Load("config.env") // мягко, без фатала

	// 1. Загружаем конфигурацию
	cfg := config.Load()

	logger := log.NewSlogLogger(cfg)
	ctx := log.WithRequestID(context.Background(), uuid.NewString())

	logger.Info(ctx, "ethusd-converter started",
		log.String("log_format", cfg.LogFormat),
		log.String("log_level", cfg.LogLevel),
		log.String("module", "main"),
	)

	// 2. Читаем аргумент адреса
	if len(os.Args) < 2 {
		//log.Fatal("Usage: ./ethusd-converter <ethereum_address>")
	}
	rawAddr := os.Args[1]
	if !common.IsHexAddress(rawAddr) {
		//log.Fatalf("Неверный Ethereum адрес: %s", rawAddr)
	}
	address := common.HexToAddress(rawAddr)

	// 3. Подключение к Ethereum
	client, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		//log.Fatalf("Ошибка подключения к Ethereum node: %v", err)
	}
	defer client.Close()

	// 4. ABI
	erc20ABI, err := abi.LoadERC20ABI()
	if err != nil {
		//log.Fatalf("Ошибка загрузки ERC20 ABI: %v", err)
	}
	feedABI, err := abi.LoadAggregatorABI()
	if err != nil {
		//log.Fatalf("Ошибка загрузки Chainlink ABI: %v", err)
	}

	// 5. Источники данных
	sources, err := source.BuildAssetSources(cfg.Tokens, client, erc20ABI, feedABI)
	if err != nil {
		//log.Fatalf("Ошибка построения токенов и фидов: %v", err)
	}

	// 6. Бизнес-логика
	service := source.NewAssetService(sources)

	assets, err := service.GetAssets(context.Background(), address)
	if err != nil {
		//log.Fatalf("Ошибка получения активов: %v", err)
	}

	// 7. Вывод
	simpleMapper := mapper.NewDisplayAssetMapper()
	printer := display.NewTablePrinter()

	info, total, err := simpleMapper.Map(assets)
	if err != nil {
		//log.Fatalf("Ошибка преобразования активов: %v", err)
	}
	printer.Print(address.Hex(), info, total)
}
