-- name: CreateUser :exec
INSERT INTO users (id, name, email, hashed_password, created)
VALUES(?, ?, ?, ?, DATETIME('now', 'utc'));

-- name: AuthenticateUser :one
SELECT id, hashed_password 
FROM users
WHERE email = ?;

-- name: UserExists :one
SELECT EXISTS(
  SELECT true
  FROM users
  WHERE id = ?
);
