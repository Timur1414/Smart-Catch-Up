package config

import "github.com/Timur1414/Smart-Catch-Up/pkg/config"

type PostgresCredentialsConfig struct {
	Login    string `env:"GENERATOR_SERVICE_USER_LOGIN" env-required:"true"`
	Password string `env:"GENERATOR_SERVICE_USER_PASSWORD" env-required:"true"`
}

type GeneratorConfig struct {
	LogPrefix           string                    `yaml:"log_prefix"`
	Port                string                    `yaml:"port"`
	Postgres            config.PostgresConfig     `yaml:"postgres"`
	PostgresCredentials PostgresCredentialsConfig `yaml:"postgres_credentials"`
}

type Config struct {
	Common    config.Common   `yaml:"common"`
	Generator GeneratorConfig `yaml:"generator"`
}
