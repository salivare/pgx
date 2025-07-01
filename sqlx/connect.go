package sqlx

import (
	"fmt"
	"net/url"

	"github.com/avast/retry-go"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/salivare/pgx/config"
)

type Option func(*config.DBConfig)

func DSN(c *config.DBConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		url.QueryEscape(c.User),
		url.QueryEscape(c.Password),
		c.Host,
		c.Port,
		c.Name,
		c.SSLMode,
	)
}

func Connect(opts ...Option) (*sqlx.DB, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(cfg)
	}

	var db *sqlx.DB
	err = retry.Do(
		func() error {
			var err error
			db, err = sqlx.Connect(cfg.Driver, DSN(cfg))
			return err
		},
		retry.Attempts(uint(cfg.Retry.MaxAttempts)),
		retry.Delay(cfg.Retry.Delay),
		retry.MaxDelay(cfg.Retry.MaxDelay),
		retry.DelayType(retry.BackOffDelay),
	)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.Pool.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Pool.MaxIdleConns)
	return db, nil
}
