package main

import (
	"fmt"
	"os"

	"github.com/jackc/pgx"
	"go.uber.org/zap"

	"github.com/codekoala/treksum/api"
	"github.com/codekoala/treksum/db"
)

var (
	log *zap.Logger
)

func init() {
	var err error

	if log, err = zap.NewProduction(); err != nil {
		fmt.Printf("failed to setup logger: %s", err)
		os.Exit(1)
	}
}

func main() {
	var (
		pool *pgx.ConnPool
		err  error
	)

	defer log.Sync()

	if pool, err = db.Connect(); err != nil {
		log.Fatal("error connecting to database", zap.Error(err))
	}
	defer pool.Close()

	app := api.Setup(log, pool)

	if err = app.Start(":1323"); err != nil {
		log.Warn("server error", zap.Error(err))
	}
}
