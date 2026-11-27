package config

import (
	"time"

	"github.com/Timur1414/Smart-Catch-Up/pkg/config"
)

type EncryptorConfig struct {
	Secret string `env:"ENCRYPTOR_KEY" env-required:"true"`
}

type JwtConfig struct {
	Secret  string `env:"JWT_SECRET" env-required:"true"`
	Version string `env:"JWT_VERSION" env-required:"true"`
}

type PostgresCredentialsConfig struct {
	Login    string `env:"AUTH_SERVICE_USER_LOGIN" env-required:"true"`
	Password string `env:"AUTH_SERVICE_USER_PASSWORD" env-required:"true"`
}

type AuthConfig struct {
	LogPrefix                    string                    `yaml:"log_prefix"`
	Port                         string                    `yaml:"port"`
	Postgres                     config.PostgresConfig     `yaml:"postgres"`
	Redis                        config.RedisConfig        `yaml:"redis"`
	Encryptor                    EncryptorConfig           `yaml:"encryptor"`
	Jwt                          JwtConfig                 `yaml:"jwt"`
	RefreshTokenPrefix           string                    `yaml:"refresh_token_prefix"`
	TwoFactorFlowPrefix          string                    `yaml:"two_factor_flow_prefix"`
	TwoFactorAttemptsPrefix      string                    `yaml:"two_factor_attempts_prefix"`
	TwoFactorTokenExpirationTime time.Duration             `yaml:"two_factor_token_expiration_time"`
	UsedCodePrefix               string                    `yaml:"used_code_prefix"`
	UsedCodeExpirationTime       time.Duration             `yaml:"used_code_expiration_time"`
	AccessTokenType              string                    `yaml:"access_token_type"`
	RefreshTokenType             string                    `yaml:"refresh_token_type"`
	TempTokenType                string                    `yaml:"temp_token_type"`
	TotpIssuer                   string                    `env:"TOTP_ISSUER" env-required:"true"`
	PostgresCredentials          PostgresCredentialsConfig `yaml:"postgres_credentials"`
}

type Config struct {
	Common config.Common `yaml:"common"`
	Auth   AuthConfig    `yaml:"auth"`
}
