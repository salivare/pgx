package config

import (
	"fmt"
	"os"
	"time"

	"github.com/goccy/go-yaml"
)

type DBConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"sslmode"`

	Pool struct {
		MaxOpenConns int `yaml:"maxOpenConns"`
		MaxIdleConns int `yaml:"maxIdleConns"`
	} `yaml:"pool"`

	Retry struct {
		MaxAttempts int           `yaml:"maxAttempts"`
		Delay       time.Duration `yaml:"delay"`
		MaxDelay    time.Duration `yaml:"maxDelay"`
	} `yaml:"retry"`
}

type Config struct {
	Database DBConfig `yaml:"database"`
}

func LoadConfig() (*DBConfig, error) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		if env := os.Getenv("CONFIG_PATH"); env != "" {
			path = env
		} else {
			path = "./config/dbconfig.yaml"
		}
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found: %s", path)
		}
		return nil, fmt.Errorf("error accessing config file %s: %w", path, err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("config path is a directory, not a file: %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	if cfg.Database.Retry.MaxAttempts == 0 {
		cfg.Database.Retry.MaxAttempts = 10
	}
	if cfg.Database.Retry.Delay == 0 {
		cfg.Database.Retry.Delay = 2 * time.Second
	}
	if cfg.Database.Retry.MaxDelay == 0 {
		cfg.Database.Retry.MaxDelay = 10 * time.Second
	}

	return &cfg.Database, nil
}
