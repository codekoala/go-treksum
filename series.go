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
			unique("name")
		)
	`

	INSERT_SERIES = `
		insert into "series"
			("name")
		values
			($1)
		returning ("id")
	`
)

type Series struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func NewSeries(name string) (s *Series) {
	s = &Series{
		Name: name,
	}

	return s
}

func (this *Series) String() string {
	return fmt.Sprintf("%s", this.Name)
}

func (this *Series) Save(pool *pgx.ConnPool) (err error) {
	return pool.QueryRow(INSERT_SERIES, this.Name).Scan(&this.ID)
}
