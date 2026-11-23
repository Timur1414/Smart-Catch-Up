package config

import "github.com/Timur1414/Smart-Catch-Up/pkg/config"

type GeneratorConfig struct {
	LogPrefix string `yaml:"log_prefix" env:"LOG_PREFIX" env-default:"generator"`
	Port      string `yaml:"port" env:"PORT" env-default:"50052"`
}

type Config struct {
	Common    config.Common   `yaml:"common"`
	Generator GeneratorConfig `yaml:"generator"`
}
