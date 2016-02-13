package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Page struct {
	QDBName        string
	PageName       string
	WelcomeMessage string
	NewsItems      []NewsItem
}

type NewsItem struct {
	NewsText   string
	Author     string
	TimePosted time.Time
}

func HomeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func QuoteHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func SearchHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func FlagHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func DeleteHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func LatestHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func TopHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func SubmissionHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func SubmissionPostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

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

	router.GET("/submit", SubmissionHandler)
	router.POST("/submit", SubmissionPostHandler)

	fmt.Println("Starting server.")

	log.Fatal(http.ListenAndServe(":26362", router))
}
