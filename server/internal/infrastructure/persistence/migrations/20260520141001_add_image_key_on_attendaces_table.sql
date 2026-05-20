-- +goose Up
ALTER TABLE attendances ADD COLUMN photo_key VARCHAR(500);

-- +goose Down
ALTER TABLE attendances DROP COLUMN IF EXISTS photo_key;
