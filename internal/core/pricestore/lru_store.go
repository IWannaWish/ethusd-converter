package pricestore

import (
	"context"
	"github.com/IWannaWish/ethusd-converter/internal/applog"
	"github.com/IWannaWish/ethusd-converter/internal/eth/chainlink"
	"github.com/ethereum/go-ethereum/common/lru"
	"math/big"
	"time"
)

type LruStore struct {
	logger   applog.Logger
	cache    *lru.Cache[string, *big.Float]
	interval time.Duration
	feeds    map[string]chainlink.PriceFeed
}

func NewLruStore(capacity int, logger applog.Logger, interval time.Duration) *LruStore {
	var c = lru.NewCache[string, *big.Float](capacity)
	return &LruStore{
		logger:   logger,
		cache:    c,
		interval: interval,
	}
}

func (l *LruStore) Get(ctx context.Context, symbol string) (*big.Float, bool) {
	if l.cache.Contains(symbol) {
		l.logger.Debug(ctx, "Значение взято из кэша", applog.String("symbol", symbol))
		return l.cache.Get(symbol)
	}

	feed, ok := l.feeds[symbol]
	if !ok || feed == nil {
		l.logger.Error(ctx, "Feed не найден", applog.String("symbol", symbol))
		return nil, false
	}

	price, err := feed.GetUSDPrice(ctx)
	if err != nil {
		l.logger.Error(ctx, "Не удалось получить цену из Chainlink", applog.WithStack(err)...)
		return nil, false
	}

	l.cache.Add(symbol, price)
	l.logger.Debug(ctx, "Значение взято из chainlink", applog.String("symbol", symbol))
	return l.cache.Get(symbol)
}

func (l *LruStore) StartBackgroundUpdater(ctx context.Context) error {
	ticker := time.NewTicker(l.interval)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				l.updateAll()
			}
		}
	}()
	return nil
}

func (l *LruStore) updateAll() {
	keys := l.cache.Keys()
	ctx := context.Background()

	for _, key := range keys {
		feed, ok := l.feeds[key]
		if !ok || feed == nil {
			l.logger.Error(ctx, "Feed не найден во время обновления", applog.String("symbol", key))
			continue
		}

		price, err := feed.GetUSDPrice(ctx)
		if err != nil {
			l.logger.Error(ctx, "Не удалось обновить цену", applog.WithStack(err)...)
			continue
		}
		l.cache.Add(key, price)
		l.logger.Debug(ctx, "price updated in background", applog.String("symbol", key))
	}
}

func (l *LruStore) RegisterFeed(symbol string, feed chainlink.PriceFeed) {
	if l.feeds == nil {
		l.feeds = make(map[string]chainlink.PriceFeed)
	}
	l.feeds[symbol] = feed
}
