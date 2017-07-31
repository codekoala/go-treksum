package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx"
	"go.uber.org/zap"

	"github.com/codekoala/go-treksum"
	"github.com/codekoala/go-treksum/db"
)

var (
	allSeries = []*treksum.Series{
		treksum.NewSeries("Star Trek", "http://chakoteya.net/StarTrek/index.htm"),
		treksum.NewSeries("Star Trek: The Next Generation", "http://chakoteya.net/NextGen/"),
		treksum.NewSeries("Star Trek: Deep Space Nine", "http://chakoteya.net/DS9/"),
		treksum.NewSeries("Star Trek: Voyager", "http://chakoteya.net/Voyager/"),
		treksum.NewSeries("Star Trek: Enterprise", "http://chakoteya.net/Enterprise/"),
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
	var (
		pool *pgx.ConnPool
		err  error
	)

	defer log.Sync()

	if pool, err = db.Connect(); err != nil {
		log.Fatal("error connecting to database", zap.Error(err))
	}
	defer pool.Close()

	// wipe out all existing tables
	if err = startFresh(pool); err != nil {
		log.Fatal("error cleaning up tables", zap.Error(err))
	}

	// scrape all series sequentially
	for _, s := range allSeries {
		if err = fetchSeries(pool, s); err != nil {
			log.Warn("unable to fetch series", zap.String("name", s.String()), zap.Error(err))
			continue
		}
	}

	wg.Wait()
}

// startFresh drops and recreates all tables required for treksum.
func startFresh(pool *pgx.ConnPool) (err error) {
	log.Info("dropping tables")
	if _, err = pool.Exec(`drop table if exists "series", "episode", "speaker", "line" cascade`); err != nil {
		return
	}

	log.Info("creating series table")
	if _, err = pool.Exec(treksum.CREATE_SERIES_TABLE); err != nil {
		return
	}

	log.Info("creating episodes table")
	if _, err = pool.Exec(treksum.CREATE_EPISODE_TABLE); err != nil {
		return
	}

	log.Info("creating speakers table")
	if _, err = pool.Exec(treksum.CREATE_SPEAKER_TABLE); err != nil {
		return
	}

	log.Info("creating lines table")
	if _, err = pool.Exec(treksum.CREATE_LINE_TABLE); err != nil {
		return
	}

	return nil
}

// fetchSeries prepares a transaction for inserting all information about the specified series into the database.
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

// fetchEpisodes scrapes the transcript for all episodes of the specified series.
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
