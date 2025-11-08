-- +goose Up
ALTER TABLE items ADD COLUMN meta_summary TEXT DEFAULT '';
ALTER TABLE items ADD COLUMN meta_season TEXT DEFAULT '';
ALTER TABLE items ADD COLUMN meta_style TEXT DEFAULT '';
ALTER TABLE items ADD COLUMN meta_colors TEXT DEFAULT '';

-- +goose Down
-- SQLite는 컬럼 DROP을 지원하지 않으므로 다운 마이그레이션에서는 변경 사항을 되돌리지 않습니다.
SELECT 1;
