package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type EmailServiceConfig struct {
}

type RabbitMQConfig struct {
	Host         string `yaml:"host" env:"HOST"`
	Port         int    `yaml:"port" env:"PORT"`
	VirtualHost  string `yaml:"virtual_host" env:"VIRTUAL_HOST"`
	User         string `yaml:"user" env:"USER"`
	Password     string `yaml:"password" env:"PASSWORD"`
	QueueName    string `yaml:"queue_name" env:"QUEUE_NAME"`
	ExchangeName string `yaml:"exchange_name" env:"EXCHANGE_NAME"`
	RoutingKey   string `yaml:"routing_key" env:"ROUTING_KEY"`
}

type Config struct {
	EmailServiceConfig EmailServiceConfig `yaml:"email_service" env-prefix:"EMAIL_SERVICE_"`
	RabbitMQConfig     RabbitMQConfig     `yaml:"rabbitmq" env-prefix:"RABBITMQ_"`
}

func New() (Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to read env vars: %v", err)
	}

	return cfg, nil
}
