package mapper

import (
	"fmt"
	"github.com/IWannaWish/ethusd-converter/cmd/cli/display"
	"github.com/IWannaWish/ethusd-converter/internal/core"
	"math/big"
)

type DisplayAssetMapper struct{}

func NewDisplayAssetMapper() *DisplayAssetMapper {
	return &DisplayAssetMapper{}
}

func (m *DisplayAssetMapper) Map(coreAssets []core.Asset) ([]display.AssetInfo, *big.Float, error) {
	var result []display.AssetInfo
	total := big.NewFloat(0)

	for _, a := range coreAssets {
		if a.Balance == "" || a.USDValue == "" {
			return nil, nil, fmt.Errorf("некорректный токен: %s", a.Symbol)
		}

		balance, err := ParseBalance(a.Balance)
		if err != nil {
			return nil, nil, err
		}

		usdValue, err := ParseUSDValue(a.USDValue)
		if err != nil {
			return nil, nil, err
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
