package main

import (
	"time"
)

type Page struct {
	PageName       string
	WelcomeMessage string
	TimeInSQL      time.Duration
	TimeToRender   time.Duration
	NewsItems      []NewsItem
	Cfg            Config
}

func (p *Page) FinishSQL(t time.Time) {
	p.TimeInSQL = time.Since(t)
}

type PageData interface {
	FinishSQL(time.Time)
}

type QuotePage struct {
	Page
	Quotes []Quote
}

type SubmitPage struct {
	Page
	ID int64
}

type NewsItem struct {
	NewsText   string
	Author     string
	TimePosted time.Time
}

func NewPage(name string) Page {
	return Page{
		Cfg:      config,
		PageName: name,
	}
}

func (p *QuotePage) getQuotesFromDatabase(statement string, args ...interface{}) {

	var q []Quote

	SQLBeginning := time.Now()
	rows, err := db.Query(statement, args...)
	checkErr(err)
	p.TimeInSQL = time.Since(SQLBeginning)

	for rows.Next() {
		var body, notes string
		var id, rating uint
		err = rows.Scan(&id, &body, &notes, &rating)
		checkErr(err)
		newQuote := Quote{
			ID:     id,
			Text:   body,
			Notes:  notes,
			Rating: rating,
		}
		q = append(q, newQuote)
	}

	p.Quotes = q
}
