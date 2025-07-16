package core

import (
	"github.com/ethereum/go-ethereum/common"
)

type Asset struct {
	Symbol   string
	Balance  string // например "1.230000"
	USDValue string // например "$4,312.45"
}

type AssetService interface {
	GetAssets(address common.Address) ([]Asset, error)
}
