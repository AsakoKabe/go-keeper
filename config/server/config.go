package server

import (
	"log"

	"github.com/caarlos0/env/v10"
)

// Config структура для хранения конфигурации приложения
type Config struct {
	Addr        string `env:"SERVER_ADDRESS"`
	DatabaseDSN string `env:"DATABASE_DSN"`
	CertFile    string `env:"CERT_FILE"`
	KeyFile     string `env:"KEY_FILE"`
	EnableHTTPS bool   `env:"ENABLE_HTTPS"`
	SecretHash  string `env:"SECRET_HASH"`
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
