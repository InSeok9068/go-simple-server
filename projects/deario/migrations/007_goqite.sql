-- +goose Up
CREATE TABLE goqite (
    id TEXT PRIMARY KEY DEFAULT (
        'm_' || LOWER(HEX(RANDOMBLOB(16)))
    ),
    created TEXT NOT NULL DEFAULT (
        strftime('%Y-%m-%dT%H:%M:%fZ')
    ),
    updated TEXT NOT NULL DEFAULT (
        strftime('%Y-%m-%dT%H:%M:%fZ')
    ),
    queue TEXT NOT NULL,
    body BLOB NOT NULL,
    timeout TEXT NOT NULL DEFAULT (
        strftime('%Y-%m-%dT%H:%M:%fZ')
    ),
    received INTEGER NOT NULL DEFAULT 0
) STRICT;

-- +goose StatementBegin
CREATE TRIGGER goqite_updated_timestamp
AFTER UPDATE ON goqite
BEGIN
  UPDATE goqite
  SET    updated = strftime('%Y-%m-%dT%H:%M:%fZ')
  WHERE  id = OLD.id;
END;
-- +goose StatementEnd

CREATE INDEX goqite_queue_created_idx ON goqite (queue, created);

-- +goose Down
DROP TABLE goqite;