package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Port string `env:"PORT" envDefault:":3000"`

	LogLevel string `env:"LOG_LEVEL" envDefault:"INFO"`

	DbUser string `env:"DB_USER" envDefault:"postgres"`
	DbPass string `env:"DB_PASSWORD" envDefault:"postgres"`
	DbHost string `env:"Db_HOST" envDefault:"localhost"`
	DbPort string `env:"DB_PORT" envDefault:"5432"`
	DbName string `env:"DB_NAME" envDefault:"postgres"`

	SecretKey string `env:"SECRET_KEY" envDefault:"secret"`

	AccessTokenLifetimeMinutes  int `env:"ACCESS_TOKEN_LIFETIME_MINUTES" envDefault:"5"`
	RefreshTokenLifetimeMinutes int `env:"REFRESH_TOKEN_LIFETIME_MINUTES" envDefault:"1440"`
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
