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
	"golang.org/x/net/html"
)

var lineCorrections = map[string]string{
	"\xe0": "a",
	"\xe8": "e",
	"\xe9": "e",
}

type Episode struct {
	ID      int64      `json:"id"`
	Title   string     `json:"title"`
	Number  int        `json:"number"`
	Url     string     `json:"url"`
	Airdate *time.Time `json:"airdate"`
	Script  []*Line    `json:"script"`
}

func (this *Episode) AddLine(line *Line) {
	if line.Speaker != "" {
		line.Episode = this

		for old, new := range lineCorrections {
			line.Line = strings.Replace(line.Line, old, new, -1)
		}

		this.Script = append(this.Script, line)
	}
}

func (this *Episode) ScriptString() string {
	buf := bytes.NewBuffer(nil)
	for _, line := range this.Script {
		buf.WriteString(line.String() + "\n")
	}

	return buf.String()
}

func (this *Episode) String() string {
	return fmt.Sprintf("%d. %s (%s)", this.Number, this.Title, this.Airdate.Format("Jan 2, 2006"))
}

func (this *Episode) Fetch() (err error) {
	var (
		resp *http.Response
		doc  *html.Node

		scriptText = bytes.NewBuffer(nil)
	)

	if resp, err = http.Get(this.Url); err != nil {
		return
	}
	defer resp.Body.Close()

	if doc, err = htmlquery.Parse(resp.Body); err != nil {
		return
	}

	for _, n := range htmlquery.Find(doc, "//td") {
		scriptText.WriteString(htmlquery.InnerText(n))
	}

	return this.Parse(scriptText)
}

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
				line = NewLine(strings.TrimRight(first, ":"), rest)
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
