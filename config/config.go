package config

import (
	"fmt"
	"os"
	"sync"
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

var (
	cfgOnce sync.Once
	cfgInst *DBConfig
	cfgErr  error
)

func LoadConfig() (*DBConfig, error) {
	cfgOnce.Do(
		func() {
			path := os.Getenv("CONFIG_PATH")
			if path == "" {
				path = "./config/dbconfig.yaml"
			}

			info, err := os.Stat(path)
			if err != nil {
				if os.IsNotExist(err) {
					cfgErr = fmt.Errorf("config file not found: %s", path)
				} else {
					cfgErr = fmt.Errorf("error accessing config file %s: %w", path, err)
				}
				return
			}
			if info.IsDir() {
				cfgErr = fmt.Errorf("config path is a directory, not a file: %s", path)
				return
			}

			data, err := os.ReadFile(path)
			if err != nil {
				cfgErr = fmt.Errorf("read config file: %w", err)
				return
			}
			var tmp Config
			if err := yaml.Unmarshal(data, &tmp); err != nil {
				cfgErr = fmt.Errorf("unmarshal yaml: %w", err)
				return
			}

			// Default values for retry
			if tmp.Database.Retry.MaxAttempts == 0 {
				tmp.Database.Retry.MaxAttempts = 10
			}
			if tmp.Database.Retry.Delay == 0 {
				tmp.Database.Retry.Delay = 2 * time.Second
			}
			if tmp.Database.Retry.MaxDelay == 0 {
				tmp.Database.Retry.MaxDelay = 10 * time.Second
			}

			cfgInst = &tmp.Database
		},
	)

	return cfgInst, cfgErr
}
