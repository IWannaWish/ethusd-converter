package main

import (
	"context"
	"fmt"
	"github.com/IWannaWish/ethusd-converter/internal/eth"
	"log"
	"os"

	"github.com/IWannaWish/ethusd-converter/internal/config"
	"github.com/IWannaWish/ethusd-converter/internal/core"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Загружаем .env файл
	_ = godotenv.Load("config.env") // мягко, без фатала

	// 2. Читаем переменные окружения
	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		log.Fatal("RPC_URL не задан в .env или окружении")
	}

	// 3. Читаем аргументы
	if len(os.Args) < 2 {
		log.Fatal("Usage: ./ethusd-converter <ethereum_address>")
	}
	rawAddr := os.Args[1]
	if !common.IsHexAddress(rawAddr) {
		log.Fatalf("Неверный Ethereum адрес: %s", rawAddr)
	}
	address := common.HexToAddress(rawAddr)

	// 4. Подключаемся к Ethereum node
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к Ethereum node: %v", err)
	}
	defer client.Close()

	// 5. Загружаем tokenlist.yaml
	tokenList, err := config.LoadTokenList("internal/config/tokenlist.yaml")
	if err != nil {
		log.Fatalf("Не удалось загрузить tokenlist.yaml: %v", err)
	}

	// 6. Загружаем ABI
	erc20ABI, err := eth.LoadERC20ABI()
	if err != nil {
		log.Fatalf("Ошибка загрузки ERC20 ABI: %v", err)
	}
	feedABI, err := eth.LoadAggregatorABI()
	if err != nil {
		log.Fatalf("Ошибка загрузки Chainlink ABI: %v", err)
	}

	// 7. Строим источники данных (токены + прайс фиды)
	sources, err := core.BuildAssetSources(tokenList.Tokens, client, erc20ABI, feedABI)
	if err != nil {
		log.Fatalf("Ошибка построения токенов и фидов: %v", err)
	}

	// 8. Инициализируем бизнес-сервис
	service := core.NewAssetService(sources)

	// 9. Получаем активы
	assets, err := service.GetAssets(context.Background(), address)
	if err != nil {
		log.Fatalf("Ошибка получения активов: %v", err)
	}

	// 10. Печатаем результат
	fmt.Printf("Address: %s\n", address.Hex())
	for _, asset := range assets[:len(assets)-1] {
		fmt.Printf("%s: %s ≈ %s\n", asset.Symbol, asset.Balance, asset.USDValue)
	}
	total := assets[len(assets)-1] // последний — Total
	fmt.Printf("\nTotal: %s\n", total.USDValue)
}
