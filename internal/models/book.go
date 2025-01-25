package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"log/slog"
	"time"
)

// PaginatedBooks represents a paginated collection of books with metadata such as total count and pagination details.
type PaginatedBooks struct {
	Books      []Book
	Total      int
	Page       int
	TotalPages int
	PageSize   int
}

// Book represents a literary work with details such as title, author, ISBN, and publication year.
type Book struct {
	Base
	ID              int
	Title           string
	Author          string
	ISBN            string
	PublicationYear int
	Status          string
	ImageURL        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	UserId          int
}

// BookModel represents the data structure for accessing book-related data in the database.
type BookModel struct {
	DB     *sql.DB
	Logger *slog.Logger
}

// Create inserts a new book into the database, associates it with a user, and returns the book's ID or an error.
func (m *BookModel) Create(title, author, isbn, status, imageUrl string, publicationYear, userId int) (int, error) {

	// Start a transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}

	// Prepare the statement for inserting a book
	stmt := `INSERT INTO books (title, author, isbn, publication_year, image_url) 
             VALUES (?, ?, ?, ?, ?)`

	// Execute the statement and get the result
	result, err := tx.Exec(stmt, title, author, isbn, publicationYear, imageUrl)
	if err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) && errors.Is(sqliteError.ExtendedCode, sqlite3.ErrConstraintUnique) {
			if sqliteError.Error() == "UNIQUE constraint failed: books.isbn" {
				if err = tx.Rollback(); err != nil {
					return 0, err
				}
				return 0, ErrDuplicateIsbn
			}
		} else {
			if err = tx.Rollback(); err != nil {
				return 0, err
			}
		}

		bookId, err := result.LastInsertId()
		if err != nil {
			if err = tx.Rollback(); err != nil {
				m.Logger.Error(err.Error())
			}
		}
		return int(bookId), err
	}

	// Get the ID of the newly inserted book
	bookId, err := result.LastInsertId()
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	// Create the user-book relationship
	stmt = `INSERT INTO user_books (user_id, book_id, status) VALUES (?, ?, ?)`
	_, err = tx.Exec(stmt, userId, bookId, status)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return int(bookId), nil
}

// Retrieve fetches a book by its ID from the database and returns the book or an error if not found.
func (m *BookModel) Retrieve(id int) (Book, error) {
	var book Book

	stmt := `SELECT b.id, b.title, b.author, b.isbn, b.publication_year, b.created_at, b.updated_at, b.image_url, ub.user_id, ub.status
		FROM books b
		LEFT JOIN user_books ub ON b.id = ub.book_id
        WHERE b.id = ?`

	err := m.DB.QueryRow(stmt, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.ISBN,
		&book.PublicationYear,
		&book.CreatedAt,
		&book.UpdatedAt,
		&book.ImageURL,
		&book.UserId,
		&book.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Book{}, ErrNoRecord
		} else {
			return Book{}, err
		}
	}
	return book, nil
}

// Delete removes a book from the database based on the provided book ID and user ID, verifying ownership. Returns an error if unsuccessful.
func (m *BookModel) Delete(id, userId int) error {

	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				m.Logger.Error(rbErr.Error())
			}
			return
		}
		err = tx.Commit()
	}()

	var exists bool
	err = tx.QueryRow(`SELECT EXISTS(
        SELECT 1 FROM user_books 
        WHERE book_id = ? AND user_id = ?
    )`, id, userId).Scan(&exists)

	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("book with id %d does not belong to user %d", id, userId)
	}

	result, err := tx.Exec(`DELETE FROM books WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("book with id %d not found", id)
	}

	return nil
}

// Update modifies an existing book's data in the database based on the provided ID and new field values.
// Returns ErrDuplicateIsbn if the ISBN is already in use or ErrNoRecord if no record was updated.
func (m *BookModel) Update(id int, title, author, isbn, status, imageUrl string, publicationYear int) error {
	// Start a transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			m.Logger.Error(err.Error())
		}
	}(tx)

	// Update books table
	stmt := `UPDATE books SET title = ?, author = ?, isbn = ?, publication_year = ?, image_url = ? WHERE id = ?`
	result, err := tx.Exec(stmt, title, author, isbn, publicationYear, imageUrl, id)
	if err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) && errors.Is(sqliteError.ExtendedCode, sqlite3.ErrConstraintUnique) {
			if sqliteError.Error() == "UNIQUE constraint failed: books.isbn" {
				return ErrDuplicateIsbn
			}
		}
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNoRecord
	}

	// Update user_books table
	result, err = tx.Exec(`UPDATE user_books SET status = ? WHERE book_id = ?`, status, id)
	if err != nil {
		return err
	}

	affected, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNoRecord
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// List retrieves a paginated collection of books from the database, including total count and pagination metadata.
func (m *BookModel) List(page, pageSize int) (PaginatedBooks, error) {

	var total int
	err := m.DB.QueryRow("SELECT COUNT(*) FROM books").Scan(&total)
	if err != nil {
		return PaginatedBooks{}, err
	}

	totalPages := (total + pageSize - 1) / pageSize

	if page > totalPages && totalPages > 0 {
		page = totalPages
	}

	offset := (page - 1) * pageSize

	stmt := `
        SELECT b.id, b.title, b.author, b.isbn, b.publication_year, b.created_at, b.updated_at, b.image_url , ub.user_id
        FROM books b
        LEFT JOIN user_books ub ON b.id = ub.book_id
        ORDER BY created_at DESC
        LIMIT ? OFFSET ?
    `

	rows, err := m.DB.Query(stmt, pageSize, offset)
	if err != nil {
		return PaginatedBooks{}, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			m.Logger.Error(err.Error())
		}
	}()

	var books []Book
	for rows.Next() {
		var book Book
		err = rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.ISBN,
			&book.PublicationYear,
			&book.CreatedAt,
			&book.UpdatedAt,
			&book.ImageURL,
			&book.UserId,
		)
		if err != nil {
			return PaginatedBooks{}, err
		}
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return PaginatedBooks{}, err
	}

	return PaginatedBooks{
		Books:      books,
		Total:      total,
		Page:       page,
		TotalPages: totalPages,
		PageSize:   pageSize,
	}, nil
}

// Filter retrieves books where the title, notes, or reviews contain the given search term.
func (m *BookModel) Filter(searchTerm string) ([]Book, error) {
	searchTerm = "%" + searchTerm + "%"

	stmt := `
	SELECT DISTINCT b.id, b.title, b.author, b.isbn, b.publication_year, b.created_at, b.updated_at, b.image_url
	FROM books b
	LEFT JOIN notes n ON b.id = n.book_id
	LEFT JOIN reviews r ON b.id = r.book_id
	WHERE b.title LIKE ? COLLATE NOCASE 
   		OR n.note_text LIKE ? COLLATE NOCASE 
   		OR r.review_text LIKE ? COLLATE NOCASE;
	`

	rows, err := m.DB.Query(stmt, searchTerm, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			m.Logger.Error(err.Error())
		}
	}()

	var books []Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.ISBN,
			&book.PublicationYear,
			&book.CreatedAt,
			&book.UpdatedAt,
			&book.ImageURL,
		); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

// Count retrieves the total number of book records in the database and returns the count or an error if the query fails.
func (m *BookModel) Count() (int, error) {
	var count int
	err := m.DB.QueryRow("SELECT COUNT(*) FROM books").Scan(&count)
	if err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) {
			return 0, sqliteError
		}
		return 0, err
	}
	return count, nil
}

// RetrieveRecentBooks fetches the two most recently created book records from the database and returns them or an error.
func (m *BookModel) RetrieveRecentBooks(limit int) ([]Book, error) {
	var books []Book
	stmt := `SELECT id, title, author, isbn, publication_year, created_at, updated_at, image_url 
			FROM books ORDER BY created_at DESC LIMIT ?`

	rows, err := m.DB.Query(stmt, limit)
	if err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) {
			return nil, sqliteError
		}
		return nil, err
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			m.Logger.Error(err.Error())
		}
	}()

	for rows.Next() {
		var book Book
		err = rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.ISBN,
			&book.PublicationYear,
			&book.CreatedAt,
			&book.UpdatedAt,
			&book.ImageURL)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}

// CountFinishedBooks returns the number of books marked as 'finished' by a specific user, identified by userId.
func (m *BookModel) CountFinishedBooks(userId int) (int, error) {
	var count int
	stmt := `SELECT COUNT(*) FROM user_books WHERE user_id = ? AND status = 'finished'`
	err := m.DB.QueryRow(stmt, userId).Scan(&count)
	if err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) {
			return 0, sqliteError
		}
		return 0, err
	}
	return count, nil
}
