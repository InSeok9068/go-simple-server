-- sqlc generate -f ./projects/deario/sqlc.yaml

-- name: GetDiary :one
SELECT * FROM diary WHERE date = ? AND uid = ? LIMIT 1;

-- name: GetDiaryRandom :one
SELECT *
FROM diary
WHERE
    date IS NOT NULL
    AND uid = ?
ORDER BY RANDOM()
LIMIT 1;

-- name: ListDiarys :many
SELECT *
FROM diary
WHERE
    uid = ?
ORDER BY created desc
LIMIT 7
OFFSET ((? - 1) * 7);

-- name: CreateDiary :one
INSERT INTO
    diary (
        uid,
        content,
        date,
        created,
        updated
    )
VALUES (
        ?,
        ?,
        ?,
        --         strftime('%Y%m%d', 'now', 'localtime'),
        datetime('now', 'localtime'),
        datetime('now', 'localtime')
    ) RETURNING *;

-- name: UpdateDiary :one
UPDATE diary
SET
    content = ?,
    updated = datetime('now')
WHERE
    id = ? RETURNING *;

-- name: DeleteDiary :exec
DELETE FROM diary WHERE id = ?;

-- name: UpdateDiaryOfAiFeedback :exec
UPDATE diary
SET
    ai_feedback = ?,
    ai_image = ?,
    updated = datetime('now')
WHERE
    id = ?;

-- name: GetPushKey :one
SELECT * FROM push_key WHERE uid = ? LIMIT 1;

-- name: CreatePushKey :exec
INSERT INTO
    push_key (uid, token, created, updated)
VALUES (
        ?,
        ?,
        datetime('now', 'localtime'),
        datetime('now', 'localtime')
    );

-- name: UpdatePushKey :exec
UPDATE push_key
SET
    token = ?,
    updated = datetime('now')
WHERE
    uid = ?;

-- name: GetUser :one
SELECT * FROM user WHERE uid = ? LIMIT 1;

-- name: CreateUser :exec
INSERT INTO
    user (
        uid,
        name,
        email,
        created,
        updated
    )
VALUES (
        ?,
        ?,
        ?,
        datetime('now', 'localtime'),
        datetime('now', 'localtime')
    );