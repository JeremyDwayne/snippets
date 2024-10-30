-- name: UserExists :one
SELECT EXISTS(SELECT true FROM users WHERE id = ?);
