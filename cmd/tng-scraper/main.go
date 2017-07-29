package main

import (
	"log"
	"sync"

	"github.com/codekoala/treksum"
	"github.com/jackc/pgx"
)

var wg sync.WaitGroup

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

	series := treksum.NewSeries("Star Trek: The Next Generation")
	if err = series.Save(pool); err != nil {
		log.Fatalf("unable to create series: %s", err)
	}

	for ep := range fetchEpisodes(series) {
		if err = ep.Save(pool); err != nil {
			log.Printf("error saving episode: %s", err)
		}
	}

	wg.Wait()
}

func fetchEpisodes(series *treksum.Series) <-chan *treksum.Episode {
	out := make(chan *treksum.Episode, 200)

	wg.Add(1)
	go func() {
		episodes, err := treksum.ParseEpisodeList(series)
		if err != nil {
			log.Fatalf("failed to parse episode list: %s", err)
		}

		for _, ep := range episodes {
			log.Printf("Fetching episode: %s", ep)
			if err = ep.Fetch(); err != nil {
				log.Printf("failed to parse episode transcript: %s", err)
				continue
			}
			out <- ep
		}

		close(out)
		wg.Done()
	}()

	return out
}
