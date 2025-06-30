-- +goose Up
CREATE TABLE IF NOT EXISTS user (
    uid TEXT PRIMARY KEY,
    name TEXT DEFAULT '' NOT NULL,
    email TEXT UNIQUE DEFAULT '' NOT NULL,
    created TEXT DEFAULT CURRENT_TIMESTAMP,
    updated TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS diary (
    id TEXT DEFAULT (
        'r' || LOWER(HEX(RANDOMBLOB(7)))
    ) NOT NULL PRIMARY KEY,
    uid TEXT DEFAULT '' NOT NULL,
    date TEXT DEFAULT '' NOT NULL,
    content TEXT DEFAULT '' NOT NULL,
    ai_feedback TEXT DEFAULT '' NOT NULL,
    ai_image TEXT DEFAULT '' NOT NULL,
    created TEXT DEFAULT CURRENT_TIMESTAMP,
    updated TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS push_key (
    id TEXT DEFAULT (
        'r' || LOWER(HEX(RANDOMBLOB(7)))
    ) NOT NULL PRIMARY KEY,
    uid TEXT DEFAULT '' NOT NULL,
    token TEXT DEFAULT '' NOT NULL,
    created TEXT DEFAULT CURRENT_TIMESTAMP,
    updated TEXT DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO
    user (
        uid,
        name,
        email,
        created,
        updated
    )
SELECT uid, name, email, created, updated
FROM users;

INSERT INTO
    diary (
        id,
        uid,
        date,
        content,
        ai_feedback,
        ai_image,
        created,
        updated
    )
SELECT
    id,
    uid,
    date,
    content,
    aiFeedback,
    aiImage,
    created,
    updated
FROM diarys;

INSERT INTO
    push_key (
        id,
        uid,
        token,
        created,
        updated
    )
SELECT id, uid, token, created, updated
FROM push_keys;

DROP TABLE users;

DROP TABLE diarys;

DROP TABLE push_keys;

-- +goose Down
DROP TABLE user;

DROP TABLE diary;

DROP TABLE push_key;