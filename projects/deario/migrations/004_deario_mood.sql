-- +goose Up
ALTER TABLE diary
ADD COLUMN mood TEXT DEFAULT '0' NOT NULL CHECK (
    mood IN ('0', '1', '2', '3', '4', '5')
);

-- +goose Down
ALTER TABLE diary DROP COLUMN mood;