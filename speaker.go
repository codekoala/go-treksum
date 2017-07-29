package treksum

import (
	"fmt"

	"github.com/jackc/pgx"
)

const (
	CREATE_SPEAKER_TABLE = `
		create table "speaker" (
			"id" serial primary key,
			"series_id" int,
			"name" citext,
			foreign key ("series_id") references "series" ("id"),
			unique ("series_id", "name")
		)
	`

	INSERT_SPEAKER = `
		insert into "speaker"
			("series_id", "name")
		select $1, lower($2)
		where not exists (
			select "id"
			from "speaker"
			where "series_id" = $1
			  and "name" = lower($2)
		)
	`

	SELECT_SPEAKER = `
		SELECT "id"
		FROM "speaker"
		WHERE "series_id" = $1
		  AND "name" = $2
	`
)

type Speaker struct {
	ID     int64   `json:"id"`
	Series *Series `json:"series"`
	Name   string  `json:"name"`
}

func NewSpeaker(series *Series, name string) (s *Speaker) {
	s = &Speaker{
		Series: series,
		Name:   name,
	}

	return s
}

func (this *Speaker) String() string {
	return fmt.Sprintf("%s", this.Name)
}

func (this *Speaker) Save(pool *pgx.ConnPool) (err error) {
	pool.Exec(INSERT_SPEAKER, this.Series.ID, this.Name)

	return pool.QueryRow(SELECT_SPEAKER, this.Series.ID, this.Name).Scan(&this.ID)
}
