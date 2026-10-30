package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL  string   `env:"DATABASE_URL"  envDefault:"postgres://postgres:postgres@localhost:5432/lab?sslmode=disable"`
	KafkaBrokers []string `env:"KAFKA_BROKERS" envDefault:"localhost:9092"                                                  envSeparator:","`
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
