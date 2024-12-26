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

// BooksPage handles HTTP requests to display a paginated list of books using the given app's data and templates.
func BooksPage(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := app.GetTemplateData(r)

		page := 1
		if pageStr := r.URL.Query().Get("page"); pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}

		pageSize := 8
		paginated, err := app.Models.Books.List(page, pageSize)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		data.Books = paginated.Books
		data.Page = paginated.Page
		data.PageSize = paginated.PageSize
		data.Total = paginated.Total
		data.TotalPages = paginated.TotalPages

		if app.IsHtmxRequest(r) {
			app.Render(w, r, "htmxBookCard", data, http.StatusOK)
			return
		}
		app.Render(w, r, "books.tmpl", data, http.StatusOK)
	}
}

// BooksAddPage renders the "htmxCreateBook" template with the provided request-specific data using a 200 OK status.
func BooksAddPage(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form forms.BookForm
		data := app.GetTemplateData(r)
		data.Form = form
		if app.IsHtmxRequest(r) {
			app.Render(w, r, "htmxCreateBook", data, http.StatusOK)
			return
		}
		app.Render(w, r, "books_add.tmpl", data, http.StatusOK)
	}
}

// BooksDetailPage handles requests for the book detail page
// by retrieving a book record and rendering the appropriate template.
func BooksDetailPage(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.NotFound(w, r)
			return
		}

		book, err := app.Models.Books.Retrieve(id)
		if err != nil {
			app.ServerError(w, r, err)
		}

		data := app.GetTemplateData(r)
		data.Book = book
		if app.IsHtmxRequest(r) {
			app.Render(w, r, "htmxBookDetail", data, http.StatusOK)
			return
		}
		app.Render(w, r, "books_detail.tmpl", data, http.StatusOK)
	}

}

// CreateBookPost handles HTTP POST requests for creating a new book with supplied form data.
func CreateBookPost(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Limit the request body to 5MB and parse the multipart form
		r.Body = http.MaxBytesReader(w, r.Body, 5<<20)

		// Ensure the request body is limited to prevent excessively large uploads
		if err := r.ParseMultipartForm(5 << 20); err != nil {
			app.ClientError(w, r, http.StatusBadRequest, err)
			return
		}
		defer func() {
			if err := r.MultipartForm.RemoveAll(); err != nil {
				app.Logger.Error(err.Error())
			}
		}()

		var form forms.BookForm
		if err := app.FormDecoder.Decode(&form, r.PostForm); err != nil {
			app.ClientError(w, r, http.StatusBadRequest, err)
			return
		}

		form.Validate()
		if !form.Valid() {
			data := app.GetTemplateData(r)
			data.Form = form
			app.Logger.Error("form validation failed", "valid", form.Valid())
			app.Render(w, r, "htmxBookForm", data, http.StatusUnprocessableEntity)
			return
		}

		// Handle file upload
		err := form.HandleFileUpload(app, r)
		if err != nil {
			if errors.Is(err, forms.ErrFormBadRequest) {
				app.Logger.Error("form bad request", "error", err)
				app.ClientError(w, r, http.StatusBadRequest, err)
				return
			} else if errors.Is(err, forms.ErrInvalidFileType) {
				app.Logger.Error("invalid file type", "error", err)
				data := app.GetTemplateData(r)
				data.Form = form
				app.Render(w, r, "htmxBookForm", data, http.StatusUnprocessableEntity)
				return
			} else {
				app.ServerError(w, r, err)
				return
			}
		}
		userId := app.GetAuthenticatedUserId(r)
		if userId == 0 {
			app.ClientError(w, r, http.StatusUnauthorized, errors.New("user not authenticated"))
			return
		}

		bookId, err := app.Models.Books.Create(form.Title, form.Author, form.ISBN, form.Status, form.ImageURL, form.PublicationYear, userId)
		if err != nil {
			if errors.Is(err, models.ErrDuplicateIsbn) {
				form.AddFieldError("isbn", "This ISBN is already registered.")
				data := app.GetTemplateData(r)
				data.Form = form
				app.Render(w, r, "htmxBookForm", data, http.StatusUnprocessableEntity)
			} else {
				app.ServerError(w, r, err)
			}
			return
		}

		url := fmt.Sprintf("/books/%d", bookId)
		app.HtmxLocation(w, r, url, "#books-content", "innerHTML")

	}
}

// UpdateBookPage handles HTTP requests to render the book update page, populating form and book data from the database.
func UpdateBookPage(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bookId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		var form forms.BookForm
		data := app.GetTemplateData(r)
		book, err := app.Models.Books.Retrieve(bookId)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		userId := app.GetAuthenticatedUserId(r)
		if userId == 0 {
			app.ClientError(w, r, http.StatusUnauthorized, errors.New("user not authenticated"))
			return
		}

		data.Book = book
		data.Form = form
		if app.IsHtmxRequest(r) {
			app.Render(w, r, "htmxBookUpdate", data, http.StatusOK)
			return
		}
		app.Render(w, r, "books_update.tmpl", data, http.StatusOK)
	}
}

// UpdateBookPost handles the update of a book record via an HTTP POST request.
// Validates input, processes file uploads, and updates the book in the database.
// Sends appropriate responses on success or error.
func UpdateBookPost(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := app.GetTemplateData(r)

		// Limit the request body to 5MB and parse the multipart form
		r.Body = http.MaxBytesReader(w, r.Body, 5<<20)

		// Ensure the request body is limited to prevent excessively large uploads
		if err := r.ParseMultipartForm(5 << 20); err != nil {
			app.Logger.Error("form parsing failed", "error", err)
			app.ClientError(w, r, http.StatusBadRequest, err)
			return
		}
		defer func() {
			if err := r.MultipartForm.RemoveAll(); err != nil {
				app.Logger.Error(err.Error())
			}
		}()

		var form forms.BookForm
		if err := app.FormDecoder.Decode(&form, r.PostForm); err != nil {
			app.Logger.Error("form decoding failed", "error", err)
			app.ClientError(w, r, http.StatusBadRequest, err)
			return
		}

		form.Validate()
		if !form.Valid() {
			app.Logger.Error("form validation failed", "valid", form.Valid())
			data := app.GetTemplateData(r)
			data.Form = form
			app.Render(w, r, "htmxBookEdit", data, http.StatusUnprocessableEntity)
			return
		}

		// Handle file upload
		err := form.HandleFileUpload(app, r)
		if err != nil {
			if errors.Is(err, forms.ErrFormBadRequest) {
				app.Logger.Error("form bad request", "error", err)
				app.ClientError(w, r, http.StatusBadRequest, err)
				return
			} else if errors.Is(err, forms.ErrInvalidFileType) {
				app.Logger.Error("invalid file type", "error", err)
				data := app.GetTemplateData(r)
				data.Form = form
				app.Render(w, r, "htmxBookEdit", data, http.StatusUnprocessableEntity)
			} else {
				app.ServerError(w, r, err)
				return
			}
		}

		userId := app.GetAuthenticatedUserId(r)
		if userId == 0 {
			app.Logger.Error("user not authenticated")
			app.ClientError(w, r, http.StatusUnauthorized, errors.New("user not authenticated"))
			return
		}

		bookId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.NotFound(w, r)
			return
		}

		_, err = app.Models.Books.Retrieve(bookId)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.ClientError(w, r, http.StatusNotFound, err)
				app.Logger.Error("book not found", "id", bookId)
				return
			} else {
				app.ServerError(w, r, err)

			}
			return
		}

		err = app.Models.Books.Update(bookId, form.Title, form.Author, form.ISBN, form.Status, form.ImageURL, form.PublicationYear)
		if err != nil {
			if errors.Is(err, models.ErrDuplicateIsbn) {
				form.AddFieldError("isbn", "This ISBN is already registered.")
				data.Form = form
				app.Render(w, r, "htmxBookEdit", data, http.StatusUnprocessableEntity)
				app.Logger.Error("duplicate ISBN", "isbn", form.ISBN)
				return
			} else {
				app.ServerError(w, r, err)
			}
			return
		}

		url := fmt.Sprintf("/books/%d", bookId)
		app.HtmxLocation(w, r, url, "#books-content", "innerHTML")

	}
}

// DeleteBookPost handles the deletion of a book post by decoding form data and invoking the Books model's Delete method.
// It returns an appropriate response or error based on the outcome of the operation.
func DeleteBookPost(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			app.ClientError(w, r, http.StatusBadRequest, err)
			app.Logger.Error("form parsing failed", "error", err)
			return
		}

		var form forms.BookForm
		if err := app.FormDecoder.Decode(&form, r.PostForm); err != nil {
			app.ClientError(w, r, http.StatusBadRequest, err)
			app.Logger.Error("form decoding failed", "error", err)
			return
		}
		userId := app.GetAuthenticatedUserId(r)
		if userId == 0 {
			app.ClientError(w, r, http.StatusUnauthorized, errors.New("user not authenticated"))
			app.Logger.Error("user not authenticated")
			return
		}

		err = app.Models.Books.Delete(form.Id, userId)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}
		app.HtmxLocation(w, r, "/books", "#books-content", "innerHTML")
	}
}

// GetFilteredBooks handles requests to search for books using a search term and renders the filtered results as HTML templates.
// It parses form data, retrieves matching books from the application's data model, and handles errors accordingly.
func GetFilteredBooks(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			app.ClientError(w, r, http.StatusBadRequest, err)
			return
		}

		searchTerm := r.FormValue("search")

		books, err := app.Models.Books.Filter(searchTerm)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}
		data := app.GetTemplateData(r)
		data.Books = books
		app.Render(w, r, "htmxBookCard", data, http.StatusOK)
	}
}

// GetBooksCount handles HTTP requests to retrieve the count of books from the database and responds in JSON format.
func GetBooksCount(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		count, err := app.Models.Books.Count()
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

// GetRecentBooks returns an HTTP handler function that retrieves and renders the two most recent books from the database.
func GetRecentBooks(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		books, err := app.Models.Books.RetrieveRecentBooks(2)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		data := app.GetTemplateData(r)
		data.Books = books
		app.Render(w, r, "htmxRecentBooks", data, http.StatusOK)

	}
}

// GetFinishedBooks handles HTTP requests to retrieve the count of finished books for an authenticated user.
// It writes the book count as a response or a 204 status if the user is not authenticated.
func GetFinishedBooks(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userId := app.GetAuthenticatedUserId(r)
		if userId == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		count, err := app.Models.Books.CountFinishedBooks(userId)
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
