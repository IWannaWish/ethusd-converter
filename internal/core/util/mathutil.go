package util

import (
	"math/big"
)

// ScaleDown делит value на 10^decimals, возвращая *big.Float
func ScaleDown(value *big.Int, decimals int) *big.Float {
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	return new(big.Float).Quo(new(big.Float).SetInt(value), new(big.Float).SetInt(scale))
}
