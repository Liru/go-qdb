package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HomeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
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
	router.GET("/:quote", QuoteHandler)
	router.GET("/search/:query", SearchHandler)
	router.GET("/flag/:quote", FlagHandler)
	router.GET("/delete/:quote", DeleteHandler)
	router.GET("/latest", LatestHandler)
	router.GET("/top", TopHandler)

	router.GET("/submit", SubmissionHandler)
	router.POST("/submit", SubmissionPostHandler)

	log.Fatal(http.ListenAndServe(":26362", router))
}
