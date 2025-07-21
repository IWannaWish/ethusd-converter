package config

type TokenConfig struct {
	Symbol           string `yaml:"symbol"`
	Type             string `yaml:"type"` // "eth" или "erc20"
	TokenAddress     string `yaml:"tokenAddress"`
	PriceFeedAddress string `yaml:"priceFeedAddress"`
	Decimals         uint8  `yaml:"decimals"`
}
