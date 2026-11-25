package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Common struct {
	Debug                      bool          `yaml:"debug"`
	FrontUrl                   string        `yaml:"front_url"`
	ExtensionUrl               string        `yaml:"extension_url"`
	DigestCacheTtl             time.Duration `yaml:"digest_cache_ttl"`
	DigestCachePrefix          string        `yaml:"digest_cache_prefix"`
	AccessTokenExpirationTime  time.Duration `yaml:"access_token_expiration_time"`
	RefreshTokenExpirationTime time.Duration `yaml:"refresh_token_expiration_time"`
}

type PostgresConfig struct {
	Host            string        `env:"POSTGRES_HOST" env-required:"true"`
	Port            string        `env:"POSTGRES_PORT" env-required:"true"`
	DBName          string        `env:"POSTGRES_DB" env-required:"true"`
	User            string        `env:"POSTGRES_USER" env-required:"true"`
	Password        string        `env:"POSTGRES_PASSWORD" env-required:"true"`
	SSLMode         string        `yaml:"ssl_mode"`
	MaxConns        int32         `yaml:"max_conns"`
	MinConns        int32         `yaml:"min_conns"`
	MaxConnLifetime time.Duration `yaml:"max_conn_lifetime"`
}

func (obj PostgresConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		obj.User, obj.Password, obj.Host, obj.Port, obj.DBName, obj.SSLMode)
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" env-required:"true"`
	Port     string `env:"REDIS_PORT" env-required:"true"`
	Password string `env:"REDIS_PASSWORD" env-required:"true"`
	DB       int    `yaml:"db"`
}

func (obj RedisConfig) DSN() string {
	return fmt.Sprintf("redis://:%s@%s:%s/%d", obj.Password, obj.Host, obj.Port, obj.DB)
}

func Load[T any](path string) (*T, error) {
	var config T
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
