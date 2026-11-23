package config

import "github.com/Timur1414/Smart-Catch-Up/pkg/config"

type FileserverConfig struct {
	LogPrefix string `yaml:"log_prefix" env:"LOG_PREFIX" env-default:"fileserver"`
	Port      string `yaml:"port" env:"PORT" env-default:"50053"`
}

type Config struct {
	Common     config.Common    `yaml:"common"`
	Fileserver FileserverConfig `yaml:"fileserver"`
}
