package eth

import (
	_ "embed"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
)

//go:embed abi/AggregatorV3Interface.abi.json
var aggregatorABIData string

func LoadAggregatorABI() (abi.ABI, error) {
	return abi.JSON(strings.NewReader(aggregatorABIData))
}
