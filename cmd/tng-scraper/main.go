package main

import (
	"log"

	"github.com/codekoala/treksum"
	"github.com/jackc/pgx"
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

	for ep := range fetchEpisodes() {
		if err = saveEpisode(pool, ep); err != nil {
			log.Printf("error saving episode: %s", err)
		}
	}
}

func fetchEpisodes() <-chan *treksum.Episode {
	out := make(chan *treksum.Episode, 200)

	go func() {
		episodes, err := treksum.ParseEpisodeList()
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
	}()

	return out
}
