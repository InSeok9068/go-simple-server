package services

import (
	"context"
	"net/http"
	"simple-server/internal/middleware"
	"simple-server/projects/deario/db"

	"github.com/labstack/echo/v4"
)

func EnsureUser(ctx context.Context, uid string) error {
	queries, _ := db.GetQueries(ctx)

	_, err := queries.GetUser(ctx, uid)
	if err != nil {
		auth, err := middleware.App.Auth(ctx)
		if err != nil {
			return err
		}
		user, err := auth.GetUser(ctx, uid)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "유효하지 않은 사용자입니다.")
		}
		_ = queries.CreateUser(ctx, db.CreateUserParams{
			Uid:   user.UID,
			Name:  user.DisplayName,
			Email: user.Email,
		})
	}

	middleware.Enforcer.AddRoleForUser(uid, "user")

	return nil
}
