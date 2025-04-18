-- sqlc generate -f ./projects/deario/sqlc.yaml

-- name: GetDiary :one
SELECT *
FROM diarys
WHERE date = ?
  AND uid = ?
LIMIT 1;

-- name: GetDiaryRandom :one
SELECT *
FROM diarys
WHERE date IS NOT NULL
  AND uid = ?
ORDER BY RANDOM()
LIMIT 1;

-- name: ListDiarys :many
SELECT *
FROM diarys
WHERE uid = ?
ORDER BY created desc
LIMIT 10 OFFSET ((? - 1) * 10);

-- name: CreateDiary :one
INSERT INTO diarys (uid, content, date, created, updated)
VALUES (?,
        ?,
        ?,
--         strftime('%Y%m%d', 'now', 'localtime'),
        datetime('now', 'localtime'),
        datetime('now', 'localtime'))
RETURNING *;

-- name: UpdateDiary :one
UPDATE diarys
set content = ?,
    updated = datetime('now')
WHERE id = ?
RETURNING *;