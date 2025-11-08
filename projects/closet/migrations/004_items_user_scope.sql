-- +goose Up
PRAGMA foreign_keys = OFF;

ALTER TABLE items RENAME TO items_old;

CREATE TABLE items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_uid TEXT NOT NULL REFERENCES user (uid) ON DELETE CASCADE,
    kind TEXT NOT NULL CHECK (
        kind IN (
            'top',
            'bottom',
            'shoes',
            'accessory'
        )
    ),
    filename TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    bytes BLOB NOT NULL,
    thumb_bytes BLOB,
    sha256 TEXT,
    width INTEGER,
    height INTEGER,
    created_at INTEGER NOT NULL DEFAULT(strftime('%s', 'now')),
    meta_summary TEXT DEFAULT '',
    meta_season TEXT DEFAULT '',
    meta_style TEXT DEFAULT '',
    meta_colors TEXT DEFAULT ''
);

INSERT OR IGNORE INTO user (uid, name, email)
VALUES ('__legacy__', 'Legacy Closet User', 'legacy@example.com');

INSERT INTO items (
    id,
    user_uid,
    kind,
    filename,
    mime_type,
    bytes,
    thumb_bytes,
    sha256,
    width,
    height,
    created_at,
    meta_summary,
    meta_season,
    meta_style,
    meta_colors
)
SELECT
    id,
    '__legacy__' AS user_uid,
    kind,
    filename,
    mime_type,
    bytes,
    thumb_bytes,
    sha256,
    width,
    height,
    created_at,
    meta_summary,
    meta_season,
    meta_style,
    meta_colors
FROM items_old;

DROP TABLE items_old;

PRAGMA foreign_keys = ON;

CREATE INDEX IF NOT EXISTS idx_items_user_kind_created ON items (user_uid, kind, created_at DESC);
CREATE UNIQUE INDEX IF NOT EXISTS idx_items_user_sha ON items (user_uid, sha256);
CREATE INDEX IF NOT EXISTS idx_items_sha256 ON items (sha256);

-- +goose Down
PRAGMA foreign_keys = OFF;

ALTER TABLE items RENAME TO items_new;

CREATE TABLE items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    kind TEXT NOT NULL CHECK (
        kind IN (
            'top',
            'bottom',
            'shoes',
            'accessory'
        )
    ),
    filename TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    bytes BLOB NOT NULL,
    thumb_bytes BLOB,
    sha256 TEXT UNIQUE,
    width INTEGER,
    height INTEGER,
    created_at INTEGER NOT NULL DEFAULT(strftime('%s', 'now')),
    meta_summary TEXT DEFAULT '',
    meta_season TEXT DEFAULT '',
    meta_style TEXT DEFAULT '',
    meta_colors TEXT DEFAULT ''
);

INSERT INTO items (
    id,
    kind,
    filename,
    mime_type,
    bytes,
    thumb_bytes,
    sha256,
    width,
    height,
    created_at,
    meta_summary,
    meta_season,
    meta_style,
    meta_colors
)
SELECT
    id,
    kind,
    filename,
    mime_type,
    bytes,
    thumb_bytes,
    sha256,
    width,
    height,
    created_at,
    meta_summary,
    meta_season,
    meta_style,
    meta_colors
FROM items_new;

DROP TABLE items_new;

PRAGMA foreign_keys = ON;

CREATE INDEX IF NOT EXISTS idx_items_kind_created ON items (kind, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_items_sha256 ON items (sha256);
