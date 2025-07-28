package source

import (
	"context"
	"github.com/IWannaWish/ethusd-converter/internal/applog"
	"github.com/IWannaWish/ethusd-converter/internal/config"
	"github.com/IWannaWish/ethusd-converter/internal/core"
	"github.com/IWannaWish/ethusd-converter/internal/core/pricestore"
	"github.com/IWannaWish/ethusd-converter/internal/eth/chainlink"
	"github.com/IWannaWish/ethusd-converter/internal/eth/token/erc20"
	"github.com/IWannaWish/ethusd-converter/internal/eth/token/eth"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type AssetSource struct {
	Token core.TokenBalanceFetcher
	Feed  pricestore.PriceStore
}

func BuildAssetSources(
	ctx context.Context,
	logger applog.Logger,
	tokenList []config.TokenConfig,
	client *ethclient.Client,
	conf *config.Config,
	erc20ABI abi.ABI,
	feedABI abi.ABI,
) ([]AssetSource, error) {

	var sources []AssetSource
	cache := pricestore.NewLruStore(conf.LRUCacheSize, logger, conf.PriceRefreshInterval)
	logger.Info(ctx, "инициализация общего кэша цен завершена", applog.Int("tokens", len(tokenList)))
	if err := cache.StartBackgroundUpdater(ctx); err != nil {
		logger.Error(ctx, "не удалось запустить фоновое обновление цен", applog.Err(err)...)
	}

	for _, entry := range tokenList {
		feed := chainlink.NewChainlinkFeed(
			client,
			common.HexToAddress(entry.PriceFeedAddress),
			feedABI,
		)
		cache.RegisterFeed(entry.Symbol, feed)

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
			Feed:  cache,
		})
	}

	return sources, nil
}
