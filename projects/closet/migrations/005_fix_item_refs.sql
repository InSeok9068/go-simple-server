-- +goose Up
PRAGMA foreign_keys = OFF;

-- item_tags 재생성 (items FK 재연결)
CREATE TABLE item_tags_new (
    item_id INTEGER NOT NULL REFERENCES items (id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags (id) ON DELETE CASCADE,
    PRIMARY KEY (item_id, tag_id)
);

INSERT INTO item_tags_new (item_id, tag_id)
SELECT item_id, tag_id
FROM item_tags;

DROP TABLE item_tags;
ALTER TABLE item_tags_new RENAME TO item_tags;
CREATE INDEX IF NOT EXISTS idx_item_tags_tag ON item_tags (tag_id);

-- embeddings 재생성 (items FK 재연결)
CREATE TABLE embeddings_new (
    item_id INTEGER PRIMARY KEY REFERENCES items (id) ON DELETE CASCADE,
    model TEXT NOT NULL,
    dim INTEGER NOT NULL,
    vec_f32 BLOB NOT NULL
);

INSERT INTO embeddings_new (item_id, model, dim, vec_f32)
SELECT item_id, model, dim, vec_f32
FROM embeddings;

DROP TABLE embeddings;
ALTER TABLE embeddings_new RENAME TO embeddings;

PRAGMA foreign_keys = ON;

-- +goose Down
-- FK 재구성 전 상태로 되돌릴 수 있는 안전한 방법이 없어 더 이상 작업하지 않는다.
SELECT 1;
