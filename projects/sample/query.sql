-- name: GetAuthor :one 
SELECT    *
FROM      authors
WHERE     id = ?
LIMIT     1;

-- name: ListAuthors :many 
SELECT    *
FROM      authors
ORDER BY  name;

-- name: CreateAuthor :one
INSERT    INTO authors (name, bio, created, updated)
VALUES    (?, ?, datetime ('now'), datetime ('now')) RETURNING *;

-- name: UpdateAuthor :one 
UPDATE    authors
set       name = ?,
          bio = ?,
          updated = datetime ('now')
WHERE     id = ? RETURNING *;

-- name: DeleteAuthor :exec 
DELETE    FROM authors
WHERE     id = ?;