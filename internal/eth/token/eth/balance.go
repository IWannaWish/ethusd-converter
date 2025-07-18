package eth

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

func GetETHBalance(client *ethclient.Client, address common.Address) (*big.Float, error) {
	balanceWei, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return nil, err
	}

	return new(big.Float).Quo(
		new(big.Float).SetInt(balanceWei),
		big.NewFloat(1e18),
	), nil
}
