package services

import (
	"log/slog"
	"simple-server/projects/homepage/db"
)

func GetAuthors() ([]db.Author, error) {
	queries, ctx := db.DbQueries()
	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		slog.Error("저자 목록 조회 오류", "error", err.Error())
		return nil, err
	}
	return authors, nil
}

func GetAuthor(id string) (db.Author, error) {
	queries, ctx := db.DbQueries()
	author, err := queries.GetAuthor(ctx, id)
	if err != nil {
		slog.Error("저자 조회 오류", "error", err.Error())
		return db.Author{}, err
	}
	return author, nil
}

func CreateAuthor(name string, bio string) (db.Author, error) {
	queries, ctx := db.DbQueries()
	author, err := queries.CreateAuthor(ctx, db.CreateAuthorParams{
		Name: name,
		Bio:  bio,
	})
	if err != nil {
		slog.Error("저자 등록 오류", "error", err.Error())
		return db.Author{}, err
	}
	return author, nil
}

func UpdateAuthor(id string, name string, bio string) (db.Author, error) {
	queries, ctx := db.DbQueries()
	author, err := queries.UpdateAuthor(ctx, db.UpdateAuthorParams{
		ID:   id,
		Name: name,
		Bio:  bio,
	})
	if err != nil {
		slog.Error("저자 수정 오류", "error", err.Error())
		return db.Author{}, err
	}
	return author, nil
}

func DeleteAuthor(id string) error {
	queries, ctx := db.DbQueries()
	err := queries.DeleteAuthor(ctx, id)
	if err != nil {
		slog.Error("저자 삭제 오류", "error", err.Error())
		return err
	}
	return nil
}
