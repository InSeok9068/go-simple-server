-- auto-generated definition
CREATE TABLE users (
    uid TEXT PRIMARY KEY,
    name TEXT DEFAULT '' NOT NULL,
    email TEXT UNIQUE DEFAULT '' NOT NULL,
    created TEXT DEFAULT CURRENT_TIMESTAMP,
    updated TEXT DEFAULT CURRENT_TIMESTAMP
);

-- auto-generated definition
CREATE TABLE diarys (
    id TEXT DEFAULT(
        'r' || LOWER(HEX(RANDOMBLOB (7)))
    ) NOT NULL PRIMARY KEY,
    uid TEXT DEFAULT '' NOT NULL,
    date TEXT DEFAULT '' NOT NULL,
    content TEXT DEFAULT '' NOT NULL,
    aiFeedback TEXT DEFAULT '' NOT NULL,
    aiImage TEXT DEFAULT '' NOT NULL,
    created TEXT DEFAULT CURRENT_TIMESTAMP,
    updated TEXT DEFAULT CURRENT_TIMESTAMP
);

-- auto-generated definition
CREATE TABLE push_keys (
    id TEXT DEFAULT(
        'r' || LOWER(HEX(RANDOMBLOB (7)))
    ) NOT NULL PRIMARY KEY,
    uid TEXT DEFAULT '' NOT NULL,
    token TEXT DEFAULT '' NOT NULL,
    created TEXT DEFAULT CURRENT_TIMESTAMP,
    updated TEXT DEFAULT CURRENT_TIMESTAMP
);