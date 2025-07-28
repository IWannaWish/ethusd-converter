package pricestore

import (
	"golang.org/x/net/context"
	"math/big"
)

type PriceStore interface {
	Get(ctx context.Context, symbol string) (*big.Float, bool)
	StartBackgroundUpdater(ctx context.Context) error
}
