-- +goose Up
CREATE UNIQUE INDEX IF NOT EXISTS idx_diary_uid_date
ON diary (uid, date);

-- +goose Down
DROP INDEX IF EXISTS idx_diary_uid_date;
