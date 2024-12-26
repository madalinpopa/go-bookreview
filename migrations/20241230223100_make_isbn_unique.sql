-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- Create a new table with the UNIQUE constraint
CREATE TABLE books_new
(
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    title            TEXT NOT NULL,
    author           TEXT NOT NULL,
    isbn             TEXT UNIQUE,
    publication_year INTEGER,
    created_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at       DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Copy data from the old table to the new table
INSERT INTO books_new (id, title, author, isbn, publication_year, created_at, updated_at)
SELECT id, title, author, isbn, publication_year, created_at, updated_at
FROM books;

-- Drop the old table
DROP TABLE books;

-- Rename the new table to the original name
ALTER TABLE books_new
    RENAME TO books;

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

-- Create a new table without the UNIQUE constraint
CREATE TABLE books_new
(
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    title            TEXT NOT NULL,
    author           TEXT NOT NULL,
    isbn             TEXT,
    publication_year INTEGER,
    created_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at       DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Copy data from the old table to the new table
INSERT INTO books_new (id, title, author, isbn, publication_year, created_at, updated_at)
SELECT id, title, author, isbn, publication_year, created_at, updated_at
FROM books;

-- Drop the old table
DROP TABLE books;

-- Rename the new table to the original name
ALTER TABLE books_new
    RENAME TO books;