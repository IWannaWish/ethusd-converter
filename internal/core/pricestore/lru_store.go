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
	feed     chainlink.PriceFeed
	interval time.Duration
}

func NewLruStore(feed chainlink.PriceFeed, capacity int, logger applog.Logger, interval time.Duration) *LruStore {
	var c = lru.NewCache[string, *big.Float](capacity)
	return &LruStore{
		feed:     feed,
		logger:   logger,
		cache:    c,
		interval: interval,
	}
}

func (l LruStore) Get(symbol string) (*big.Float, bool) {

	ctx := context.Background()

	if exist := l.cache.Contains(symbol); exist {
		l.logger.Debug(ctx, "Значение взято из кэша", applog.String("symbol", symbol))
		return l.cache.Get(symbol)
	}

	price, err := l.feed.GetUSDPrice(ctx)
	if err != nil {
		l.logger.Error(context.Background(), "Не удалось получить значение", applog.WithStack(err)...)
		return nil, false
	}

	l.cache.Add(symbol, price)
	l.logger.Debug(ctx, "Значение взято из chainlink", applog.String("symbol", symbol))
	return l.cache.Get(symbol)
}

func (l LruStore) StartBackgroundAdapter(ctx context.Context) error {
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

func (l LruStore) updateAll() {
	keys := l.cache.Keys()
	ctx := context.Background()

	for _, key := range keys {
		price, err := l.feed.GetUSDPrice(ctx)
		if err != nil {
			l.logger.Error(ctx, "Не удалось получить цену", applog.WithStack(err)...)
			continue
		}
		l.cache.Add(key, price)
		l.logger.Debug(context.Background(), "price updated in background", applog.String("symbol", key))
	}
}
