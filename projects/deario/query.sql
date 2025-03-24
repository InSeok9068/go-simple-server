-- name: GetDiary :one
SELECT *
FROM diarys
WHERE date = ?
  AND uid = ?
LIMIT 1;

-- name: CreateDiary :one
INSERT INTO diarys (uid, content, date, created, updated)
VALUES (?,
        ?,
        strftime('%Y%m%d', 'now', 'localtime'),
        datetime('now', 'localtime'),
        datetime('now', 'localtime'))
RETURNING *;

-- name: UpdateDiary :one
UPDATE diarys
set content = ?,
    updated = datetime('now')
WHERE id = ?
RETURNING *;