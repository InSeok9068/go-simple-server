-- +goose Up
CREATE TABLE IF NOT EXISTS user_setting (
    uid TEXT PRIMARY KEY,
    is_push INTEGER DEFAULT 0 NOT NULL CHECK (is_push IN (0, 1)),
    push_token TEXT DEFAULT '' NOT NULL,
    push_time TEXT DEFAULT '' NOT NULL CHECK (
        push_time = ''
        OR (
            push_time GLOB '[0-2][0-9]:[0-5][0-9]'
            AND substr(push_time, 1, 2) BETWEEN '00' AND '23'
        )
    ),
    random_range INTEGER DEFAULT 365 NOT NULL,
    created TEXT DEFAULT CURRENT_TIMESTAMP,
    updated TEXT DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO
    user_setting (uid, push_token)
SELECT uid, token
FROM push_key;

DROP TABLE push_key;

-- +goose Down
DROP TABLE user_setting;