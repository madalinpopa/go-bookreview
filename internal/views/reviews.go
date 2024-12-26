package views

import (
	"errors"
	"fmt"
	"github.com/madalinpopa/go-bookreview/internal/app"
	"github.com/madalinpopa/go-bookreview/internal/forms"
	"github.com/madalinpopa/go-bookreview/internal/models"
	"net/http"
	"strconv"
)

// CreateReview handles the creation of a book review by rendering the appropriate template with book-related data.
// It retrieves the book's information using its ID from the URL, ensuring validation and error handling for possible failures.
func CreateReview(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form forms.BookReviewForm

		bookId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			app.ClientError(w, r, http.StatusBadRequest, err)
			return
		}

		book, err := app.Models.Books.Retrieve(bookId)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.ClientError(w, r, http.StatusNotFound, err)
				return
			}
			app.ServerError(w, r, err)
			return
		}

		data := app.GetTemplateData(r)
		data.Book = book
		data.Form = form
		app.Render(w, r, "htmxBookReviewForm", data, http.StatusOK)
	}
}

func CreateReviewPost(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()
		if err != nil {
			app.ClientError(w, r, http.StatusBadRequest, err)
			return
		}

		var form forms.BookReviewForm
		if err := app.FormDecoder.Decode(&form, r.PostForm); err != nil {
			app.ClientError(w, r, http.StatusBadRequest, err)
			return
		}

		if !form.Valid() {
			data := app.GetTemplateData(r)
			data.Form = form
			app.Render(w, r, "htmxBookReviewForm", data, http.StatusUnprocessableEntity)
		}

		book, err := app.Models.Books.Retrieve(form.BookId)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.ClientError(w, r, http.StatusNotFound, err)
				return
			}
			app.ServerError(w, r, err)
			return
		}
		userId := app.GetAuthenticatedUserId(r)
		if userId == 0 {
			app.ClientError(w, r, http.StatusUnauthorized, errors.New("user not authenticated"))
			return
		}

		_, err = app.Models.Reviews.Create(userId, book.ID, form.Rating, form.ReviewText)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		w.Header().Set("HX-Trigger", "update-reviews")
		w.WriteHeader(http.StatusNoContent)

	}
}

// UpdateReview handles a request to update an existing book review for an authenticated user.
// It fetches the review and associated book, validates the user's access, and renders the review update form.
func UpdateReview(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form forms.BookReviewForm

		reviewId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			app.ClientError(w, r, http.StatusBadRequest, err)
			return
		}

		userId := app.GetAuthenticatedUserId(r)
		if userId == 0 {
			app.ClientError(w, r, http.StatusUnauthorized, errors.New("user not authenticated"))
			return
		}

		review, err := app.Models.Reviews.Retrieve(reviewId, userId)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.ClientError(w, r, http.StatusNotFound, err)
				return
			}
			app.ServerError(w, r, err)
			return
		}

		book, err := app.Models.Books.Retrieve(review.BookId)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.ClientError(w, r, http.StatusNotFound, err)
				return
			}
			app.ServerError(w, r, err)
			return
		}

		form.Id = review.ID
		form.Rating = review.Rating
		form.ReviewText = review.ReviewText

		data := app.GetTemplateData(r)
		data.Review = review
		data.Form = form
		data.Book = book
		app.Render(w, r, "htmxBookReviewForm", data, http.StatusOK)
	}
}

// UpdateReviewPost handles the updating of a book review by parsing form data, validating input, and saving changes to the database.
func UpdateReviewPost(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			app.ClientError(w, r, http.StatusBadRequest, err)
			return
		}

		var form forms.BookReviewForm
		if err := app.FormDecoder.Decode(&form, r.PostForm); err != nil {
			app.ClientError(w, r, http.StatusBadRequest, err)
			return
		}

		if !form.Valid() {
			data := app.GetTemplateData(r)
			data.Form = form
			app.Render(w, r, "htmxBookReviewForm", data, http.StatusUnprocessableEntity)
		}

		userId := app.GetAuthenticatedUserId(r)
		if userId == 0 {
			app.ClientError(w, r, http.StatusUnauthorized, errors.New("user not authenticated"))
			return
		}

		review, err := app.Models.Reviews.Retrieve(form.Id, userId)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.ClientError(w, r, http.StatusNotFound, err)
				return
			}
			app.ServerError(w, r, err)
			return
		}

		err = app.Models.Reviews.Update(review.ID, review.UserId, form.Rating, form.ReviewText)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.ClientError(w, r, http.StatusNotFound, err)
				return
			}
			app.ServerError(w, r, err)
			return
		}

		url := fmt.Sprintf("/books/%d", review.BookId)
		app.HtmxLocation(w, r, url, "#books-content", "innerHTML")
	}
}

// DeleteReviewPost handles the deletion of a review post based on the review ID and authenticated user.
// It validates the review ID, checks user authentication, and performs removal if the user is authorized.
// If successful, it redirects to the book's detailed page while updating the target content dynamically.
func DeleteReviewPost(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			app.ClientError(w, r, http.StatusBadRequest, err)
			return
		}

		var form forms.BookReviewForm
		if err := app.FormDecoder.Decode(&form, r.PostForm); err != nil {
			app.ClientError(w, r, http.StatusBadRequest, err)
			return
		}

		userId := app.GetAuthenticatedUserId(r)
		if userId == 0 {
			app.ClientError(w, r, http.StatusUnauthorized, errors.New("user not authenticated"))
			return
		}

		review, err := app.Models.Reviews.Retrieve(form.Id, userId)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.ClientError(w, r, http.StatusNotFound, err)
				return
			}
			app.ServerError(w, r, err)
			return
		}

		err = app.Models.Reviews.Delete(review.ID, review.UserId)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}
	}
}

// ListReviews handles HTTP requests for listing reviews for a specific book by its ID and renders the result as a response.
func ListReviews(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bookId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		book, err := app.Models.Books.Retrieve(bookId)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.ClientError(w, r, http.StatusNotFound, err)
				return
			}
			app.ServerError(w, r, err)
			return
		}

		data := app.GetTemplateData(r)
		data.Book = book

		reviews, err := app.Models.Reviews.List(book.ID)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}
		data.Reviews = reviews
		app.Render(w, r, "htmxBookReviews", data, http.StatusOK)
	}
}

// GetReviewsCount handles HTTP requests to fetch and return the total count of reviews in the database as a plain text response.
func GetReviewsCount(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userId := app.GetAuthenticatedUserId(r)
		if userId == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		count, err := app.Models.Reviews.Count(userId)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		_, err = fmt.Fprintf(w, "%d", count)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}
	}
}

// GetRecentReviews handles HTTP requests to retrieve and display the most recent reviews in the application.
func GetRecentReviews(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reviews, err := app.Models.Reviews.RetrieveRecentReviews(2)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		data := app.GetTemplateData(r)
		data.Reviews = reviews
		app.Render(w, r, "htmxRecentReviews", data, http.StatusOK)
	}
}
