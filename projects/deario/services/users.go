package services

import (
	"context"
	"log/slog"
	"net/http"
	"simple-server/internal/middleware"
	"simple-server/projects/deario/db"

	"github.com/labstack/echo/v4"
)

func EnsureUser(ctx context.Context, uid string) error {
	slog.Info("EnsureUser called", "uid", uid)

	queries, _ := db.GetQueries(ctx)

	_, err := queries.GetUser(ctx, uid)
	if err != nil {
		slog.Info("User not found in database, creating new user", "uid", uid, "error", err)

		auth, err := middleware.App.Auth(ctx)
		if err != nil {
			slog.Error("Failed to get auth client", "error", err)
			return err
		}

		user, err := auth.GetUser(ctx, uid)
		if err != nil {
			slog.Error("Failed to get user from Firebase", "uid", uid, "error", err)
			return echo.NewHTTPError(http.StatusUnauthorized, "유효하지 않은 사용자입니다.")
		}

		err = queries.CreateUser(ctx, db.CreateUserParams{
			Uid:   user.UID,
			Name:  user.DisplayName,
			Email: user.Email,
		})
		if err != nil {
			slog.Error("Failed to create user in database", "uid", uid, "error", err)
			return err
		}
		slog.Info("Successfully created new user", "uid", uid)
	}

	if _, err := middleware.Enforcer.AddRoleForUser(uid, "user"); err != nil {
		slog.Error("Failed to add role for user", "uid", uid, "error", err)
		return err
	}

	return nil
}
