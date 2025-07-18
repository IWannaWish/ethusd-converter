package mapper

import (
	"math/big"
	"strings"
)

func ParseBalance(raw string) (*big.Float, error) {
	f, _, err := big.ParseFloat(raw, 10, 256, big.ToNearestEven)
	return f, err
}

func ParseUSDValue(raw string) (*big.Float, error) {
	clean := strings.TrimPrefix(raw, "$")
	f, _, err := big.ParseFloat(clean, 10, 256, big.ToNearestEven)
	return f, err
}
