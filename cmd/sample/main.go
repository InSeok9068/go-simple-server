package main

import (
	"simple-server/projects/sample/views"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return templ.Handler(views.Index("Sample")).Component.Render(c.Request().Context(), c.Response().Writer)
	})

	e.GET("/radio", func(c echo.Context) error {
		return templ.Handler(views.Radio()).Component.Render(c.Request().Context(), c.Response().Writer)
	})

	e.GET("/radio2", func(c echo.Context) error {
		return templ.Handler(views.Radio2()).Component.Render(c.Request().Context(), c.Response().Writer)
	})

	e.Logger.Fatal(e.Start(":8002"))
}
