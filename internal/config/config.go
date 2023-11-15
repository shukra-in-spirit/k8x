package config

import (
	"fmt"

	configEnv "github.com/caarlos0/env/v10"
)

type AWSConfig struct {
	AWSRegion    string `env:"AWS_REGION"     envDefault:"us-west-2"`
	AWSAccessID  string `env:"AWS_ACCESS_ID"`
	AWSSecretKey string `env:"AWS_SECRET_KEY"`
	DBEndpoint   string `env:"DB_ENDPOINT"`
}

type ServiceConfig struct {
	AWSConfig AWSConfig
	PromUrl   string `env:"PROM_URL"`
}

func NewServiceConfig() (ServiceConfig, error) {
	sConfig := ServiceConfig{}

	err := configEnv.Parse(&sConfig)
	if err != nil {
		return ServiceConfig{}, fmt.Errorf("parse config error: %w", err)
	}

	return sConfig, nil
}
