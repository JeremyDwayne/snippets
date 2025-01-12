-- +goose Up
CREATE TABLE IF NOT EXISTS users (
  id TEXT NOT NULL PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  hashed_password TEXT NOT NULL,
  created DATETIME NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS users;
