package services

import (
	"context"
	"fmt"
	"net/http"
	"simple-server/internal/middleware"
	"simple-server/projects/deario/db"

	"github.com/labstack/echo/v4"
)

func EnsureUser(ctx context.Context, uid string) error {
	queries, err := db.GetQueries(ctx)
	if err != nil {
		return fmt.Errorf("쿼리 로드 실패: %w", err)
	}

	if _, err := queries.GetUser(ctx, uid); err != nil {
		auth, err := middleware.App.Auth(ctx)
		if err != nil {
			return fmt.Errorf("인증 클라이언트 생성 실패: %w", err)
		}

		user, err := auth.GetUser(ctx, uid)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "유효하지 않은 사용자입니다.")
		}

		if err := queries.CreateUser(ctx, db.CreateUserParams{
			Uid:   user.UID,
			Name:  user.DisplayName,
			Email: user.Email,
		}); err != nil {
			return fmt.Errorf("사용자 생성 실패: %w", err)
		}
	}

	if _, err := middleware.Enforcer.AddRoleForUser(uid, "user"); err != nil {
		return fmt.Errorf("역할 추가 실패: %w", err)
	}

	return nil
}
