-- sqlc generate -f ./projects/closet/sqlc.yaml

-- name: GetUser :one
SELECT * FROM user WHERE uid = ?;