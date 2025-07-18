package abi

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
)

func LoadERC20ABI() (abi.ABI, error) {
	return abi.JSON(strings.NewReader(ERC20ABI))
}

func LoadAggregatorABI() (abi.ABI, error) {
	return abi.JSON(strings.NewReader(AggregatorV3InterfaceABI))
}
