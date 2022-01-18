package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RedisAddr string `split_words:"true" default:":6379"`
	RedisPass string `split_words:"true"`
	HttpAddr  string `split_words:"true"`
	HostsFile string `split_words:"true"`
}

func LoadConfig() (*Config, error) {
	var c Config
	if err := envconfig.Process("foxylock", &c); err != nil {
		return nil, fmt.Errorf("failed to process config, %w", err)
	}
	if err := validator.New().Struct(c); err != nil {
		return nil, fmt.Errorf("invalid config, %w", err)
	}

	return &c, nil
}
