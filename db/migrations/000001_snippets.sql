-- +goose Up
CREATE TABLE snippets (
  id TEXT NOT NULL PRIMARY KEY,
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  expires DATETIME NOT NULL,
  created DATETIME NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets (created);

-- +goose Down
DROP TABLE IF EXISTS snippets;
