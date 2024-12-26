package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"log/slog"
	"time"
)

// Note represents a user's note on a specific page of a book with associated user and book identifiers.
type Note struct {
	ID         int
	UserId     int
	BookId     int
	NoteText   string
	PageNumber int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NoteModel wraps a database connection pool for managing operations related to notes.
type NoteModel struct {
	DB     *sql.DB
	Logger *slog.Logger
}

// Create inserts a new note into the database and returns its ID or an error if the operation fails.
func (n *NoteModel) Create(userId, bookId int, noteText string, pageNumber int) (int, error) {

	stmt := `INSERT INTO notes (user_id, book_id, note_text, page_number) VALUES (?, ?, ?, ?)`

	result, err := n.DB.Exec(stmt, userId, bookId, noteText, pageNumber)
	if err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) {
			return 0, sqliteError
		} else {
			return 0, err
		}
	}

	noteId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(noteId), nil
}

// Retrieve fetches a note from the database for a given user ID and book ID, returning the note or an error.
func (n *NoteModel) Retrieve(userId, noteId int) (Note, error) {
	var note Note
	stmt := `SELECT id, user_id, book_id, note_text, page_number, created_at, updated_at FROM notes WHERE user_id = ? AND id = ?`

	err := n.DB.QueryRow(stmt, userId, noteId).Scan(
		&note.ID,
		&note.UserId,
		&note.BookId,
		&note.NoteText,
		&note.PageNumber,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Note{}, ErrNoRecord
		} else {
			return Note{}, err
		}
	}
	return note, nil
}

// Update modifies an existing note's text and page number using the provided user ID and book ID.
// Returns an error if the update fails or no record is found.
func (n *NoteModel) Update(userId, noteId int, noteText string, pageNumber int) error {
	stmt := `UPDATE notes SET note_text = ?, page_number = ? WHERE user_id = ? AND id = ?`

	result, err := n.DB.Exec(stmt, noteText, pageNumber, userId, noteId)
	if err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) {
			return sqliteError
		} else {
			return err
		}
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNoRecord
	}

	return nil
}

// Delete removes a note from the database based on the given user ID and book ID. Returns an error if the operation fails.
func (n *NoteModel) Delete(userId, noteId int) error {
	stmt := `DELETE FROM notes WHERE user_id = ? AND id = ?`

	result, err := n.DB.Exec(stmt, userId, noteId)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("note with user_id %d and id %d not found", userId, noteId)
	}

	return nil
}

// List retrieves all notes associated with a specific book ID from the database and returns them or an error if it fails.
func (n *NoteModel) List(bookId, userId int) ([]Note, error) {
	stmt := `SELECT id, user_id, book_id, note_text, page_number, created_at, updated_at 
			FROM notes WHERE book_id = ? AND user_id = ? 
			`

	rows, err := n.DB.Query(stmt, bookId, userId)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			n.Logger.Error(err.Error())
		}
	}()

	var notes []Note
	for rows.Next() {
		var note Note
		err = rows.Scan(
			&note.ID,
			&note.UserId,
			&note.BookId,
			&note.NoteText,
			&note.PageNumber,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		notes = append(notes, note)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return notes, nil
}

// Count retrieves the total number of notes associated with a specific user ID and returns it or an error if it fails.
func (n *NoteModel) Count(userId int) (int, error) {
	var count int
	stmt := `SELECT COUNT(*) FROM notes WHERE user_id = ?`
	err := n.DB.QueryRow(stmt, userId).Scan(&count)
	if err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) {
			return 0, sqliteError
		}
		return 0, err
	}
	return count, nil
}
