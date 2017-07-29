package main

import (
	"log"

	"github.com/codekoala/treksum"
	"github.com/jackc/pgx"
)

func startFresh(pool *pgx.ConnPool) (err error) {
	log.Printf("dropping tables")
	if _, err = pool.Exec(`drop table if exists "series", "episode", "speaker", "line" cascade`); err != nil {
		return
	}

	log.Printf("creating series table")
	if _, err = pool.Exec(treksum.CREATE_SERIES_TABLE); err != nil {
		return
	}

	log.Printf("creating episodes table")
	if _, err = pool.Exec(treksum.CREATE_EPISODE_TABLE); err != nil {
		return
	}

	log.Printf("creating speakers table")
	if _, err = pool.Exec(treksum.CREATE_SPEAKER_TABLE); err != nil {
		return
	}

	log.Printf("creating lines table")
	if _, err = pool.Exec(treksum.CREATE_LINE_TABLE); err != nil {
		return
	}

	return nil
}
