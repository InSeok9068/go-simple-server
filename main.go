package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"simple-server/database"
	"simple-server/views"

	"github.com/a-h/templ"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(views.Index()).ServeHTTP(w, r)
	})
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(views.Text("Hello")).ServeHTTP(w, r)
	})
	http.HandleFunc("/buy", func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(views.Text("Buy")).ServeHTTP(w, r)
	})
	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		db, err := sql.Open("sqlite3", "file:./database/data.db")
		if err != nil {
			slog.Error(err.Error())
		}

		queries := database.New(db)
		authors, err := queries.ListAuthors(ctx)
		if err != nil {
			slog.Error(err.Error())
		}
		templ.Handler(views.Authors(authors)).ServeHTTP(w, r)
	})

	// 서버 실행
	slog.Error(http.ListenAndServe(":8000", nil).Error())
}
