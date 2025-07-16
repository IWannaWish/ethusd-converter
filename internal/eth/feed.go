package eth

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type PriceFeed interface {
	GetUSDPrice() (*big.Float, error)
}

type ChainlinkFeed struct {
	Address common.Address
	Client  *ethclient.Client
	ABI     abi.ABI
}

type roundData struct {
	RoundID         *big.Int `abi:"roundId"`
	Answer          *big.Int `abi:"answer"`
	StartedAt       *big.Int `abi:"startedAt"`
	UpdatedAt       *big.Int `abi:"updatedAt"`
	AnsweredInRound *big.Int `abi:"answeredInRound"`
}

func NewChainlinkFeed(client *ethclient.Client, feedAddress common.Address, abi abi.ABI) *ChainlinkFeed {
	return &ChainlinkFeed{
		Address: feedAddress,
		Client:  client,
		ABI:     abi,
	}
}

const methodLatestRound = "latestRoundData"

func (c *ChainlinkFeed) GetUSDPrice() (*big.Float, error) {
	callData, err := c.ABI.Pack(methodLatestRound)
	if err != nil {
		return nil, fmt.Errorf("pack error: %w", err)
	}

	result, err := c.Client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &c.Address,
		Data: callData,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("eth_call failed: %w", err)
	}

	var rd roundData
	if err := c.ABI.UnpackIntoInterface(&rd, methodLatestRound, result); err != nil {
		return nil, fmt.Errorf("unpack error: %w", err)
	}

	return new(big.Float).Quo(new(big.Float).SetInt(rd.Answer), big.NewFloat(1e8)), nil
}
