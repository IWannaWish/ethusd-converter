package erc20

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

type ERC20Token struct {
	symbol      string
	address     common.Address
	decimals    uint8
	client      *ethclient.Client
	contractABI abi.ABI
}

func NewERC20Token(
	symbol string, address common.Address, decimals uint8, client *ethclient.Client, contractABI abi.ABI) *ERC20Token {
	return &ERC20Token{
		symbol:      symbol,
		address:     address,
		decimals:    decimals,
		client:      client,
		contractABI: contractABI,
	}
}

func (t ERC20Token) GetBalance(ctx context.Context, holder common.Address) (*big.Float, error) {
	callData, err := t.contractABI.Pack("balanceOf", holder)
	if err != nil {
		return nil, err
	}
	msg := ethereum.CallMsg{
		To:   &t.address,
		Data: callData,
	}

	// Можно явно указать последний блок
	latestBlock, err := t.client.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}

	result, err := t.client.CallContract(ctx, msg, new(big.Int).SetUint64(latestBlock))
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("empty result from balanceOf")
	}

	// Логирование для диагностики
	log.Printf("Raw balanceOf result for %s: %x", t.symbol, result)

	balance := new(big.Int).SetBytes(result)

	decimalsBigInt := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(t.decimals)), nil)
	scale := new(big.Float).SetInt(decimalsBigInt)

	return new(big.Float).Quo(new(big.Float).SetInt(balance), scale), nil
}

func (t ERC20Token) GetSymbol() string {
	return t.symbol
}
