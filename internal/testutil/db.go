package testutil

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// NewTestDB creates an in-memory SQLite database with predefined test tables for unit testing purposes.
// It initializes the database connection, enforces foreign key constraints, and sets up test table schemas.
// Accepts a *testing.T instance for error handling and terminates the test upon any setup failure.
// Returns a pointer to the initialized *sql.DB instance.
func NewTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		t.Fatal(err)
	}

	if err := setupTestTables(db); err != nil {
		t.Fatal(err)
	}

	return db
}

// setupTestTables initializes tables in the given database with predefined schemas for testing purposes.
// It accepts a database connection and executes relevant SQL create table statements.
// Returns an error if any SQL execution fails.
func setupTestTables(db *sql.DB) error {
	// Using the same schema as in your migrations
	statements := []string{
		`CREATE TABLE users (
            id         INTEGER PRIMARY KEY AUTOINCREMENT,
            username   TEXT NOT NULL UNIQUE,
            email      TEXT NOT NULL UNIQUE,
            password   TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )`,
		`CREATE TABLE books (
            id               INTEGER PRIMARY KEY AUTOINCREMENT,
            title            TEXT NOT NULL,
            author           TEXT NOT NULL,
            isbn             TEXT UNIQUE,
            publication_year INTEGER,
            image_url        TEXT DEFAULT NULL,
            created_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at       DATETIME DEFAULT CURRENT_TIMESTAMP
        )`,
		`CREATE TABLE user_books (
            id       INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id  INTEGER NOT NULL,
            book_id  INTEGER NOT NULL,
            status   TEXT CHECK (status IN ('want_to_read', 'reading', 'finished')) DEFAULT 'want_to_read',
            added_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
            FOREIGN KEY (book_id) REFERENCES books (id) ON DELETE CASCADE,
            UNIQUE (user_id, book_id)
        )`,
	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}

	return nil
}
