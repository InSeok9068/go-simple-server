package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"simple-server/internal/config"
	"simple-server/internal/connection"
	"simple-server/pkg/util/authutil"

	firebase "firebase.google.com/go/v4"
	sqladapter "github.com/Blank-Xu/sql-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
	_ "modernc.org/sqlite"
)

var App *firebase.App

/* 인증 */
func InitFirebase() error {
	config := config.EnvMap["FIREBASE_CONFIG"]

	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON([]byte(config)))
	if err != nil {
		slog.Error("파이어베이스 초기화 실패", "error", err.Error())
		return err
	}

	App = app

	return nil
}

func RegisterFirebaseAuthMiddleware(e *echo.Echo) error {
	err := InitFirebase()
	if err != nil {
		return err
	}

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

		slog.Info("data", "token", data["token"].(string))

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
			Secure:   config.IsProdEnv(),
			SameSite: http.SameSiteLaxMode,
		}
		sess.Values["uid"] = user.UID

		if err := sess.Save(c.Request(), c.Response()); err != nil {
			return err
		}

		return c.NoContent(http.StatusOK)
	})

	return nil
}

/* 인증 */

/* 권한 */
func InitCasbin() (*casbin.Enforcer, error) {
	serviceName := os.Getenv("SERVICE_NAME")

	db, _ := connection.AppDBOpen()
	adapter, err := sqladapter.NewAdapter(db, "sqlite", "")
	if err != nil {
		slog.Error("casbin adapter 초기화 실패", "error", err.Error())
		return nil, err
	}
	enforcer, err := casbin.NewEnforcer(fmt.Sprintf("./projects/%s/model.conf", serviceName), adapter)
	if err != nil {
		slog.Error("casbin enforcer 초기화 실패", "error", err.Error())
		return nil, err
	}
	return enforcer, enforcer.LoadPolicy()
}

func RegisterCasbinMiddleware(e *echo.Echo) {
	enforcer, err := InitCasbin()
	if err != nil {
		return
	}

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			uid, err := authutil.SessionUID(c)
			if err != nil {
				return err
			}

			obj := c.Path()           // 요청 경로
			act := c.Request().Method // 요청 메서드 (GET/POST 등)

			ok, err := enforcer.Enforce(uid, obj, act)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "권한 검사 실패")
			}
			if !ok {
				return echo.NewHTTPError(http.StatusForbidden, "권한이 없습니다.")
			}

			return next(c)
		}
	})

}
