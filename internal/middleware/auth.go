package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"simple-server/internal/config"

	firebase "firebase.google.com/go/v4"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
)

var App *firebase.App

func FirebaseInit() {
	config := config.EnvMap["FIREBASE_CONFIG"]

	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON([]byte(config)))
	if err != nil {
		slog.Error("파이어베이스 초기화 실패", "error", err.Error())
	}

	App = app
}

func RegisterFirebaseAuthMiddleware(e *echo.Echo) {
	FirebaseInit()

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	e.POST("/create-session", func(c echo.Context) error {
		ctx := c.Request().Context()
		auth, err := App.Auth(ctx)
		if err != nil {
			return err
		}

		var data map[string]interface{}
		if err := c.Bind(&data); err != nil {
			return err
		}

		slog.Info("data", data["token"].(string))

		user, err := auth.VerifyIDToken(ctx, data["token"].(string))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "유효하지 않은 사용자입니다.")
		}

		sess, err := session.Get("session", c)
		if err != nil {
			return err
		}
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
		}
		sess.Values["uid"] = user.UID

		if err := sess.Save(c.Request(), c.Response()); err != nil {
			return err
		}

		return c.NoContent(http.StatusOK)
	})
}
