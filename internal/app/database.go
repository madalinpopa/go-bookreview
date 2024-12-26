package app

import (
	"database/sql"
	"fmt"
	"github.com/madalinpopa/go-bookreview/migrations"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

// sqliteDriver specifies the driver name for connecting to a SQLite database using the database/sql package.
const sqliteDriver = "sqlite3"

// MakeMigrations applies database migrations using the provided database connection.
// It sets the migration file system and database dialect before running the migrations.
// Returns an error if setting the dialect or applying the migrations fails.
func MakeMigrations(DB *sql.DB) error {

	goose.SetBaseFS(migrations.MigrationFiles)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	if err := goose.Up(DB, "."); err != nil {
		return err
	}
	return nil
}

// CreateDatabaseConnection opens a SQLite database using the provided DSN and verifies the connection.
// Returns a database connection or an error if opening or pinging the database fails.
func CreateDatabaseConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open(sqliteDriver, fmt.Sprintf("%s?_foreign_keys=1", dsn))
	if err != nil {
		return nil, handleDatabaseError(db, err)
	}

	if err = db.Ping(); err != nil {
		return nil, handleDatabaseError(db, err)
	}

	return db, nil
}

// handleDatabaseError ensures the database is closed and returns the error.
func handleDatabaseError(db *sql.DB, err error) error {
	if db != nil {
		_ = db.Close()
	}
	return err
}
