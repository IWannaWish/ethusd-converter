package pricestore

import (
	"golang.org/x/net/context"
	"math/big"
)

type PriceStore interface {
	Get(symbol string) (*big.Float, bool)
	StartBackgroundAdapter(ctx context.Context) error
}
