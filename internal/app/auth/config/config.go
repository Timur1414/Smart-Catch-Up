package config

import "github.com/Timur1414/Smart-Catch-Up/pkg/config"

type EncryptorConfig struct {
	Secret string `env:"ENCRYPTOR_KEY" env-required:"true"`
}

type JwtConfig struct {
	Secret  string `env:"JWT_SECRET" env-required:"true"`
	Version string `env:"JWT_VERSION" env-required:"true"`
}

type AuthConfig struct {
	LogPrefix string                `yaml:"log_prefix" env-default:"auth"`
	Port      string                `yaml:"port" env-default:"50051"`
	Postgres  config.PostgresConfig `yaml:"postgres"`
	Redis     config.RedisConfig    `yaml:"redis"`
	Encryptor EncryptorConfig       `yaml:"encryptor"`
	Jwt       JwtConfig             `yaml:"jwt"`
}

type Config struct {
	Common config.Common `yaml:"common"`
	Auth   AuthConfig    `yaml:"auth"`
}
