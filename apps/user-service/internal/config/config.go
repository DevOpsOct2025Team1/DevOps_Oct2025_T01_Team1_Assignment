package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port                 string `env:"PORT" env-default:"8080"`
	MongoDBURI           string `env:"MONGODB_URI" env-required:"true"`
	MongoDBDatabase      string `env:"MONGODB_DATABASE" env-default:"user_service"`
	AxiomToken           string `env:"AXIOM_API_TOKEN"`
	AxiomEndpoint        string `env:"AXIOM_ENDPOINT" env-default:"us-east-1.aws.edge.axiom.co"`
	AxiomDataset         string `env:"AXIOM_DATASET" env-default:"traces"`
	AxiomMetricsDataset  string `env:"AXIOM_METRICS_DATASET" env-default:"metrics"`
	Environment          string `env:"ENVIRONMENT" env-default:"development"`
	DefaultAdminUsername string `env:"DEFAULT_ADMIN_USERNAME" env-default:"admin"`
	DefaultAdminPassword string `env:"DEFAULT_ADMIN_PASSWORD" env-default:"changeme"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
