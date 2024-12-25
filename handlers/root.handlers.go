package handlers

import (
	"context"
	"database/sql"
	"log/slog"
	"simple-server/database"
	"simple-server/views"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func RootHandler(c echo.Context) error {
	return templ.Handler(views.Index()).Component.Render(c.Request().Context(), c.Response().Writer)
}

func HelloHandler(c echo.Context) error {
	return templ.Handler(views.Text("Hello")).Component.Render(c.Request().Context(), c.Response().Writer)
}

func BuyHandler(c echo.Context) error {
	return templ.Handler(views.Text("Buy")).Component.Render(c.Request().Context(), c.Response().Writer)
}

func ListHanlder(c echo.Context) error {
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
	return templ.Handler(views.Authors(authors)).Component.Render(c.Request().Context(), c.Response().Writer)
}
