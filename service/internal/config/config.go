package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type EmailServiceConfig struct {
}

type Config struct {
	EmailServiceConfig EmailServiceConfig `yaml:"email_service" env-prefix:"EMAIL_SERVICE"`
}

func New() (Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to read env vars: %v", err)
	}

	return cfg, nil
}
