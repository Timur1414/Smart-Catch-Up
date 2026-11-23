package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Common struct {
	Debug    bool           `yaml:"debug" env:"DEBUG" env-default:"false"`
	Postgres PostgresConfig `yaml:"postgres"`
}

type PostgresConfig struct {
	Host            string        `yaml:"host" env:"POSTGRES_HOST" env-required:"true"`
	Port            string        `yaml:"port" env:"POSTGRES_PORT" env-required:"true"`
	DBName          string        `yaml:"dbname" env:"POSTGRES_DB" env-required:"true"`
	User            string        `yaml:"user" env:"POSTGRES_USER" env-required:"true"`
	Password        string        `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
	SSLMode         string        `yaml:"ssl_mode" env:"POSTGRES_SSLMODE" env-default:"disable"`
	MaxConns        int32         `yaml:"max_conns" env:"POSTGRES_MAX_CONNS" env-default:"20"`
	MinConns        int32         `yaml:"min_conns" env:"POSTGRES_MIN_CONNS" env-default:"5"`
	MaxConnLifetime time.Duration `yaml:"max_conn_lifetime" env:"POSTGRES_MAX_CONN_LIFETIME" env-default:"1h"`
}

func (obj PostgresConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		obj.User, obj.Password, obj.Host, obj.Port, obj.DBName, obj.SSLMode)
}

func Load[T any](path string) (*T, error) {
	var config T
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
