package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	GrpcServer GrpcServerConfig `mapstructure:"grpcServer"`
}

type GrpcServerConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	TLS      bool   `mapstructure:"tls"`
	CertFile string `mapstructure:"tls-cert"`
	KeyFile  string `mapstructure:"tls-key"`
	Enabled  bool   `mapstructure:"enabled"`
}

func New(configFile string) (Config, error) {
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("failed to read the configuration file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("failed to parse the configuration: %v", err)
	}

	return config, nil
}
