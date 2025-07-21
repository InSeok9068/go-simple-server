-- +goose Up
ALTER TABLE user_setting ADD COLUMN pin_enabled INTEGER DEFAULT 0 NOT NULL CHECK (pin_enabled IN (0,1));
ALTER TABLE user_setting ADD COLUMN pin TEXT DEFAULT '' NOT NULL;
ALTER TABLE user_setting ADD COLUMN pin_cycle INTEGER DEFAULT -1 NOT NULL;
ALTER TABLE user_setting ADD COLUMN pin_last_at TEXT DEFAULT '' NOT NULL;

-- +goose Down
ALTER TABLE user_setting DROP COLUMN pin_enabled;
ALTER TABLE user_setting DROP COLUMN pin;
ALTER TABLE user_setting DROP COLUMN pin_cycle;
ALTER TABLE user_setting DROP COLUMN pin_last_at;
