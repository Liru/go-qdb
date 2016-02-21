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

func GetQuote(sid string) []Quote {

	q := make([]Quote, 0)

	rows, err := db.Query("SELECT id,body,notes,rating FROM quotes WHERE id = ?", sid)
	checkErr(err)

	for rows.Next() {
		var body, notes string
		var id, rating uint
		err = rows.Scan(&id, &body, &notes, &rating)
		checkErr(err)
		fmt.Println(body)
		fmt.Println(notes)
		q = append(q, Quote{ID: id, Text: body, Notes: notes, Rating: rating})
	}

	return q
}

func Browse(page int) []Quote {
	q := make([]Quote, 0)

	rows, err := db.Query("SELECT id,body,notes,rating FROM quotes ORDER BY id ASC LIMIT 20 OFFSET ?", page*20)
	checkErr(err)

	fmt.Println("iterating")

	for rows.Next() {
		var body, notes string
		var id, rating uint
		err = rows.Scan(&id, &body, &notes, &rating)
		checkErr(err)
		fmt.Println(body)
		fmt.Println(notes)
		newQuote := Quote{ID: id, Text: body, Notes: notes, Rating: rating}
		fmt.Println(newQuote)
		q = append(q, newQuote)
	}
	fmt.Println("returning browse")

	return q
}

func Latest() []Quote {
	q := make([]Quote, 0)

	rows, err := db.Query("SELECT id,body,notes,rating FROM quotes ORDER BY id DESC LIMIT 20")
	checkErr(err)

	for rows.Next() {
		var body, notes string
		var id, rating uint
		err = rows.Scan(&id, &body, &notes, &rating)
		checkErr(err)
		fmt.Println(body)
		fmt.Println(notes)
		newQuote := Quote{ID: id, Text: body, Notes: notes, Rating: rating}
		fmt.Println(newQuote)
		q = append(q, newQuote)
	}
	return q
}

func Top() []Quote {
	q := make([]Quote, 0)

	rows, err := db.Query("SELECT id,body,notes,rating FROM quotes ORDER BY rating DESC LIMIT 20")
	checkErr(err)

	for rows.Next() {
		var body, notes string
		var id, rating uint
		err = rows.Scan(&id, &body, &notes, &rating)
		checkErr(err)
		fmt.Println(body)
		fmt.Println(notes)
		newQuote := Quote{ID: id, Text: body, Notes: notes, Rating: rating}
		fmt.Println(newQuote)
		q = append(q, newQuote)
	}
	return q
}

func (q *Quote) String() string {
	return fmt.Sprint("%d: %s", q.ID, q.Text)
}
