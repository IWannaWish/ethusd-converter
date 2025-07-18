package eth

import (
	"context"
	"github.com/IWannaWish/ethusd-converter/internal/core/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

// NativeTokenReader реализует TokenBalanceFetcher для ETH
type NativeTokenReader struct {
	client *ethclient.Client
}

func NewNativeTokenReader(client *ethclient.Client) *NativeTokenReader {
	return &NativeTokenReader{client: client}
}

func (n NativeTokenReader) GetBalance(ctx context.Context, address common.Address) (*big.Float, error) {
	balanceWei, err := n.client.BalanceAt(ctx, address, nil)
	if err != nil {
		return nil, err
	}
	return util.ScaleDown(balanceWei, 18), nil
}

func (n NativeTokenReader) GetSymbol() string {
	return "ETH"
}
