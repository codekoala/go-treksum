package treksum

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/antchfx/xquery/html"
	"github.com/jackc/pgx"
	"go.uber.org/zap"
	"golang.org/x/net/html"
)

const (
	CREATE_EPISODE_TABLE = `
		create table "episode" (
			"id" serial primary key,
			"series_id" int,
			"season" int,
			"episode" int,
			"title" text,
			"url" text,
			"airdate" date,
			foreign key ("series_id") references "series" ("id"),
			unique ("series_id", "season", "episode")
		)
	`

	INSERT_EPISODE = `
		insert into "episode"
			("series_id", "season", "episode", "title", "url", "airdate")
		values
			($1, $2, $3, $4, $5, $6::date)
		returning ("id")
	`
)

type Episode struct {
	ID      int64      `json:"id"`
	Series  *Series    `json:"series"`
	Season  int        `json:"season"`
	Episode int        `json:"episode"`
	Title   string     `json:"title"`
	Url     string     `json:"url"`
	Airdate *time.Time `json:"airdate"`
	Script  []*Line    `json:"script"`

	Log *zap.Logger
}

// AddLine adds a new line from the episode's script, cleaning up any problematic characters in the process.
func (this *Episode) AddLine(line *Line) {
	if line.Speaker != "" {
		line.Episode = this
		line.Line = CleanUnicode(line.Line)
		this.Script = append(this.Script, line)
	}
}

// ScriptString combines all lines into a single string for the entire episode's script.
func (this *Episode) ScriptString() string {
	buf := bytes.NewBuffer(nil)
	for _, line := range this.Script {
		buf.WriteString(line.String() + "\n")
	}

	return buf.String()
}

// GetAbbrev returns the season and episode in the form of S00E00.
func (this *Episode) GetAbbrev() string {
	return fmt.Sprintf("S%02dE%02d", this.Season, this.Episode)
}

func (this *Episode) String() string {
	return fmt.Sprintf("%s %s (%s)", this.GetAbbrev(), this.Title, this.Airdate.Format("Jan 2, 2006"))
}

// Fetch scrapes the episode's transcript.
func (this *Episode) Fetch() (err error) {
	var (
		resp *http.Response
		doc  *html.Node

		scriptText = bytes.NewBuffer(nil)
	)

	this.Log.Info("fetching episode")
	if resp, err = http.Get(this.Url); err != nil {
		return
	}
	defer resp.Body.Close()

	if doc, err = htmlquery.Parse(resp.Body); err != nil {
		return
	}

	// add all regular text to a single buffer that we can parse later
	for _, n := range htmlquery.Find(doc, "//td") {
		scriptText.WriteString(htmlquery.InnerText(n))
	}

	return this.Parse(scriptText)
}

// Parse extracts individual lines from the episode's transcript.
func (this *Episode) Parse(scriptText io.Reader) (err error) {
	var (
		first string
		rest  string
		aside bool

		line    = new(Line)
		scanner = bufio.NewScanner(scriptText)
	)

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		} else if text[0] == '(' || text[0] == '[' {
			aside = true
			continue
		}

		// get the first word and the rest of the line
		toks := strings.SplitN(text, ":", 2)
		if len(toks) == 2 {
			first, rest = toks[0], toks[1]

			if strings.Contains(first, "[") && strings.Contains(first, "]") {
				first = strings.TrimSpace(first[:strings.Index(first, "[")])
			}

			// see if the first word looks like it's a character speaking
			if strings.ToUpper(first) == first {
				this.AddLine(line)
				aside = false
				line = NewLine(this.Log, strings.TrimRight(first, ":"), rest)
				continue
			}
		}

		if !aside {
			line.AddText(text)
		}
	}
	this.AddLine(line)

	if err = scanner.Err(); err != nil {
		return
	}

	return nil
}

// Save persists the episode and its transcript to the database.
func (this *Episode) Save(tx *pgx.Tx) (err error) {
	this.Log.Info("inserting episode")
	err = tx.QueryRow(INSERT_EPISODE, this.Series.ID, this.Season, this.Episode, this.Title, this.Url, this.Airdate).Scan(&this.ID)
	if err != nil {
		return
	}

	for _, line := range this.Script {
		if err = line.Save(tx); err != nil {
			return
		}
	}

	return nil
}
