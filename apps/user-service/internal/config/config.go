package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port            string `env:"PORT" env-default:"8080"`
	MongoDBURI      string `env:"MONGODB_URI" env-required:"true"`
	MongoDBDatabase string `env:"MONGODB_DATABASE" env-default:"user_service"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
