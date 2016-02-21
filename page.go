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
