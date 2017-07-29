package main

import (
	"log"
	"sync"

	"github.com/codekoala/treksum"
	"github.com/jackc/pgx"
)

var (
	allSeries = []*treksum.Series{
		treksum.NewSeries("Star Trek", "http://chakoteya.net/StarTrek/index.htm"),
		treksum.NewSeries("Star Trek: The Next Generation", "http://chakoteya.net/NextGen/"),
		treksum.NewSeries("Star Trek: Deep Space Nine", "http://chakoteya.net/DS9/"),
		treksum.NewSeries("Star Trek: Voyager", "http://chakoteya.net/Voyager/"),
		treksum.NewSeries("Enterprise", "http://chakoteya.net/Enterprise/"),
	}

	wg sync.WaitGroup
)

func main() {
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
		log.Fatal(err)
	}
	defer pool.Close()

	if err = startFresh(pool); err != nil {
		log.Fatal(err)
	}

	for _, s := range allSeries {
		if err = fetchSeries(pool, s); err != nil {
			log.Printf(`unable to fetch series "%s": %s`, s, err)
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

	if err = series.Save(tx); err != nil {
		log.Printf("unable to create series: %s (%s)", series, err)
		return
	}

	for ep := range fetchEpisodes(series) {
		if err = ep.Save(tx); err != nil {
			log.Printf("error saving episode: %s", err)
		}
	}

	if err = tx.Commit(); err != nil {
		log.Printf("failed to commit transaction: %s", err)
		return
	}

	return nil
}

func fetchEpisodes(series *treksum.Series) <-chan *treksum.Episode {
	out := make(chan *treksum.Episode, 200)

	wg.Add(1)
	go func() {
		episodes, err := treksum.ParseEpisodeList(series)
		if err != nil {
			log.Fatalf("failed to parse episode list (%s): %s", series, err)
		}

		for _, ep := range episodes {
			log.Printf("Fetching episode: %s", ep)
			if err = ep.Fetch(); err != nil {
				log.Printf("failed to parse episode transcript (%s; %s): %s", series, ep, err)
				continue
			}
			out <- ep
		}

		close(out)
		wg.Done()
	}()

	return out
}
