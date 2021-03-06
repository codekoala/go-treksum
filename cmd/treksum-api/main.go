package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jackc/pgx"
	"go.uber.org/zap"

	"github.com/codekoala/go-treksum"
	"github.com/codekoala/go-treksum/api"
	"github.com/codekoala/go-treksum/config"
	"github.com/codekoala/go-treksum/db"
)

const appname = "treksum-api"

var (
	log *zap.Logger
)

func init() {
	var err error

	treksum.AppInfo.App = appname

	flag.Parse()
	if flag.NArg() > 0 {
		switch flag.Arg(0) {
		case "version":
			fmt.Println(treksum.AppInfo.String())
			os.Exit(0)
		}
	}

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

	if err = app.Start(config.Global.ApiAddr); err != nil {
		log.Warn("server error", zap.Error(err))
	}
}
