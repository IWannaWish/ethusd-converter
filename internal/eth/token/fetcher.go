package token

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type TokenBalanceFetcher interface {
	GetBalance(address common.Address) (*big.Float, error)
	GetSymbol() string
}
