package main

import (
	"flag"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"simple-server/cmd"
)

func main() {
	appMode := flag.Bool("app", false, "Run the admin server")
	flag.Parse()

	if *appMode {
		err := godotenv.Load()
		if err != nil {
			slog.Error("Failed to load .env file", "err", err)
		}
		slog.Info(os.Getenv("TEST_ENV"))

		cmd.AppServer()
	} else {
		cmd.AdminServer()
	}
}
