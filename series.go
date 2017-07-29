package treksum

import (
	"fmt"

	"github.com/jackc/pgx"
)

const (
	CREATE_SERIES_TABLE = `
		create table "series" (
			"id" serial primary key,
			"name" citext,
			"url" text,
			unique("name")
		)
	`

	INSERT_SERIES = `
		insert into "series"
			("name", "url")
		values
			($1, $2)
		returning ("id")
	`
)

type Series struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

func NewSeries(name, url string) (s *Series) {
	s = &Series{
		Name: name,
		Url:  url,
	}

	return s
}

func (this *Series) String() string {
	return fmt.Sprintf("%s", this.Name)
}

func (this *Series) Save(tx *pgx.Tx) (err error) {
	return tx.QueryRow(INSERT_SERIES, this.Name, this.Url).Scan(&this.ID)
}
