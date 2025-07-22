package source

import (
	"context"
	"github.com/IWannaWish/ethusd-converter/internal/applog"

	"github.com/IWannaWish/ethusd-converter/internal/core"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type ethAssetService struct {
	sources []AssetSource
	logger  applog.Logger
}

func NewAssetService(sources []AssetSource, logger applog.Logger) core.AssetService {
	return &ethAssetService{
		sources: sources,
		logger:  logger,
	}
}
func (s *ethAssetService) GetAssets(ctx context.Context, address common.Address) ([]core.Asset, error) {
	var result []core.Asset

	for _, src := range s.sources {
		symbol := src.Token.GetSymbol()
		balance, err := src.Token.GetBalance(ctx, address)
		if err != nil {
			s.logger.Error(ctx, "не удалось получить баланс", applog.Err(err, applog.String("symbol", symbol))...)
			continue
		}
		s.logger.Debug(ctx, "токен: баланс получен", applog.String("symbol", symbol), applog.String("balance", balance.Text('f', 6)))
		price, err := src.Feed.GetUSDPrice(ctx)
		if err != nil {
			s.logger.Error(ctx, "не удалось получить цену", applog.Err(err, applog.String("symbol", symbol))...)
			continue
		}
		s.logger.Debug(ctx, "токен: цена получена", applog.String("symbol", symbol), applog.String("price", price.Text('f', 6)))

		usdValue := new(big.Float).Mul(balance, price)

		result = append(result, core.Asset{
			Symbol:   symbol,
			Balance:  balance.Text('f', 6),
			USDValue: "$" + usdValue.Text('f', 2),
		})
	}
	return result, nil
}
