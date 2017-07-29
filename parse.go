package treksum

import (
	"net/http"
	"strings"
	"time"

	"github.com/antchfx/xquery/html"
	"go.uber.org/zap"
	"golang.org/x/net/html"
)

var (
	fmts = []string{
		"2 Jan, 2006",
		"2 Jan,2006",
		"2 Jan 2006",
		"2 Jan. 2006",
	}

	lineCorrections = map[string]string{
		"\xe0": "a",
		"\xe4": "a",
		"\xe7": "c",
		"\xe8": "e",
		"\xe9": "e",
		"\xea": "e",
		"\xef": "i",
		"\x91": "^",
		"\x92": "'",
		"\xdf": "ss",
	}
)

func CleanUnicode(in string) string {
	for old, new := range lineCorrections {
		in = strings.Replace(in, old, new, -1)
	}

	return in
}

func SwitchPage(url, page string) string {
	return url[:strings.LastIndex(url, "/")+1] + page
}

func FindEpisodesLink(url string) (next string, err error) {
	var (
		resp *http.Response
		doc  *html.Node
	)

	if resp, err = http.Get(url); err != nil {
		return
	}
	defer resp.Body.Close()

	if doc, err = htmlquery.Parse(resp.Body); err != nil {
		return
	}

	a := htmlquery.FindOne(doc, "//a[contains(., 'Episode')]")
	url = SwitchPage(url, htmlquery.SelectAttr(a, "href"))

	return url, nil
}

func ParseEpisodeList(log *zap.Logger, series *Series) (episodes []*Episode, err error) {
	var (
		resp    *http.Response
		doc     *html.Node
		episode *Episode

		url    string
		season int
		epNum  int
	)

	log.Info("finding episodes")
	if url, err = FindEpisodesLink(series.Url); err != nil {
		return
	}

	if resp, err = http.Get(url); err != nil {
		return
	}
	defer resp.Body.Close()

	if doc, err = htmlquery.Parse(resp.Body); err != nil {
		return
	}

	path := "//td/table"
	if strings.Contains(url, "Voyager") {
		path = "//div/table"
	}

	for _, table := range htmlquery.Find(doc, path) {
		season++
		epNum = 0

		for _, td := range htmlquery.Find(table, "//tr/td") {
			text := strings.TrimSpace(htmlquery.InnerText(td))
			text = strings.Replace(text, "\n", " ", -1)

			// ignore table headers
			if strings.Contains(text, "Episode Name") {
				continue
			}

			// make sure there's a link to the episode transcript
			a := htmlquery.FindOne(td, "//a/@href")
			if a != nil && episode == nil {
				epNum++
				episode = &Episode{
					Series:  series,
					Season:  season,
					Episode: epNum,
					Title:   text,
					Url:     SwitchPage(series.Url, htmlquery.SelectAttr(a, "href")),
				}
				episode.Log = log.With(zap.String("episode", episode.GetAbbrev()))
				continue
			}

			if episode != nil && episode.Airdate == nil {
				// fix some dates
				text = strings.Replace(text, "Sept", "Sep", -1)

				// try matching against a series of date formats
				for _, f := range fmts {
					if t, err := time.Parse(f, text); err == nil {
						episode.Airdate = &t

						// found a matching date format
						break
					}
				}

				// skip if we don't have a valid airdate
				if episode.Airdate == nil {
					continue
				}

				episodes = append(episodes, episode)
				episode = nil
			}

		}
	}

	return episodes, nil
}
