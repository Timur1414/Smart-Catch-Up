package config

import (
	"time"

	"github.com/Timur1414/Smart-Catch-Up/pkg/config"
)

type CsrfConfig struct {
	Secret string `env:"CSRF_SECRET" env-required:"true"`
}

type ApiConfig struct {
	LogPrefix        string                `yaml:"log_prefix" env-default:"api"`
	AuthTarget       string                `yaml:"auth_target" env:"AUTH_TARGET" env-default:"auth:50051"`
	GeneratorTarget  string                `yaml:"generator_target" env:"GENERATOR_TARGET" env-default:"generator:50052"`
	FileserverTarget string                `yaml:"fileserver_target" env:"FILESERVER_TARGET" env-default:"fileserver:50053"`
	Port             string                `yaml:"port" env:"PORT" env-default:"8080"`
	ReadTimeout      time.Duration         `yaml:"read_timeout" env:"READ_TIMEOUT" env-default:"10s"`
	WriteTimeout     time.Duration         `yaml:"write_timeout" env:"WRITE_TIMEOUT" env-default:"120s"`
	Postgres         config.PostgresConfig `yaml:"postgres"`
	Redis            config.RedisConfig    `yaml:"redis"`
	Csrf             CsrfConfig            `yaml:"csrf"`
}

type Config struct {
	Common config.Common `yaml:"common"`
	Api    ApiConfig     `yaml:"api_gateway"`
}
