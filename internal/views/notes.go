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

// CreateNote handles the display of a form for creating notes associated with a specific book, identified by its ID.
// It retrieves the book using the provided ID, validates input, and sends the form template in the response.
// If the book is not found or input validation fails, appropriate error responses are returned.
func CreateNote(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form forms.BookNoteForm

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
		app.Render(w, r, "htmxBookNoteForm", data, http.StatusOK)
	}
}

// CreateNotePost handles an HTTP POST request to create a new note associated with a book for the authenticated user.
// It parses and validates the form data, checks user authentication, retrieves the book, creates the note,
// and responds with a status indicating success or error. Adds an "HX-Trigger" header for HTMX requests if successful.
func CreateNotePost(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			app.ClientError(w, r, http.StatusBadRequest, err)
			return
		}

		var form forms.BookNoteForm
		if err := app.FormDecoder.Decode(&form, r.PostForm); err != nil {
			app.ClientError(w, r, http.StatusBadRequest, err)
		}

		if !form.Valid() {
			data := app.GetTemplateData(r)
			data.Form = form
			app.Render(w, r, "htmxBookNoteForm", data, http.StatusUnprocessableEntity)
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

		_, err = app.Models.Notes.Create(userId, book.ID, form.NoteText, form.PageNumber)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		w.Header().Set("HX-Trigger", "update-notes")
		w.WriteHeader(http.StatusNoContent)
	}
}

// UpdateNote handles the HTTP request to display the update form for a specific book note based on the provided note ID.
// It validates the note ID, retrieves the associated note and book details, and renders the HTMX update form.
func UpdateNote(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form forms.BookNoteForm
		noteId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			app.ClientError(w, r, http.StatusBadRequest, err)
			return
		}

		userId := app.GetAuthenticatedUserId(r)
		if userId == 0 {
			app.ClientError(w, r, http.StatusUnauthorized, errors.New("user not authenticated"))
		}

		note, err := app.Models.Notes.Retrieve(userId, noteId)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.ClientError(w, r, http.StatusNotFound, err)
				return
			}
			app.ServerError(w, r, err)
			return
		}
		book, err := app.Models.Books.Retrieve(note.BookId)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.ClientError(w, r, http.StatusNotFound, err)
			}
		}

		form.Id = note.ID
		form.NoteText = note.NoteText
		form.PageNumber = note.PageNumber
		data := app.GetTemplateData(r)
		data.Note = note
		data.Form = form
		data.Book = book
		app.Render(w, r, "htmxBookNoteForm", data, http.StatusOK)
	}
}

// UpdateNotePost handles HTTP POST requests to update an existing note, performing validation and authentication checks.
// It retrieves the note by ID, validates user ownership, updates its content, and sends an appropriate HTTP response.
func UpdateNotePost(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			app.ClientError(w, r, http.StatusBadRequest, err)
			return
		}

		var form forms.BookNoteForm
		if err := app.FormDecoder.Decode(&form, r.PostForm); err != nil {
			app.ClientError(w, r, http.StatusBadRequest, err)
			return
		}

		if !form.Valid() {
			data := app.GetTemplateData(r)
			data.Form = form
			app.Render(w, r, "htmxBookNoteForm", data, http.StatusUnprocessableEntity)
		}

		userId := app.GetAuthenticatedUserId(r)
		if userId == 0 {
			app.ClientError(w, r, http.StatusUnauthorized, errors.New("user not authenticated"))
			return
		}

		note, err := app.Models.Notes.Retrieve(form.Id, userId)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.ClientError(w, r, http.StatusNotFound, err)
				return
			}
			app.ServerError(w, r, err)
			return
		}

		err = app.Models.Notes.Update(note.ID, note.UserId, form.NoteText, form.PageNumber)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.ClientError(w, r, http.StatusNotFound, err)
				return
			}
			app.ServerError(w, r, err)
			return
		}

		w.Header().Set("HX-Trigger", "update-notes")
		w.WriteHeader(http.StatusNoContent)
	}
}

// DeleteNotePost handles the deletion of a note based on the provided form data and the authenticated user ID.
// Validates request form data and user authentication before deleting the note.
// Responds with appropriate HTTP status codes for errors such as unauthorized access, bad request, or not found.
func DeleteNotePost(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			app.ClientError(w, r, http.StatusBadRequest, err)
			return
		}

		var form forms.BookNoteForm
		if err := app.FormDecoder.Decode(&form, r.PostForm); err != nil {
			app.ClientError(w, r, http.StatusBadRequest, err)
		}

		userId := app.GetAuthenticatedUserId(r)
		if userId == 0 {
			app.ClientError(w, r, http.StatusUnauthorized, errors.New("user not authenticated"))
			return
		}

		note, err := app.Models.Notes.Retrieve(userId, form.Id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.ClientError(w, r, http.StatusNotFound, err)
				return
			}
			app.ServerError(w, r, err)
			return
		}

		err = app.Models.Notes.Delete(note.UserId, note.ID)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}
	}
}

// ListNotes returns an HTTP handler that lists all notes associated with a specific book for the provided book ID.
// It retrieves the book and its notes, populating the template data, and renders the "htmxBookNotes" template.
func ListNotes(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

		userId := app.GetAuthenticatedUserId(r)

		notes, err := app.Models.Notes.List(book.ID, userId)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}
		data := app.GetTemplateData(r)
		data.Notes = notes
		data.Book = book
		app.Render(w, r, "htmxBookNotes", data, http.StatusOK)
	}
}

// GetNotesCount handles HTTP requests to retrieve the total count of notes from the application's data model.
func GetNotesCount(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := app.GetAuthenticatedUserId(r)
		if userId == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		count, err := app.Models.Notes.Count(userId)
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
