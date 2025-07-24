package mapper

import (
	"github.com/IWannaWish/ethusd-converter/cmd/cli/display"
	"github.com/IWannaWish/ethusd-converter/proto/ethusd/gen"
	"math/big"
)

func FromProtoAssets(protoAssets []*gen.Asset) ([]display.AssetInfo, *big.Float, error) {
	var assets []display.AssetInfo
	total := big.NewFloat(0)

	for _, a := range protoAssets {
		balance, err := ParseBalance(a.Balance)
		if err != nil {
			return nil, nil, err
		}

		usdValue, err := ParseUSDValue(a.UsdValue)
		if err != nil {
			return nil, nil, err
		}

		usdPrice := ComputeUSDPrice(usdValue, balance)

		assets = append(assets, display.AssetInfo{
			Symbol:   a.Symbol,
			Balance:  balance,
			USDValue: usdValue,
			USDPrice: usdPrice,
		})

		total = new(big.Float).Add(total, usdValue)
	}

	return assets, total, nil
}
