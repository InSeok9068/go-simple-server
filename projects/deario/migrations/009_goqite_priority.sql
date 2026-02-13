-- +goose Up
ALTER TABLE goqite ADD COLUMN priority INTEGER NOT NULL DEFAULT 0;
CREATE INDEX IF NOT EXISTS goqite_queue_priority_created_idx ON goqite (queue, priority DESC, created);

-- +goose Down
DROP INDEX IF EXISTS goqite_queue_priority_created_idx;
