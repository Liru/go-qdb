package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

type QuoteHandle func(http.ResponseWriter, *http.Request, httprouter.Params, *QuotePage)

type Config struct {
	QDBName   string
	SiteRoot  string
	ShowQueue bool
}

var (
	config      Config
	db          *sql.DB
	tmplFuncMap = template.FuncMap{
		"nl2br": nl2br,
	}
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
}

func QuoteHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, p *QuotePage) {
	quoteID := ps.ByName("quote")
	p.Page = NewPage("Quote #" + quoteID)

	p.getQuotesFromDatabase("SELECT id,body,notes,rating FROM quotes WHERE id = ?", quoteID)
}

func SearchHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params, p *QuotePage) {

	if r.URL.RawQuery != "" {

		terms := strings.Split(r.URL.Query().Get("q"), " ")
		query := "SELECT id,body,notes,rating FROM quotes WHERE 1=1"

		for i := 0; i < len(terms); i++ {
			// This took WAY too long for what it was.
			// Note to future self: Go doesn't like '%?%'. It takes it literally and
			// ignores the question mark as a binding parameter.
			query += " AND body LIKE '%' || ? || '%'"
		}
		query += " ORDER BY id DESC"

		// We have to cast `terms` to []interface{} because Go sucks
		args := make([]interface{}, len(terms))

		for i := range terms {
			args[i] = terms[i]
		}
		p.getQuotesFromDatabase(query, args...)

	}
}

func FlagHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	Incomplete(w)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	Incomplete(w)
}

func LatestHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params, p *QuotePage) {
	p.Page = NewPage("Latest Quotes")

	p.getQuotesFromDatabase("SELECT id,body,notes,rating FROM quotes ORDER BY id DESC LIMIT 20")
}

func TopHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params, p *QuotePage) {
	p.Page = NewPage("Top Quotes")

	p.getQuotesFromDatabase("SELECT id,body,notes,rating FROM quotes ORDER BY rating DESC LIMIT 20")
}

func SubmissionHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	now := time.Now()
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
}

func BrowseHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, p *QuotePage) {
	p.Page = NewPage("Browse Quotes")

	sid := r.URL.Query().Get("p")
	page := 0
	if sid != "" {
		id, err := strconv.Atoi(sid)
		checkErr(err)
		if id >= 0 {
			page = id
		}
	}

	p.getQuotesFromDatabase("SELECT id,body,notes,rating FROM quotes ORDER BY id ASC LIMIT 20 OFFSET ?", page*20)

}

/////////////////////////////////////

func Incomplete(w http.ResponseWriter) {
	p := NewPage("INCOMPLETE")
	t, err := template.ParseFiles("tmpl/base.tmpl", "tmpl/todo.tmpl")
	checkErr(err)

	err = t.Execute(w, p)
	checkErr(err)

}

var _ = DebugTime

func DebugTime(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		now := time.Now()
		fn(w, r, ps)
		fmt.Println("Took", time.Since(now), "to run.")
	}
}

func QuoteWrapper(handle QuoteHandle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		now := time.Now()

		p := &QuotePage{}

		t, err := template.New("base.tmpl").Funcs(tmplFuncMap).ParseFiles("tmpl/base.tmpl", "tmpl/quotes.tmpl")
		checkErr(err)

		handle(w, r, ps, p)

		ttr := time.Since(now) - p.TimeInSQL
		p.TimeToRender = ttr
		err = t.Execute(w, p)
		checkErr(err)

	}
}

func main() {
	router := httprouter.New()
	router.GET("/", HomeHandler)
	router.GET("/q/:quote", QuoteWrapper(QuoteHandler))
	router.GET("/search", QuoteWrapper(SearchHandler))
	router.GET("/q/:quote/flag", FlagHandler)
	router.GET("/q/:quote/delete", DeleteHandler)
	router.GET("/latest", QuoteWrapper(LatestHandler))
	router.GET("/top", QuoteWrapper(TopHandler))
	router.GET("/browse", QuoteWrapper(BrowseHandler))

	router.GET("/submit", SubmissionHandler)
	router.POST("/submit", SubmissionHandler)

	router.ServeFiles("/css/*filepath", http.Dir("./public/css/"))

	fmt.Println("Starting server.")

	log.Fatal(http.ListenAndServe(":26362", router))
}
