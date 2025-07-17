package core

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"math/big"
)

type ethAssetService struct {
	sources []AssetSource
}

func NewAssetService(sources []AssetSource) AssetService {
	return &ethAssetService{
		sources: sources,
	}
}
func (s *ethAssetService) GetAssets(ctx context.Context, address common.Address) ([]Asset, error) {
	var result []Asset
	var totalUSD = big.NewFloat(0)

	for _, src := range s.sources {
		balance, err := src.Token.GetBalance(ctx, address)
		if err != nil {
			log.Printf("Error getting balance for %s: %v", src.Token.GetSymbol(), err)
			continue
		}
		log.Printf("%s balance: %s", src.Token.GetSymbol(), balance.Text('f', 6))

		price, err := src.Feed.GetUSDPrice(ctx)
		if err != nil {
			log.Printf("Error getting price for %s: %v", src.Token.GetSymbol(), err)
			continue
		}
		log.Printf("%s price: %s", src.Token.GetSymbol(), price.Text('f', 6))

		usdValue := new(big.Float).Mul(balance, price)
		totalUSD.Add(totalUSD, usdValue)

		result = append(result, Asset{
			Symbol:   src.Token.GetSymbol(),
			Balance:  balance.Text('f', 6),
			USDValue: "$" + usdValue.Text('f', 2),
		})
	}

	// Можем добавить итоговую сумму в конце (если хочешь)
	result = append(result, Asset{
		Symbol:   "Total",
		Balance:  "",
		USDValue: "$" + totalUSD.Text('f', 2),
	})

	return result, nil
}
