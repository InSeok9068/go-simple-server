-- sqlc generate -f ./projects/deario/sqlc.yaml

-- name: GetDiary :one
SELECT * FROM diary WHERE date = ? AND uid = ? LIMIT 1;

-- name: GetDiaryRandom :one
SELECT *
FROM diary
WHERE
    date IS NOT NULL
    AND uid = ?
    AND date >= ?
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
    diary (uid, content, date)
VALUES (?, ?, ?) RETURNING *;

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

-- name: UpdateDiaryOfMood :exec
UPDATE diary
SET
    mood = ?,
    updated = datetime('now')
WHERE
    id = ?;

-- name: GetUserSetting :one
SELECT * FROM user_setting WHERE uid = ? LIMIT 1;

-- name: UpsertUserSetting :exec
INSERT INTO
    user_setting (
        uid,
        is_push,
        push_time,
        random_range
    )
VALUES (?, ?, ?, ?)
ON CONFLICT (uid) DO
UPDATE
SET
    is_push = excluded.is_push,
    push_time = excluded.push_time,
    random_range = excluded.random_range,
    updated = datetime('now');

-- name: UpsertPushKey :exec
INSERT INTO
    user_setting (uid, push_token)
VALUES (?, ?)
ON CONFLICT (uid) DO
UPDATE
SET
    push_token = excluded.push_token,
    updated = datetime('now');

-- name: ListPushTargets :many
SELECT
    uid,
    push_token,
    push_time,
    random_range
FROM user_setting
WHERE
    is_push = 1
    AND push_token != ''
    AND push_time != '';

-- name: GetUser :one
SELECT * FROM user WHERE uid = ? LIMIT 1;

-- name: CreateUser :exec
INSERT INTO user (uid, name, email) VALUES (?, ?, ?);