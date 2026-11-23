package config

import "github.com/Timur1414/Smart-Catch-Up/pkg/config"

type AuthConfig struct {
	LogPrefix string `yaml:"log_prefix" env:"LOG_PREFIX" env-default:"auth"`
	Port      string `yaml:"port" env:"PORT" env-default:"50051"`
}

type Config struct {
	Common config.Common `yaml:"common"`
	Auth   AuthConfig    `yaml:"auth"`
}
