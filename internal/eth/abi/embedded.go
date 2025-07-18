package abi

import _ "embed"

//go:embed files/ERC20.abi.json
var ERC20ABI string

//go:embed files/AggregatorV3Interface.abi.json
var AggregatorV3InterfaceABI string
