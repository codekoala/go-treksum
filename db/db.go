package db

import (
	"github.com/jackc/pgx"
)

func Connect() (*pgx.ConnPool, error) {
	return pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			Port:     5433,
			Database: "wheaties",
			User:     "wheaties",
		},
		MaxConnections: 20,
	})
}
