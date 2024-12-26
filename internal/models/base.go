package models

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"
)

var (

	// ErrNoRecord indicates that no matching record was found in the database or data source.
	ErrNoRecord = errors.New("models: no matching record found")

	// ErrInvalidCredentials indicates that the provided credentials are invalid or do not match any existing records.
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	// ErrDuplicateEmail indicates that the provided email address already exists in the system and cannot be used again.
	ErrDuplicateEmail = errors.New("models: duplicate email")

	// ErrDuplicateUsername indicates that the provided username already exists in the system and cannot be used again.
	ErrDuplicateUsername = errors.New("models: duplicate username")

	// ErrDuplicateIsbn indicates that the provided ISBN already exists in the system and cannot be used again.
	ErrDuplicateIsbn = errors.New("models: duplicate isbn")
)

// LookupField represents an enumeration used to specify fields for lookup operations in user-related database queries.
type LookupField int

// Base represents a common structure containing fields for ID and timestamps to be embedded in other types.
type Base struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Models aggregates different data models to provide centralized access within the application.
type Models struct {
	Users   UserModel
	Books   BookModel
	Notes   NoteModel
	Reviews ReviewModel
}

// NewModels initializes and returns a Models instance with the provided database connection.
func NewModels(db *sql.DB, logger *slog.Logger) *Models {
	return &Models{
		Users:   UserModel{DB: db, Logger: logger},
		Books:   BookModel{DB: db, Logger: logger},
		Notes:   NoteModel{DB: db, Logger: logger},
		Reviews: ReviewModel{DB: db, Logger: logger},
	}
}
