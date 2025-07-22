package source

import (
	"context"
	"github.com/IWannaWish/ethusd-converter/internal/applog"
	"github.com/IWannaWish/ethusd-converter/internal/config"
	"github.com/IWannaWish/ethusd-converter/internal/core"
	"github.com/IWannaWish/ethusd-converter/internal/eth/chainlink"
	"github.com/IWannaWish/ethusd-converter/internal/eth/token/erc20"
	"github.com/IWannaWish/ethusd-converter/internal/eth/token/eth"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type AssetSource struct {
	Token core.TokenBalanceFetcher
	Feed  chainlink.PriceFeed
}

func BuildAssetSources(
	ctx context.Context,
	logger applog.Logger,
	tokenList []config.TokenConfig,
	client *ethclient.Client,
	erc20ABI abi.ABI,
	feedABI abi.ABI,
) ([]AssetSource, error) {

	var sources []AssetSource

	for _, entry := range tokenList {
		feed := chainlink.NewChainlinkFeed(
			client,
			common.HexToAddress(entry.PriceFeedAddress),
			feedABI,
		)

		var token core.TokenBalanceFetcher

		switch entry.Type {
		case "eth":
			token = eth.NewNativeTokenReader(client)

		case "erc20":
			tokenAddress := common.HexToAddress(entry.TokenAddress)
			token = erc20.NewERC20Token(entry.Symbol, tokenAddress, entry.Decimals, client, erc20ABI)

		default:
			logger.Error(ctx, "неизвестный тип токена — пропускаем",
				applog.String("type", entry.Type),
				applog.String("symbol", entry.Symbol),
			)
			continue
		}

		sources = append(sources, AssetSource{
			Token: token,
			Feed:  feed,
		})
	}

	return sources, nil
}
