// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package db

import (
	"database/sql"
)

type Diary struct {
	ID         string
	Uid        string
	Date       string
	Content    string
	AiFeedback string
	AiImage    string
	Created    sql.NullString
	Updated    sql.NullString
}

type PushKey struct {
	ID      string
	Uid     string
	Token   string
	Created sql.NullString
	Updated sql.NullString
}

type User struct {
	Uid     string
	Name    string
	Email   string
	Created sql.NullString
	Updated sql.NullString
}
