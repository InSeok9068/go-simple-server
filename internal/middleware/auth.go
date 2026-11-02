package middleware

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"simple-server/internal/config"
	"simple-server/internal/connection"
	"simple-server/pkg/util/authutil"

	resources "simple-server"

	firebase "firebase.google.com/go/v4"
	sqladapter "github.com/Blank-Xu/sql-adapter"
	"github.com/casbin/casbin/v3"
	casbinlog "github.com/casbin/casbin/v3/log"
	"github.com/casbin/casbin/v3/model"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
	_ "modernc.org/sqlite"
)

var App *firebase.App
var Enforcer *casbin.Enforcer

/* 인증 */
func InitFirebase() error {
	firebaseConfig := config.GetEnv("FIREBASE_CONFIG")
	if firebaseConfig == "" {
		return fmt.Errorf("FIREBASE_CONFIG 환경 변수가 비어 있습니다")
	}

	app, err := firebase.NewApp(
		context.Background(),
		nil,
		option.WithAuthCredentialsJSON(option.ServiceAccount, []byte(firebaseConfig)),
	)
	if err != nil {
		return fmt.Errorf("파이어베이스 초기화 실패: %w", err)
	}

	App = app

	return nil
}

func RegisterFirebaseAuthMiddleware(e *echo.Echo, ensureUserFn func(ctx context.Context, uid string) error) error {
	if err := InitFirebase(); err != nil {
		return err
	}

	store := sessions.NewCookieStore([]byte(config.GetEnv("SESSION_SECRET")))
	e.Use(session.Middleware(store))

	e.POST("/create-session", func(c echo.Context) error {
		ctx := c.Request().Context()
		auth, err := App.Auth(ctx)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "인증 서비스를 초기화하지 못했습니다.").SetInternal(err)
		}

		var req struct {
			Token string `json:"token"`
		}
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "로그인 요청 형식이 올바르지 않습니다.").SetInternal(err)
		}

		token, err := auth.VerifyIDToken(ctx, req.Token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "유효하지 않은 사용자입니다.").SetInternal(err)
		}

		if ensureUserFn != nil {
			if err := ensureUserFn(ctx, token.UID); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "로그인 사용자 정보를 준비하지 못했습니다.").SetInternal(err)
			}
		}

		if err := saveFirebaseSession(c, store, token.UID); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "로그인 세션을 저장하지 못했습니다.").SetInternal(err)
		}

		return c.NoContent(http.StatusNoContent)
	})

	return nil
}

func saveFirebaseSession(c echo.Context, store sessions.Store, uid string) error {
	// Firebase 토큰 검증 후에는 기존 쿠키를 읽지 않고 새 세션으로 교체한다.
	sess := sessions.NewSession(store, "session_v2")
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 180,
		HttpOnly: true,
		Secure:   config.IsProdEnv(),
		SameSite: http.SameSiteLaxMode,
	}
	sess.Values["uid"] = uid
	return sess.Save(c.Request(), c.Response())
}

/* 인증 */

/* 권한 */
func InitCasbin() error {
	db, err := connection.AppDBOpen()
	if err != nil {
		return fmt.Errorf("데이터베이스 연결 실패: %w", err)
	}
	adapter, err := sqladapter.NewAdapter(db, "sqlite", "")
	if err != nil {
		return fmt.Errorf("casbin adapter 초기화 실패: %w", err)
	}
	modelData, err := fs.ReadFile(resources.EmbeddedFiles, "model.conf")
	if err != nil {
		return fmt.Errorf("모델 파일 읽기 실패: %w", err)
	}
	m, err := model.NewModelFromString(string(modelData))
	if err != nil {
		return fmt.Errorf("casbin 모델 생성 실패: %w", err)
	}
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return fmt.Errorf("casbin enforcer 초기화 실패: %w", err)
	}

	Enforcer = enforcer

	return enforcer.LoadPolicy()
}

func RegisterCasbinMiddleware(e *echo.Group) error {
	if err := InitCasbin(); err != nil {
		return err
	}

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			uid, err := authutil.SessionUID(c)
			if err != nil {
				return err
			}

			obj := c.Path()
			act := c.Request().Method

			if config.IsDevEnv() {
				Enforcer.SetLogger(casbinlog.NewDefaultLogger())
			}

			ok, err := Enforcer.Enforce(uid, obj, act)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "권한 검사 실패")
			}
			if !ok {
				return echo.NewHTTPError(http.StatusForbidden, "권한이 없습니다.")
			}

			return next(c)
		}
	})

	return nil
}
