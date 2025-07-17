package core

import (
	"github.com/IWannaWish/ethusd-converter/internal/config"
	"github.com/IWannaWish/ethusd-converter/internal/eth"
	"github.com/IWannaWish/ethusd-converter/internal/eth/token"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type AssetSource struct {
	Token token.TokenBalanceFetcher
	Feed  eth.PriceFeed
}

func BuildAssetSources(
	tokenList []config.TokenConfig,
	client *ethclient.Client,
	erc20ABI abi.ABI,
	feedABI abi.ABI,
) ([]AssetSource, error) {

	var sources []AssetSource

	for _, entry := range tokenList {
		tokenAddress := common.HexToAddress(entry.TokenAddress)
		priceFeedAddress := common.HexToAddress(entry.PriceFeedAddress)

		erc20 := token.NewERC20Token(entry.Symbol, tokenAddress, entry.Decimals, client, erc20ABI)
		feed := eth.NewChainlinkFeed(client, priceFeedAddress, feedABI)

		sources = append(sources, AssetSource{
			Token: erc20,
			Feed:  feed,
		})
	}

	return sources, nil
}
