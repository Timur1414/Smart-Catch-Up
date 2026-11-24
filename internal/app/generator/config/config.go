package config

import "github.com/Timur1414/Smart-Catch-Up/pkg/config"

type GeneratorConfig struct {
	LogPrefix string                `yaml:"log_prefix" env-default:"generator"`
	Port      string                `yaml:"port" env-default:"50052"`
	Postgres  config.PostgresConfig `yaml:"postgres"`
}

type Config struct {
	Common    config.Common   `yaml:"common"`
	Generator GeneratorConfig `yaml:"generator"`
}
