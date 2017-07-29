package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/codekoala/treksum"
	"github.com/jackc/pgx"
	"go.uber.org/zap"
)

var (
	allSeries = []*treksum.Series{
		treksum.NewSeries("Star Trek", "http://chakoteya.net/StarTrek/index.htm"),
		treksum.NewSeries("Star Trek: The Next Generation", "http://chakoteya.net/NextGen/"),
		treksum.NewSeries("Star Trek: Deep Space Nine", "http://chakoteya.net/DS9/"),
		treksum.NewSeries("Star Trek: Voyager", "http://chakoteya.net/Voyager/"),
		treksum.NewSeries("Enterprise", "http://chakoteya.net/Enterprise/"),
	}

	wg  sync.WaitGroup
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
	defer log.Sync()

	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			Port:     5433,
			Database: "wheaties",
			User:     "wheaties",
		},
		MaxConnections: 20,
	})
	if err != nil {
		log.Fatal("error connecting to database", zap.Error(err))
	}
	defer pool.Close()

	if err = startFresh(pool); err != nil {
		log.Fatal("error cleaning up tables", zap.Error(err))
	}

	for _, s := range allSeries {
		if err = fetchSeries(pool, s); err != nil {
			log.Warn("unable to fetch series", zap.String("name", s.String()), zap.Error(err))
			continue
		}
	}

	wg.Wait()
}

func fetchSeries(pool *pgx.ConnPool, series *treksum.Series) (err error) {
	var tx *pgx.Tx

	if tx, err = pool.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	l := log.With(zap.String("series", series.Name))
	if err = series.Save(tx); err != nil {
		l.Warn("unable to create series", zap.Error(err))
		return
	}

	for ep := range fetchEpisodes(l, series) {
		if err = ep.Save(tx); err != nil {
			l.Warn("error saving episode", zap.String("episode", ep.String()), zap.Error(err))
		}
	}

	if err = tx.Commit(); err != nil {
		l.Warn("failed to commit transaction", zap.Error(err))
		return
	}

	return nil
}

func fetchEpisodes(log *zap.Logger, series *treksum.Series) <-chan *treksum.Episode {
	out := make(chan *treksum.Episode, 200)

	wg.Add(1)
	go func() {
		episodes, err := treksum.ParseEpisodeList(log, series)
		if err != nil {
			log.Warn("failed to parse episode list", zap.Error(err))
		}

		for _, ep := range episodes {
			if err = ep.Fetch(); err != nil {
				log.Warn("failed to parse episode transcript", zap.Error(err))
				continue
			}
			out <- ep
		}

		close(out)
		wg.Done()
	}()

	return out
}
