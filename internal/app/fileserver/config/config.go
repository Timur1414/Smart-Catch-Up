package config

import "github.com/Timur1414/Smart-Catch-Up/pkg/config"

type FileserverConfig struct {
	LogPrefix string `yaml:"log_prefix"`
	Port      string `yaml:"port"`
	StaticDir string `yaml:"static_dir"`
}

type Config struct {
	Common     config.Common    `yaml:"common"`
	Fileserver FileserverConfig `yaml:"fileserver"`
}
