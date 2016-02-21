package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	QDBName   string
	SiteRoot  string
	ShowQueue bool
}

var (
	config Config
	db     *sql.DB
)

func init() {
	database, err := sql.Open("sqlite3", "./qdb.db")
	checkErr(err)

	_, err = toml.DecodeFile("./config.toml", &config)
	checkErr(err)

	fmt.Println(config)

	db = database
}

func HomeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	now := time.Now()
	fmt.Println("Handling homepage.")

	p := Page{
		Cfg:            config,
		PageName:       "Home",
		WelcomeMessage: "Hello!",
		NewsItems: []NewsItem{
			{
				NewsText:   "This is a first attempt.",
				Author:     "Liru",
				TimePosted: time.Now(),
			},
			{
				NewsText:   "This is a zeroth attempt.",
				Author:     "Liru",
				TimePosted: time.Now().Add(-1 * time.Hour),
			},
		},
	}

	t, err := template.ParseFiles("tmpl/base.tmpl", "tmpl/home.tmpl")
	if err != nil {
		panic(err)
	}

	ttr := time.Since(now)
	p.TimeToRender = ttr
	err = t.Execute(w, p)
	checkErr(err)

	fmt.Println("Took", time.Since(now), "to run.")

}

func QuoteHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	now := time.Now()
	fmt.Println("Handling single quote.")
	quoteID := ps.ByName("quote")

	p := QuotePage{
		Page: NewPage(quoteID),
	}

	myFuncMap := template.FuncMap{
		"nl2br": nl2br,
	}

	t, err := template.New("base.tmpl").Funcs(myFuncMap).ParseFiles("tmpl/base.tmpl", "tmpl/quotes.tmpl")
	if err != nil {
		panic(err)
	}

	SQLBeginning := time.Now()
	p.Quotes = GetQuote(quoteID)
	p.TimeInSQL = time.Since(SQLBeginning)

	ttr := time.Since(now) - p.TimeInSQL
	p.TimeToRender = ttr
	err = t.Execute(w, p)
	checkErr(err)
	fmt.Println("Took", time.Since(now), "to run.")
}

func SearchHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	Incomplete(w)
}

func FlagHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	Incomplete(w)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	Incomplete(w)
}

func LatestHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	now := time.Now()
	fmt.Println("Handling latest.")

	p := QuotePage{
		Page: NewPage("Latest"),
	}

	myFuncMap := template.FuncMap{
		"nl2br": nl2br,
	}

	t, err := template.New("base.tmpl").Funcs(myFuncMap).ParseFiles("tmpl/base.tmpl", "tmpl/quotes.tmpl")
	checkErr(err)

	SQLBeginning := time.Now()
	p.Quotes = Latest()
	p.TimeInSQL = time.Since(SQLBeginning)

	ttr := time.Since(now) - p.TimeInSQL
	p.TimeToRender = ttr
	err = t.Execute(w, p)
	checkErr(err)
	fmt.Println("Took", time.Since(now), "to run.")
}

func TopHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	Incomplete(w)
}

func SubmissionHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	now := time.Now()
	fmt.Println("Handling submission.")

	p := SubmitPage{
		Page: NewPage("Submit"),
	}

	var SQLBeginning time.Time
	if r.Method == "POST" {
		err := r.ParseForm()
		checkErr(err)

		fmt.Println(r.Form)

		if r.Form.Get("quote") != "" {
			SQLBeginning = time.Now()
			id, err := AddQuote(r.Form)
			checkErr(err)
			p.TimeInSQL = time.Since(SQLBeginning)
			p.ID = id
		}
	}

	t, err := template.ParseFiles("tmpl/base.tmpl", "tmpl/submit.tmpl")
	checkErr(err)

	p.TimeToRender = time.Since(now) - p.TimeInSQL
	err = t.Execute(w, p)
	checkErr(err)
	fmt.Println("Took", time.Since(now), "to run.")
}

func BrowseHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	now := time.Now()
	fmt.Println("Handling browse.")

	p := QuotePage{
		Page: NewPage("Browse"),
	}

	myFuncMap := template.FuncMap{
		"nl2br": nl2br,
	}

	t, err := template.New("base.tmpl").Funcs(myFuncMap).ParseFiles("tmpl/base.tmpl", "tmpl/quotes.tmpl")
	checkErr(err)

	var id int

	sid := r.URL.Query().Get("p")

	if sid == "" {
		id = 1
	} else {
		id2, err := strconv.Atoi(sid)
		checkErr(err)
		id = id2
	}

	SQLBeginning := time.Now()
	p.Quotes = Browse(id - 1)
	p.TimeInSQL = time.Since(SQLBeginning)

	ttr := time.Since(now) - p.TimeInSQL
	p.TimeToRender = ttr
	err = t.Execute(w, p)
	checkErr(err)
	fmt.Println("Took", time.Since(now), "to run.")
}

func Incomplete(w http.ResponseWriter) {
	p := NewPage("INCOMPLETE")
	t, err := template.ParseFiles("tmpl/base.tmpl", "tmpl/todo.tmpl")
	checkErr(err)

	err = t.Execute(w, p)
	checkErr(err)

}

func main() {
	router := httprouter.New()
	router.GET("/", HomeHandler)
	router.GET("/q/:quote", QuoteHandler)
	router.GET("/search/:query", SearchHandler)
	router.GET("/q/:quote/flag", FlagHandler)
	router.GET("/q/:quote/delete", DeleteHandler)
	router.GET("/latest", LatestHandler)
	router.GET("/top", TopHandler)
	router.GET("/browse", BrowseHandler)

	router.GET("/submit", SubmissionHandler)
	router.POST("/submit", SubmissionHandler)

	router.ServeFiles("/css/*filepath", http.Dir("./public/css/"))

	fmt.Println("Starting server.")

	log.Fatal(http.ListenAndServe(":26362", router))
}
