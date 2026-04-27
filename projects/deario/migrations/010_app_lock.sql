-- +goose Up
ALTER TABLE user_setting ADD COLUMN app_lock_enabled INTEGER DEFAULT 0 NOT NULL CHECK (app_lock_enabled IN (0, 1));

ALTER TABLE user_setting ADD COLUMN app_lock_pin_hash TEXT DEFAULT '' NOT NULL;

-- +goose Down
ALTER TABLE user_setting DROP COLUMN app_lock_enabled;

ALTER TABLE user_setting DROP COLUMN app_lock_pin_hash;
