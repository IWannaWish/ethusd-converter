package core

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// TokenBalanceFetcher — универсальный интерфейс для получения баланса токена (ETH, ERC-20)
type TokenBalanceFetcher interface {
	GetBalance(ctx context.Context, holder common.Address) (*big.Float, error)
	GetSymbol() string
}
