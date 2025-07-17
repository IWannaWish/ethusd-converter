package eth

import (
	_ "embed"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
)

//go:embed abi/AggregatorV3Interface.abi.json
var aggregatorABIData string

//go:embed abi/ERC20.abi.json
var erc20ABIData string

func LoadAggregatorABI() (abi.ABI, error) {
	return abi.JSON(strings.NewReader(aggregatorABIData))
}

func LoadERC20ABI() (abi.ABI, error) {
	return abi.JSON(strings.NewReader(erc20ABIData))
}
