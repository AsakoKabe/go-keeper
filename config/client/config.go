package client

import (
	"log"

	"github.com/caarlos0/env/v10"
)

// Config struct responsible for storing client config
type Config struct {
	ServerAddr string `env:"SERVER_ADDR"`
	Timeout    int    `env:"TIMEOUT"`
}

// LoadConfig функция для загрузки конфигурации.
func LoadConfig() (*Config, error) {
	cfg := new(Config)

	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg, nil
}
