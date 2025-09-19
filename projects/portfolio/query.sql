-- sqlc generate -f ./projects/portfolio/sqlc.yaml

-- name: GetUser :one
SELECT * FROM user WHERE uid = ? LIMIT 1;

-- name: CreateUser :exec
INSERT INTO user (uid, name, email) VALUES (?, ?, ?);