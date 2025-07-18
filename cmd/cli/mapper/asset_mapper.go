package mapper

import (
	"fmt"
	"github.com/IWannaWish/ethusd-converter/cmd/cli/display"
	"github.com/IWannaWish/ethusd-converter/internal/core"
	"math/big"
	"strings"
)

type AssetMapper interface {
	Map(coreAssets []core.Asset) ([]display.AssetInfo, *big.Float, error)
}

type SimpleAssetMapper struct{}

func NewSimpleAssetMapper() *SimpleAssetMapper {
	return &SimpleAssetMapper{}
}

func (m *SimpleAssetMapper) Map(coreAssets []core.Asset) ([]display.AssetInfo, *big.Float, error) {
	return mapToDisplayAssets(coreAssets)
}

func mapToDisplayAssets(coreAssets []core.Asset) ([]display.AssetInfo, *big.Float, error) {
	var result []display.AssetInfo
	total := big.NewFloat(0)

	for _, a := range coreAssets {
		if a.Balance == "" {
			return nil, nil, fmt.Errorf("пустой баланс для токена %s", a.Symbol)
		}
		if a.USDValue == "" {
			return nil, nil, fmt.Errorf("пустое значение USD для токена %s", a.Symbol)
		}

		balance, _, err := big.ParseFloat(a.Balance, 10, 256, big.ToNearestEven)
		if err != nil {
			return nil, nil, fmt.Errorf("ошибка парсинга баланса %s: %w", a.Balance, err)
		}

		usdStr := strings.TrimPrefix(a.USDValue, "$")
		usdValue, _, err := big.ParseFloat(usdStr, 10, 256, big.ToNearestEven)
		if err != nil {
			return nil, nil, fmt.Errorf("ошибка парсинга USD значения %s: %w", a.USDValue, err)
		}

		result = append(result, display.AssetInfo{
			Symbol:   a.Symbol,
			Balance:  balance,
			USDValue: usdValue,
			USDPrice: usdValue,
		})

		total = new(big.Float).Add(total, usdValue)
	}
	return result, total, nil
}
