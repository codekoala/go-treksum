package treksum

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx"
	"go.uber.org/zap"
)

const (
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

	INSERT_LINE = `
		insert into "line"
			("episode_id", "speaker_id", "line")
		values
			($1, $2, $3)
		returning ("id")
	`
)

type Line struct {
	ID      int64    `json:"id"`
	Episode *Episode `json:"episode,omitempty"`
	Speaker string   `json:"speaker"`
	Line    string   `json:"line"`

	log *zap.Logger
}

func NewLine(log *zap.Logger, speaker, text string) (l *Line) {
	l = &Line{
		Speaker: speaker,

		log: log,
	}
	l.AddText(text)

	return l
}

// AddText appends text to existing text in a line from the transcript.
func (this *Line) AddText(text string) {
	this.Line = strings.TrimSpace(fmt.Sprintf("%s %s", this.Line, text))
}

func (this *Line) String() string {
	return fmt.Sprintf("%s: %s", this.Speaker, this.Line)
}

// Save persists the line to the database;
func (this *Line) Save(tx *pgx.Tx) (err error) {
	speaker := NewSpeaker(this.Episode.Series, this.Speaker)
	if err = speaker.Save(tx); err != nil {
		this.log.Warn("failed to save speaker", zap.Error(err))
	}

	err = tx.QueryRow(INSERT_LINE, this.Episode.ID, speaker.ID, this.Line).Scan(&this.ID)
	if err != nil {
		this.log.Warn("failed to save line", zap.Error(err))
	}

	return
}
