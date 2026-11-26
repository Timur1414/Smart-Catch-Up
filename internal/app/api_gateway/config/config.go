package config

import (
	"time"

	"github.com/Timur1414/Smart-Catch-Up/pkg/config"
)

type CsrfConfig struct {
	HeaderName           string        `yaml:"header_name"`
	CookieName           string        `yaml:"cookie_name"`
	CookieExpirationTime time.Duration `yaml:"cookie_expiration_time"`
	Secret               string        `env:"CSRF_SECRET" env-required:"true"`
}

type TargetConfig struct {
	AuthTarget       string `yaml:"auth_target"`
	GeneratorTarget  string `yaml:"generator_target"`
	FileserverTarget string `yaml:"fileserver_target"`
}

type ApiConfig struct {
	LogPrefix             string                `yaml:"log_prefix"`
	Port                  string                `yaml:"port"`
	ReadTimeout           time.Duration         `yaml:"read_timeout"`
	WriteTimeout          time.Duration         `yaml:"write_timeout"`
	Postgres              config.PostgresConfig `yaml:"postgres"`
	Redis                 config.RedisConfig    `yaml:"redis"`
	Csrf                  CsrfConfig            `yaml:"csrf"`
	AccessTokenName       string                `yaml:"access_token_name"`
	RefreshTokenName      string                `yaml:"refresh_token_name"`
	ExtensionUserId       int                   `yaml:"extension_user_id"`
	MaxAvatarSize         int64                 `yaml:"max_avatar_size"`
	CountRawNotifications int                   `yaml:"count_raw_notifications"`
	Targets               TargetConfig          `yaml:"targets"`
}

type Config struct {
	Common config.Common `yaml:"common"`
	Api    ApiConfig     `yaml:"api_gateway"`
}
