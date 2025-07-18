package core

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type TokenBalanceFetcher interface {
	GetBalance(ctx context.Context, holder common.Address) (*big.Float, error)
	GetSymbol() string
}
