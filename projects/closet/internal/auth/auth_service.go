package auth

import (
	"context"
	"fmt"

	"simple-server/internal/middleware"
	"simple-server/projects/closet/db"
)

// EnsureUser는 Firebase 인증 사용자 정보를 DB와 권한 시스템에 반영한다.
func EnsureUser(ctx context.Context, uid string) error {
	queries, err := db.GetQueries()
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
			return fmt.Errorf("인증 클라이언트 사용자 조회 실패: %w", err)
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
