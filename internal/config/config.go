package config

import (
	"github.com/caarlos0/env/v10"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

type Config struct {
	RPCURL          string `env:"RPC_URL,required"`
	ChainlinkETHUSD string `env:"CHAINLINK_ETH_USD,required"`

	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
	UseZap   bool   `env:"LOG_DEV" envDefault:"false"`

	LRUCacheSize         int           `env:"LRU_CACHE_SIZE" default:"100"`
	PriceRefreshInterval time.Duration `env:"PRICE_REFRESH_INTERVAL" default:"5s"`

	Tokens []TokenConfig `yaml:"tokens"`
}

func Load() *Config {
	var cfg Config

	// загрузка ENV
	if err := env.Parse(&cfg); err != nil {
		//todo EC-12 Убрать fatal
		log.Fatalf("failed to load config from env: %v", err)
	}

	// загрузка tokenlist.yaml
	rawYaml, err := os.ReadFile("internal/config/tokenlist.yaml")
	if err != nil {
		log.Fatalf("failed to read tokenlist.yaml: %v", err)
	}
	if err := yaml.Unmarshal(rawYaml, &cfg); err != nil {
		log.Fatalf("failed to parse tokenlist.yaml: %v", err)
	}

	return &cfg
}
