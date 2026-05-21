-- +goose Up
ALTER TABLE attendances ADD COLUMN location_status VARCHAR(20) DEFAULT 'IN_AREA';

-- +goose Down
ALTER TABLE attendances DROP COLUMN location_status;
