package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Port     string `env:"PORT" envDefault:":3000"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"INFO"`
}

func NewConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Errorf("Error loading .env file: %v", err)
		return nil, err
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Errorf("Error parsing env vars: %v", err)
		return nil, err
	}

	return cfg, nil
}
