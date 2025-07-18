package display

import "math/big"

type AssetInfo struct {
	Symbol   string
	Balance  *big.Float
	USDPrice *big.Float
	USDValue *big.Float
}

type AssetPrinter interface {
	Print(address string, assets []AssetInfo, totalUSD *big.Float)
}
