package main

import (
	"fmt"
	"html/template"
	"net/url"
	"strings"
	"time"
)

// A Quote stores all the info needed to render a quote on the site.
type Quote struct {
	ID         uint
	Text       string
	Notes      string
	CreatedAt  int64
	Rating     uint
	Up         uint
	Down       uint
	TotalVotes uint
	Score      uint // Unseen. This is generated for sorting.
}

func NewQuote(text string) Quote {
	return Quote{
		Text:      text,
		CreatedAt: time.Now().Unix(),
	}
}

func nl2br(text string) template.HTML {
	return template.HTML(strings.Replace(template.HTMLEscapeString(text), "\n", "<br>", -1))
}

func AddQuote(v url.Values) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO quotes (body, notes) VALUES (?, ?)")
	checkErr(err)

	result, err := stmt.Exec(v.Get("quote"), v.Get("comment"))
	checkErr(err)

	return result.LastInsertId()
}

func (q *Quote) String() string {
	return fmt.Sprint("%d: %s", q.ID, q.Text)
}
