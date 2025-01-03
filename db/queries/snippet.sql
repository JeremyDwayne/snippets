-- name: CreateSnippet :exec
INSERT INTO snippets (id, title, content, created, expires)
VALUES(?, ?, ?, DATETIME('now', 'utc'), ?);

-- name: GetSnippet :one
SELECT *
FROM snippets
WHERE expires > DATETIME('now', 'utc') AND id = ?;

-- name: LatestSnippets :many
SELECT *
FROM snippets
WHERE expires > DATETIME('now', 'utc') ORDER BY created DESC LIMIT 10;
