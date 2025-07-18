package source

import (
	"github.com/IWannaWish/ethusd-converter/internal/config"
	"github.com/IWannaWish/ethusd-converter/internal/core"
	"github.com/IWannaWish/ethusd-converter/internal/eth/chainlink"
	"github.com/IWannaWish/ethusd-converter/internal/eth/token/erc20"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type AssetSource struct {
	Token core.TokenBalanceFetcher
	Feed  chainlink.PriceFeed
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

		erc20 := erc20.NewERC20Token(entry.Symbol, tokenAddress, entry.Decimals, client, erc20ABI)
		feed := chainlink.NewChainlinkFeed(client, priceFeedAddress, feedABI)

		sources = append(sources, AssetSource{
			Token: erc20,
			Feed:  feed,
		})
	}

	return sources, nil
}
