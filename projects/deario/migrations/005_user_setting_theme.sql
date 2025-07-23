-- +goose Up
ALTER TABLE user_setting
ADD COLUMN theme TEXT DEFAULT 'light' NOT NULL CHECK (theme IN ('light', 'dark'));

-- +goose Down
ALTER TABLE user_setting DROP COLUMN theme;
