package treksum

import (
	"net/http"
	"strconv"
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

func ParseEpisodeList() (episodes []*Episode, err error) {
	var (
		resp *http.Response
		doc  *html.Node

		episode *Episode
	)

	if resp, err = http.Get(url); err != nil {
		return
	}
	defer resp.Body.Close()

	if doc, err = htmlquery.Parse(resp.Body); err != nil {
		return
	}

	for _, n := range htmlquery.Find(doc, "//td") {
		text := strings.TrimSpace(htmlquery.InnerText(n))
		text = strings.Replace(text, "\n", " ", -1)

		// ignore table headers
		if strings.Contains(text, "Episode Name") {
			continue
		}

		// make sure there's a link to the episode transcript
		a := htmlquery.FindOne(n, "//a/@href")
		if a != nil && episode == nil {
			episode = &Episode{
				Title: text,
				Url:   "http://chakoteya.net/NextGen/" + htmlquery.SelectAttr(a, "href"),
			}
			continue
		}

		if episode != nil {
			if episode.Number == 0 {
				if text == "101 + 102" {
					text = "101"
				}

				// skip if we can't conver the episode number to an int
				if episode.Number, err = strconv.Atoi(text); err != nil {
					continue
				}
			} else if episode.Airdate == nil {
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
