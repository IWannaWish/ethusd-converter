package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type TokenConfig struct {
	Symbol           string `yaml:"symbol"`
	TokenAddress     string `yaml:"tokenAddress"`
	PriceFeedAddress string `yaml:"priceFeedAddress"`
	Decimals         uint8  `yaml:"decimals"`
}

type TokenList struct {
	Tokens []TokenConfig `yaml:"tokens"`
}

func LoadTokenList(path string) (*TokenList, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var list TokenList
	if err := yaml.Unmarshal(data, &list); err != nil {
		return nil, err
	}

	return &list, nil
}
