package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

const (
	ID LookupField = iota
	Email
	Username
)

// minPasswordCost specifies the minimum computational cost used for hashing passwords to ensure adequate security.
const minPasswordCost = 12

// User represents an application user including identification, authentication, and contact information.
type User struct {
	Base
	Email          string
	Username       string
	HashedPassword []byte
}

// UserModel provides methods to interact with the users table in the database. It wraps a SQL database connection.
type UserModel struct {
	DB     *sql.DB
	Logger *slog.Logger
}

// Create adds a new user record to the database with a hashed password. It returns an error if insertion fails or data is invalid.
func (m *UserModel) Create(username, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), minPasswordCost)
	if err != nil {
		return err
	}

	stmt := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"

	_, err = m.DB.Exec(stmt, username, email, hashedPassword)
	if err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) && errors.Is(sqliteError.ExtendedCode, sqlite3.ErrConstraintUnique) {
			if sqliteError.Error() == "UNIQUE constraint failed: users.username" {
				return ErrDuplicateUsername
			} else if sqliteError.Error() == "UNIQUE constraint failed: users.email" {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

// Exists checks if a record exists in the users table based on the specified lookup field and value. Returns a boolean and error.
func (m *UserModel) Exists(field LookupField, value any) (bool, error) {
	var exists bool
	var stmt string

	switch field {
	case ID:
		stmt = "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"
	case Email:
		stmt = "SELECT EXISTS(SELECT true FROM users WHERE email = ?)"
	case Username:
		stmt = "SELECT EXISTS(SELECT true FROM users WHERE username = ?)"
	default:
		return false, fmt.Errorf("invalid lookup field")
	}

	err := m.DB.QueryRow(stmt, value).Scan(&exists)
	return exists, err
}

// Authenticate verifies a user's identity by comparing the provided password with the hashed password in the database.
// It returns the user's ID if authentication is successful or an error
// if the credentials are invalid or other issues occur.
func (m *UserModel) Authenticate(username, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, password FROM users WHERE username = ?"
	err := m.DB.QueryRow(stmt, username).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}
