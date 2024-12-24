package main

import (
	"log"
	"net/http"
	"simple-server/views"

	"github.com/a-h/templ"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(views.Index()).ServeHTTP(w, r)
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
