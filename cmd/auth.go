package cmd

import (
	"context"
	"firebase.google.com/go/v4"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/api/option"
	"log/slog"
	"os"
)

var App *firebase.App

func FirebaseInit() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Failed to load .env file", "err", err)
	}
	config := os.Getenv("FIREBASE_CONFIG")

	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON([]byte(config)))
	if err != nil {
		slog.Error("파이어베이스 초기화 실패", err)
	}

	App = app
}

func FirebaseAuth() middleware.KeyAuthConfig {
	return middleware.KeyAuthConfig{
		KeyLookup:  "header:Authorization",
		AuthScheme: "Bearer",
		Validator: func(key string, c echo.Context) (bool, error) {
			auth, err := App.Auth(context.Background())
			if err != nil {
				return false, nil
			}

			user, err := auth.VerifyIDToken(context.Background(), key)
			if err != nil {
				return false, nil
			}

			c.Set("user", user)

			return true, nil
		},
		Skipper: func(c echo.Context) bool {
			key := c.Request().Header.Get(echo.HeaderAuthorization)

			if key == "" {
				return true // Skip KeyAuth
			}

			return false
		},
	}
}
