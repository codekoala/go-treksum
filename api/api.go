package api

import (
	"net/http"

	"github.com/jackc/pgx"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

var (
	log  *zap.Logger
	pool *pgx.ConnPool
)

type Quote struct {
	Series  string `json:"series"`
	Title   string `json:"title"`
	Season  int    `json:"season"`
	Episode int    `json:"episode"`
	Airdate string `json:"airdate"`
	Speaker string `json:"speaker"`
	Line    string `json:"line"`
}

func Setup(l *zap.Logger, p *pgx.ConnPool) (app *echo.Echo) {
	log = l
	pool = p

	app = echo.New()

	app.GET("/api/v1/random", RandomQuote)

	return
}

func RandomQuote(c echo.Context) (err error) {
	var quote Quote

	err = pool.QueryRow(`
		select
			s.name as series,
			e.title as title,
			e.season,
			e.episode,
			e.airdate::text,
			k.name as speaker,
			l.line
		from line l
		join episode e
		  on (e.id = l.episode_id)
		join speaker k
		  on (k.id = l.speaker_id and k.series_id = e.series_id)
		join series s
		  on (s.id = e.series_id)
		order by random()
		limit 1
	`).Scan(
		&quote.Series,
		&quote.Title,
		&quote.Season,
		&quote.Episode,
		&quote.Airdate,
		&quote.Speaker,
		&quote.Line,
	)
	if err != nil {
		log.Warn("failed to fetch quote", zap.Error(err))
		return
	}

	return c.JSON(http.StatusOK, quote)
}
