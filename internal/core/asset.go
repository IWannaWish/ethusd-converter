package core

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
)

// Asset представляет токен на аккаунте в виде пользовательского вывода
type Asset struct {
	Symbol   string // Символ токена, например "ETH", "DAI"
	Balance  string // Баланс токена в читаемом виде, например "1.230000"
	USDValue string // Оценка в долларах, например "$4,312.45"
}

type AssetService interface {
	GetAssets(ctx context.Context, address common.Address) ([]Asset, error)
}
