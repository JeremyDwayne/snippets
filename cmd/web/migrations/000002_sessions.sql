-- +goose Up
CREATE TABLE IF NOT EXISTS sessions (
  token TEXT PRIMARY KEY,
  data BLOB NOT NULL,
  expiry DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions (expiry);

-- +goose Down
DROP TABLE IF EXISTS sessions;
