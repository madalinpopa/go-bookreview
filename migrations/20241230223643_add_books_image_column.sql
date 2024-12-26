-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- Add image column to books table
ALTER TABLE books ADD COLUMN image_url TEXT DEFAULT NULL;

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

-- Remove image column from books table
ALTER TABLE books DROP COLUMN image_url;