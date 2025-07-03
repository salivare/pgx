package sqlx

import (
	"github.com/salivare/pgx/config"
	"net"
	"net/url"
	"strconv"
	"strings"
)

func WithDSN(rawDSN string) Option {
	return func(c *config.DBConfig) {
		u, err := url.Parse(rawDSN)
		if err != nil {
			return
		}

		c.Driver = u.Scheme

		if u.User != nil {
			c.User = u.User.Username()
			if p, ok := u.User.Password(); ok {
				c.Password = p
			}
		}

		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			host = u.Host
		}

		c.Host = host
		if port != "" {
			if pi, err := strconv.Atoi(port); err == nil {
				c.Port = pi
			}
		}

		c.Name = strings.TrimPrefix(u.Path, "/")

		if m := u.Query().Get("sslmode"); m != "" {
			c.SSLMode = m
		}
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
