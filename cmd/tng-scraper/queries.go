package main

import (
	"github.com/codekoala/treksum"
	"github.com/jackc/pgx"
)

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
