package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"log/slog"
	"time"
)

// Review represents a user's review of a book, including their rating, user ID, book ID, and optional review text.
type Review struct {
	Base
	UserId     int
	BookId     int
	Rating     int
	ReviewText string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Username   string
	BookTitle  string
}

// ReviewModel provides methods to interact with the reviews data in the database.
type ReviewModel struct {
	DB     *sql.DB
	Logger *slog.Logger
}

// Create inserts a new review into the database and returns the ID of the created review or an error if the operation fails.
func (m *ReviewModel) Create(userId, bookId, rating int, reviewText string) (int, error) {

	stmt := `INSERT INTO reviews (user_id, book_id, rating, review_text) VALUES (?, ?, ?, ?)`

	// Execute the statement and get the result
	result, err := m.DB.Exec(stmt, userId, bookId, rating, reviewText)
	if err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) {
			return 0, sqliteError
		} else {
			return 0, err
		}
	}

	reviewId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(reviewId), nil
}

// Retrieve fetches a review by its ID and associated user ID from the database. Returns the review or an error if not found.
func (m *ReviewModel) Retrieve(id, userId int) (Review, error) {
	var review Review
	stmt := `SELECT id, user_id, book_id, rating, review_text, created_at, updated_at FROM reviews WHERE id = ? AND user_id = ?`

	err := m.DB.QueryRow(stmt, id, userId).Scan(
		&review.ID,
		&review.UserId,
		&review.BookId,
		&review.Rating,
		&review.ReviewText,
		&review.CreatedAt,
		&review.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Review{}, ErrNoRecord
		} else {
			return Review{}, err
		}
	}
	return review, nil
}

// Update modifies the rating and review text of an existing review for a specific user and book in the database.
// Returns an error if the update fails or no matching record is found.
func (m *ReviewModel) Update(id, userId int, rating int, reviewText string) error {
	stmt := `UPDATE reviews SET rating = ?, review_text = ? WHERE id = ? AND user_id = ?`

	result, err := m.DB.Exec(stmt, rating, reviewText, id, userId)
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

// Delete removes a review from the database for a specified ID and userID. Returns an error if no matching record is found.
func (m *ReviewModel) Delete(id, userId int) error {
	stmt := `DELETE FROM reviews WHERE id = ? AND user_id = ?`

	result, err := m.DB.Exec(stmt, id, userId)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("review with id %d and user_id %d not found", id, userId)
	}

	return nil
}

// List retrieves all reviews for a specified book ID from the database and returns them or an error if the query fails.
func (m *ReviewModel) List(bookId int) ([]Review, error) {

	stmt := `
        SELECT r.id, r.user_id, r.book_id, r.rating, r.review_text, r.created_at, r.updated_at, u.username 
        FROM reviews r 
        LEFT JOIN users u ON r.user_id = u.id 
        WHERE r.book_id = ?
    `

	rows, err := m.DB.Query(stmt, bookId)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			m.Logger.Error(err.Error())
		}
	}()

	var reviews []Review
	for rows.Next() {
		var review Review
		err = rows.Scan(
			&review.ID,
			&review.UserId,
			&review.BookId,
			&review.Rating,
			&review.ReviewText,
			&review.CreatedAt,
			&review.UpdatedAt,
			&review.Username,
		)
		reviews = append(reviews, review)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

// Count returns the total number of reviews associated with a specific user ID or an error if the query fails.
func (m *ReviewModel) Count(userId int) (int, error) {
	var count int
	stmt := `SELECT COUNT(*) FROM reviews WHERE user_id = ?`
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

// RetrieveRecentReviews fetches the most recent reviews up to a specified limit, ordered by creation date in descending order.
func (m *ReviewModel) RetrieveRecentReviews(limit int) ([]Review, error) {
	stmt := `SELECT r.id, r.user_id, r.book_id, r.rating, r.review_text, b.title 
        FROM reviews r 
        JOIN books b ON r.book_id = b.id 
        ORDER BY r.created_at DESC 
        LIMIT ?`

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

	var reviews []Review
	for rows.Next() {
		var review Review
		err = rows.Scan(
			&review.ID,
			&review.UserId,
			&review.BookId,
			&review.Rating,
			&review.ReviewText,
			&review.BookTitle,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reviews, nil
}
