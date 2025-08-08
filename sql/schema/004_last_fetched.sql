-- +goose Up
Alter TABLE feeds
add column last_fetched_at TIMESTAMP WITH TIME ZONE;

-- +goose Down
ALTER TABLE feeds
DROP COLUMN last_fetched_at;