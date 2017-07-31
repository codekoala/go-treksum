package db

import (
	"github.com/jackc/pgx"

	"github.com/codekoala/go-treksum/config"
)

func Connect() (*pgx.ConnPool, error) {
	return pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     config.Global.DbHost,
			Port:     config.Global.DbPort,
			Database: config.Global.DbName,
			User:     config.Global.DbUser,
			Password: config.Global.DbPassword,
		},
		MaxConnections: 20,
	})
}
