package treksum

import (
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx"
)

const (
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

type Line struct {
	ID      int64    `json:"id"`
	Episode *Episode `json:"episode,omitempty"`
	Speaker string   `json:"speaker"`
	Line    string   `json:"line"`
}

func NewLine(speaker, text string) (l *Line) {
	l = &Line{
		Speaker: speaker,
	}
	l.AddText(text)

	return l
}

func (this *Line) AddText(text string) {
	this.Line = strings.TrimSpace(fmt.Sprintf("%s %s", this.Line, text))
}

func (this *Line) String() string {
	return fmt.Sprintf("%s: %s", this.Speaker, this.Line)
}

func (this *Line) Save(pool *pgx.ConnPool) (err error) {
	pool.Exec(INSERT_SPEAKER, this.Speaker)
	err = pool.QueryRow(INSERT_LINE, this.Episode.ID, this.Speaker, this.Line).Scan(&this.ID)
	if err != nil {
		log.Printf("failed to save line: %s (%s)", this, err)
	}

	return
}
