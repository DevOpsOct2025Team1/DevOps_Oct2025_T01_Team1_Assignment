package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port                string `env:"PORT" env-default:"3000"`
	AuthServiceAddr     string `env:"AUTH_SERVICE_ADDR" env-default:"localhost:8081"`
	UserServiceAddr     string `env:"USER_SERVICE_ADDR" env-default:"localhost:8080"`
	AxiomToken          string `env:"AXIOM_API_TOKEN"`
	AxiomEndpoint       string `env:"AXIOM_ENDPOINT" env-default:"us-east-1.aws.edge.axiom.co"`
	AxiomDataset        string `env:"AXIOM_DATASET" env-default:"traces"`
	AxiomMetricsDataset string `env:"AXIOM_METRICS_DATASET" env-default:"metrics"`
	Environment         string `env:"ENVIRONMENT" env-default:"development"`
	FrontendURL         string `env:"FRONTEND_URL" env-default:"http://localhost:3000"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
