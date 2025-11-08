-- sqlc generate -f ./projects/closet/sqlc.yaml

-- name: GetUser :one
SELECT * FROM user WHERE uid = ?;

-- name: InsertItem :one
INSERT INTO items (
    user_uid,
    kind,
    filename,
    mime_type,
    bytes,
    thumb_bytes,
    sha256,
    width,
    height,
    meta_summary,
    meta_season,
    meta_style,
    meta_colors
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id;

-- name: GetItemIDBySha :one
SELECT id
FROM items
WHERE user_uid = sqlc.arg(user_uid)
  AND sha256 = sqlc.arg(sha256);

-- name: UpsertTag :one
INSERT INTO tags (name) VALUES (?)
ON CONFLICT(name) DO UPDATE SET name = excluded.name
RETURNING id;

-- name: AttachTag :exec
INSERT OR IGNORE INTO item_tags (item_id, tag_id) VALUES (?, ?);

-- name: ListItems :many
SELECT
    i.id,
    i.kind,
    i.filename,
    i.mime_type,
    i.width,
    i.height,
    i.created_at,
    IFNULL(GROUP_CONCAT(t.name, ','), '') AS tags
FROM items i
LEFT JOIN item_tags it ON it.item_id = i.id
LEFT JOIN tags t ON t.id = it.tag_id
WHERE i.user_uid = sqlc.arg(user_uid)
  AND (sqlc.arg(kind_filter) = '' OR i.kind = sqlc.arg(kind_filter))
GROUP BY i.id
HAVING (
    CASE
        WHEN json_array_length(sqlc.arg(tag_json)) = 0 THEN 1
        ELSE IFNULL((
            SELECT COUNT(DISTINCT t2.name)
            FROM item_tags it2
            JOIN tags t2 ON t2.id = it2.tag_id
            WHERE it2.item_id = i.id
              AND t2.name IN (
                  SELECT je.value
                  FROM (SELECT sqlc.arg(tag_json) AS json_data) payload
                  CROSS JOIN json_each(payload.json_data) AS je
              )
        ), 0)
    END
) >= CASE
        WHEN json_array_length(sqlc.arg(tag_json)) = 0 THEN 1
        ELSE json_array_length(sqlc.arg(tag_json))
    END
ORDER BY i.created_at DESC
LIMIT sqlc.arg(limit)
OFFSET sqlc.arg(offset);

-- name: GetItemContent :one
SELECT bytes, mime_type
FROM items
WHERE id = sqlc.arg(id)
  AND user_uid = sqlc.arg(user_uid);

-- name: PutEmbedding :exec
INSERT INTO embeddings (item_id, model, dim, vec_f32)
VALUES (?, ?, ?, ?)
ON CONFLICT(item_id) DO UPDATE
SET model = excluded.model,
    dim = excluded.dim,
    vec_f32 = excluded.vec_f32;

-- name: LoadEmbeddingsByIDs :many
SELECT e.item_id, e.dim, e.vec_f32
FROM embeddings e
WHERE e.item_id IN (sqlc.slice('ids'));

-- name: ListEmbeddingItems :many
SELECT
    i.id,
    i.kind,
    i.meta_season,
    i.meta_style,
    i.meta_colors,
    e.dim,
    e.vec_f32
FROM embeddings e
JOIN items i ON i.id = e.item_id
WHERE i.user_uid = sqlc.arg(user_uid);

-- name: ListItemsByIDs :many
SELECT
    i.id,
    i.kind,
    i.filename,
    i.mime_type,
    i.width,
    i.height,
    i.created_at,
    IFNULL(GROUP_CONCAT(t.name, ','), '') AS tags
FROM items i
LEFT JOIN item_tags it ON it.item_id = i.id
LEFT JOIN tags t ON t.id = it.tag_id
WHERE i.user_uid = sqlc.arg(user_uid)
  AND i.id IN (sqlc.slice('ids'))
GROUP BY i.id;

-- name: DeleteItem :exec
DELETE FROM items
WHERE id = sqlc.arg(id)
  AND user_uid = sqlc.arg(user_uid);
