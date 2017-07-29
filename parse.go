package treksum

import (
	"net/http"
	"strings"
	"time"

	"github.com/antchfx/xquery/html"
	"golang.org/x/net/html"
)

var (
	url  = "http://chakoteya.net/NextGen/episodes.htm"
	fmts = []string{
		"2 Jan, 2006",
		"2 Jan,2006",
		"2 Jan 2006",
		"2 Jan. 2006",
	}
)

func ParseEpisodeList(series *Series) (episodes []*Episode, err error) {
	var (
		resp    *http.Response
		doc     *html.Node
		episode *Episode

		season int
	)

	if resp, err = http.Get(url); err != nil {
		return
	}
	defer resp.Body.Close()

	if doc, err = htmlquery.Parse(resp.Body); err != nil {
		return
	}

	for _, table := range htmlquery.Find(doc, "//td/table") {
		var episodeNum int
		season++

		for _, tr := range htmlquery.Find(table, "//tr") {
			for _, td := range htmlquery.Find(tr, "//td") {
				text := strings.TrimSpace(htmlquery.InnerText(td))
				text = strings.Replace(text, "\n", " ", -1)

				// ignore table headers
				if strings.Contains(text, "Episode Name") {
					continue
				}

				// make sure there's a link to the episode transcript
				a := htmlquery.FindOne(td, "//a/@href")
				if a != nil && episode == nil {
					episodeNum++
					episode = &Episode{
						Series:  series,
						Season:  season,
						Episode: episodeNum,
						Title:   text,
						Url:     "http://chakoteya.net/NextGen/" + htmlquery.SelectAttr(a, "href"),
					}
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
	}

	return episodes, nil
}
