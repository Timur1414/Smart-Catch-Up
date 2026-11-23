package config

import "github.com/Timur1414/Smart-Catch-Up/pkg/config"

type ApiConfig struct {
	LogPrefix        string `yaml:"log_prefix" env:"LOG_PREFIX" env-default:"api"`
	AuthTarget       string `yaml:"auth_target" env:"AUTH_TARGET" env-default:"auth:50051"`
	GeneratorTarget  string `yaml:"generator_target" env:"GENERATOR_TARGET" env-default:"generator:50052"`
	FileserverTarget string `yaml:"fileserver_target" env:"FILESERVER_TARGET" env-default:"fileserver:50053"`
	Port             string `yaml:"port" env:"PORT" env-default:"8080"`
}

type Config struct {
	Common   config.Common         `yaml:"common"`
	Api      ApiConfig             `yaml:"api_gateway"`
	Postgres config.PostgresConfig `yaml:"postgres"`
}
