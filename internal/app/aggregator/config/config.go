package config

import (
	"time"

	"github.com/Timur1414/Smart-Catch-Up/pkg/config"
)

type LlmConfig struct {
	ApiKey  string `env:"OPENROUTER_API_KEY" env-required:"true"`
	BaseUrl string `env:"OPENROUTER_BASE_URL" env-required:"true"`
	Model   string `env:"OPENROUTER_MODEL" env-required:"true"`
}

type ProxyConfig struct {
	On       bool   `yaml:"on"`
	Host     string `env:"PROXY_HOST" env-required:"true"`
	Login    string `env:"PROXY_LOGIN" env-required:"true"`
	Password string `env:"PROXY_PASSWORD" env-required:"true"`
	Scheme   string `yaml:"scheme" end-default:"socks5"`
}

type AggregatorConfig struct {
	LogPrefix   string                `yaml:"log_prefix"`
	BatchSize   int                   `yaml:"batch_size"`
	WorkerCount int                   `yaml:"worker_count"`
	Interval    time.Duration         `yaml:"interval"`
	Postgres    config.PostgresConfig `yaml:"postgres"`
	Redis       config.RedisConfig    `yaml:"redis"`
	Proxy       ProxyConfig           `yaml:"proxy"`
	Llm         LlmConfig             `yaml:"llm"`
}

type Config struct {
	Common     config.Common    `yaml:"common"`
	Aggregator AggregatorConfig `yaml:"aggregator"`
}
