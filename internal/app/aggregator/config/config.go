package config

import (
	"time"

	"github.com/Timur1414/Smart-Catch-Up/pkg/config"
)

type AggregatorConfig struct {
	LogPrefix   string        `yaml:"log_prefix" env:"LOG_PREFIX" env-default:"aggregator"`
	BatchSize   int           `yaml:"batch_size" env:"BATCH_SIZE" env-default:"100"`
	WorkerCount int           `yaml:"worker_count" env:"WORKER_COUNT" env-default:"1"`
	Interval    time.Duration `yaml:"interval" env:"INTERVAL" env-default:"10s"`
}

type Config struct {
	Common     config.Common    `yaml:"common"`
	Aggregator AggregatorConfig `yaml:"aggregator"`
}
