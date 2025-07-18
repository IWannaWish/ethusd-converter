package util

import (
	"math/big"
	"testing"
)

func TestScaleDown(t *testing.T) {
	tests := []struct {
		name     string
		value    string // целочисленный input (wei)
		decimals int
		expected string // ожидаемое *big.Float значение
	}{
		{"1 ETH", "1000000000000000000", 18, "1"},
		{"0.5 ETH", "500000000000000000", 18, "0.5"},
		{"123.456 USDC", "123456000", 6, "123.456"},
		{"0", "0", 18, "0"},
		{"Round down", "999", 3, "0.999"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valueInt := new(big.Int)
			valueInt.SetString(tt.value, 10)

			result := ScaleDown(valueInt, tt.decimals)
			expected := new(big.Float)
			expected.SetString(tt.expected)

			if result.Cmp(expected) != 0 {
				t.Errorf("expected %s, got %s", expected.Text('f', -1), result.Text('f', -1))
			}
		})
	}
}
