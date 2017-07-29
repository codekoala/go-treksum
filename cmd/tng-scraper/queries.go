package main

import (
	"log"

	"github.com/codekoala/treksum"
	"github.com/jackc/pgx"
)

const (
	CREATE_EPISODE_TABLE = `
		create table "episode" (
			"id" serial primary key,
			"title" text,
			"number" int,
			"url" text,
			"airdate" date
		)
	`

	CREATE_SPEAKER_TABLE = `
		create table "speaker" (
			"id" serial primary key,
			"name" citext,
			unique("name")
		)
	`

	CREATE_LINE_TABLE = `
		create table "line" (
			"id" serial primary key,
			"episode_id" int,
			"speaker_id" int,
			"line" text,
			foreign key ("episode_id") references "episode" ("id"),
			foreign key ("speaker_id") references "speaker" ("id")
		)
	`

	INSERT_EPISODE = `
		insert into "episode"
			("title", "number", "url", "airdate")
		values
			($1, $2, $3, $4::date)
		returning ("id")
	`

	INSERT_SPEAKER = `
		insert into "speaker"
			("name")
		select lower($1)
		where not exists (
			select "id"
			from "speaker"
			where "name" = lower($1)
		)
	`

	INSERT_LINE = `
		insert into "line"
			("episode_id", "speaker_id", "line")
		select $1, "id", $3
		from "speaker"
		where "name" = $2
		returning ("id")
	`
)

func startFresh(pool *pgx.ConnPool) (err error) {
	log.Printf("dropping episodes")
	if _, err = pool.Exec(`drop table if exists "episode" cascade`); err != nil {
		return
	}

	log.Printf("dropping speakers")
	if _, err = pool.Exec(`drop table if exists "speaker" cascade`); err != nil {
		return
	}

	log.Printf("dropping lines")
	if _, err = pool.Exec(`drop table if exists "line" cascade`); err != nil {
		return
	}

	log.Printf("creating episodes")
	if _, err = pool.Exec(CREATE_EPISODE_TABLE); err != nil {
		return
	}

	log.Printf("creating speakers")
	if _, err = pool.Exec(CREATE_SPEAKER_TABLE); err != nil {
		return
	}

	log.Printf("creating lines")
	if _, err = pool.Exec(CREATE_LINE_TABLE); err != nil {
		return
	}

	return nil
}

func saveEpisode(pool *pgx.ConnPool, ep *treksum.Episode) (err error) {
	log.Printf("inserting episode: %s", ep)
	err = pool.QueryRow(INSERT_EPISODE, ep.Title, ep.Number, ep.Url, ep.Airdate).Scan(&ep.ID)
	if err != nil {
		return
	}

	for _, line := range ep.Script {
		pool.Exec(INSERT_SPEAKER, line.Speaker)
		err = pool.QueryRow(INSERT_LINE, ep.ID, line.Speaker, line.Line).Scan(&line.ID)
		if err != nil {
			log.Printf("failed to save line: %s (%s)", line, err)
		}
	}

	return nil
}
