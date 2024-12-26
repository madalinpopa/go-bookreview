-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- Add status column to user_books table with a check constraint and default value
ALTER TABLE user_books
    ADD COLUMN status TEXT
        CHECK (status IN ('want_to_read', 'reading', 'finished'))
        DEFAULT 'want_to_read';

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

-- Remove status column from user_books table
ALTER TABLE user_books
    DROP COLUMN status;