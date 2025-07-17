package core

import (
	"context"
	"fmt"
	"github.com/IWannaWish/ethusd-converter/internal/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type ethAssetService struct {
	client *ethclient.Client
	feed   eth.PriceFeed
}

func NewAssetService(client *ethclient.Client, feed eth.PriceFeed) AssetService {
	return &ethAssetService{
		client: client,
		feed:   feed,
	}
}

func (s *ethAssetService) GetAssets(ctx context.Context, address common.Address) ([]Asset, error) {
	balance, err := eth.GetETHBalance(s.client, address)
	if err != nil {
		return nil, fmt.Errorf("error reading balance: %w", err)
	}

	price, err := s.feed.GetUSDPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("error reading price: %w", err)
	}

	total := new(big.Float).Mul(balance, price)

	asset := Asset{
		Symbol:   "ETH",
		Balance:  balance.Text('f', 6),
		USDValue: "$" + total.Text('f', 2),
	}

	return []Asset{asset}, nil
}
