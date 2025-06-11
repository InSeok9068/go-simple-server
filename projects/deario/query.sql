-- sqlc generate -f ./projects/deario/sqlc.yaml

-- name: GetDiary :one
SELECT * FROM diarys WHERE date = ? AND uid = ? LIMIT 1;

-- name: GetDiaryRandom :one
SELECT *
FROM diarys
WHERE
    date IS NOT NULL
    AND uid = ?
ORDER BY RANDOM()
LIMIT 1;

-- name: ListDiarys :many
SELECT *
FROM diarys
WHERE
    uid = ?
ORDER BY created desc
LIMIT 7
OFFSET ((? - 1) * 7);

-- name: CreateDiary :one
INSERT INTO
    diarys (
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
UPDATE diarys
SET
    content = ?,
    updated = datetime('now')
WHERE
    id = ? RETURNING *;

-- name: DeleteDiary :exec
DELETE FROM diarys WHERE id = ?;

-- name: UpdateDiaryOfAiFeedback :exec
UPDATE diarys
SET
    aiFeedback = ?,
    aiImage = ?,
    updated = datetime('now')
WHERE
    id = ?;

-- name: GetPushKey :one
SELECT * FROM push_keys WHERE uid = ? LIMIT 1;

-- name: CreatePushKey :exec
INSERT INTO
    push_keys (uid, token, created, updated)
VALUES (
        ?,
        ?,
        datetime('now', 'localtime'),
        datetime('now', 'localtime')
    );

-- name: UpdatePushKey :exec
UPDATE push_keys
SET
    token = ?,
    updated = datetime('now')
WHERE
    uid = ?;