package sqlx

import "github.com/salivare/pgx/config"

func WithDSN(dsn string) Option {
	return func(c *config.DBConfig) {
		c.User = ""
		c.Password = ""
		c.Host = ""
		c.Port = 0
		c.Name = ""
		c.SSLMode = ""
	}
}

func WithMaxOpenConns(n int) Option {
	return func(c *config.DBConfig) {
		c.Pool.MaxOpenConns = n
	}
}

func WithRetryAttempts(n int) Option {
	return func(c *config.DBConfig) {
		c.Retry.MaxAttempts = n
	}
}
