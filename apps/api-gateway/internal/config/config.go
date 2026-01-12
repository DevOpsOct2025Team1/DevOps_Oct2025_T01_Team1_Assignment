package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port            string `env:"PORT" env-default:"3000"`
	AuthServiceAddr string `env:"AUTH_SERVICE_ADDR" env-default:"localhost:8081"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
