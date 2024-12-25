package main

import (
	"simple-server/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", handlers.RootHandler)
	e.GET("/hello", handlers.HelloHandler)
	e.GET("/buy", handlers.BuyHandler)
	e.GET("/list", handlers.ListHanlder)

	e.Logger.Fatal(e.Start(":8000"))
}
