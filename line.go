package treksum

import (
	"fmt"
	"strings"
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
