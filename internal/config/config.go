package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Port string `env:"PORT" envDefault:":3000"`

	LogLevel string `env:"LOG_LEVEL" envDefault:"INFO"`

	DB_USER string `env:"DB_USER" envDefault:"postgres"`
	DB_PASS string `env:"DB_PASSWORD" envDefault:"postgres"`
	DB_HOST string `env:"DB_HOST" envDefault:"localhost"`
	DB_PORT string `env:"DB_PORT" envDefault:"5432"`
	DB_NAME string `env:"DB_NAME" envDefault:"postgres"`

	SECRET_KEY string `env:"SECRET_KEY" envDefault:"secret"`
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load(".env")

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Errorf("Error parsing env vars: %v", err)
		return nil, err
	}

	return cfg, nil
}
