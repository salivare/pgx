package db

import (
	_ "github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}
