package main

import (
	"github.com/IWannaWish/ethusd-converter/internal/core"
	"github.com/IWannaWish/ethusd-converter/internal/eth"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Загружаем .env переменные
	if err := godotenv.Load("config.env"); err != nil {
		log.Println(".env файл не найден, переменные должны быть экспортированы вручную")
	}

	// 2. Читаем переменные из окружения
	rpcURL := os.Getenv("RPC_URL")
	feedAddrStr := os.Getenv("CHAINLINK_ETH_USD")

	if rpcURL == "" || feedAddrStr == "" {
		log.Fatal("RPC_URL и CHAINLINK_ETH_USD должны быть заданы в .env или окружении")
	}

	// 3. Читаем адрес из аргументов
	if len(os.Args) < 2 {
		log.Fatal("Usage: ./ethusd-converter <ethereum_address>")
	}
	rawAddr := os.Args[1]
	if !common.IsHexAddress(rawAddr) {
		log.Fatalf("Invalid Ethereum address: %s", rawAddr)
	}
	address := common.HexToAddress(rawAddr)

	// 4. Подключаемся к Ethereum
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum node: %v", err)
	}
	defer client.Close()

	// 5. Загружаем ABI
	aggregatorABI, err := eth.LoadAggregatorABI()
	if err != nil {
		log.Fatalf("Failed to load Chainlink ABI: %v", err)
	}

	// 6. Собираем зависимость feed (через интерфейс)
	feedAddress := common.HexToAddress(feedAddrStr)
	feed := eth.NewChainlinkFeed(client, feedAddress, aggregatorABI)

	// 7. Собираем бизнес-сервис (через интерфейс)
	service := core.NewAssetService(client, feed)

	// 8. Получаем активы
	assets, err := service.GetAssets(address)
	if err != nil {
		log.Fatalf("Failed to get assets: %v", err)
	}

	// 9. Выводим красиво
	log.Printf("Address: %s", rawAddr)
	var total string
	for _, asset := range assets {
		log.Printf("%s: %s ≈ %s", asset.Symbol, asset.Balance, asset.USDValue)
		total = asset.USDValue // позже заменим на суммирование
	}
	log.Printf("Total: %s", total)
}
