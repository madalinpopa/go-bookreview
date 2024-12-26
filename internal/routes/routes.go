// Package routes provided the application routes
package routes

import (
	"github.com/justinas/alice"
	"github.com/madalinpopa/go-bookreview/internal/app"
	"github.com/madalinpopa/go-bookreview/internal/middleware"
	"github.com/madalinpopa/go-bookreview/internal/views"
	"github.com/madalinpopa/go-bookreview/ui"
	"net/http"
)

// UrlPatterns sets up and returns an HTTP handler with predefined routing patterns for the given application.
func UrlPatterns(app *app.App) http.Handler {

	// Create a new HTTP request multiplexer to manage URL-to-handler routing.
	mux := http.NewServeMux()

	// Instantiate the middleware chain for handling requests with recovery capabilities.
	m := middleware.NewMiddleware(app)

	// Setup dynamic middleware
	dynamic := alice.New(app.SessionManager.LoadAndSave, m.CsrfToken, m.Authenticate, m.RedirectAuthenticatedUsers)

	// Serve static and media files under /static/ and /uploads/ path.
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))
	mux.Handle("GET /uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(app.Config.UploadDir))))

	// Public routes for pages
	mux.Handle("GET /{$}", dynamic.Then(views.HomePage(app)))
	mux.Handle("GET /login", dynamic.Then(views.LoginPage(app)))
	mux.Handle("POST /login", dynamic.Then(views.LoginPost(app)))
	mux.Handle("POST /logout", dynamic.Then(views.LogoutPost(app)))
	mux.Handle("GET /register", dynamic.Then(views.RegisterPage(app)))
	mux.Handle("POST /register", dynamic.Then(views.RegisterPost(app)))
	mux.Handle("GET /books", dynamic.Then(views.BooksPage(app)))
	mux.Handle("GET /books/{id}", dynamic.Then(views.BooksDetailPage(app)))
	mux.Handle("GET /books/{id}/reviews", dynamic.Then(views.ListReviews(app)))
	mux.Handle("GET /books/{id}/notes", dynamic.Then(views.ListNotes(app)))

	// Public routes for HTMX
	mux.Handle("GET /api/search", dynamic.Then(views.GetFilteredBooks(app)))
	mux.Handle("GET /api/books/count", dynamic.Then(views.GetBooksCount(app)))
	mux.Handle("GET /api/reviews/count", dynamic.Then(views.GetReviewsCount(app)))
	mux.Handle("GET /api/notes/count", dynamic.Then(views.GetNotesCount(app)))
	mux.Handle("GET /api/books/recent", dynamic.Then(views.GetRecentBooks(app)))
	mux.Handle("GET /api/books/read", dynamic.Then(views.GetFinishedBooks(app)))
	mux.Handle("GET /api/reviews/recent", dynamic.Then(views.GetRecentReviews(app)))

	protected := dynamic.Append(m.LoginRequired)

	mux.Handle("GET /books/new", protected.Then(views.BooksAddPage(app)))
	mux.Handle("POST /books/new", protected.Then(views.CreateBookPost(app)))
	mux.Handle("GET /books/{id}/edit", protected.Then(views.UpdateBookPage(app)))
	mux.Handle("POST /books/{id}/edit", protected.Then(views.UpdateBookPost(app)))
	mux.Handle("POST /books/delete", protected.Then(views.DeleteBookPost(app)))
	mux.Handle("GET /books/{id}/review/new", protected.Then(views.CreateReview(app)))
	mux.Handle("POST /books/review/new", protected.Then(views.CreateReviewPost(app)))
	mux.Handle("GET /books/review/{id}/edit", protected.Then(views.UpdateReview(app)))
	mux.Handle("POST /books/review/edit", protected.Then(views.UpdateReviewPost(app)))
	mux.Handle("POST /books/review/delete", protected.Then(views.DeleteReviewPost(app)))
	mux.Handle("GET /books/{id}/note/new", protected.Then(views.CreateNote(app)))
	mux.Handle("POST /books/note/new", protected.Then(views.CreateNotePost(app)))
	mux.Handle("GET /books/note/{id}/edit", protected.Then(views.UpdateNote(app)))
	mux.Handle("POST /books/note/edit", protected.Then(views.UpdateNotePost(app)))
	mux.Handle("POST /books/note/delete", protected.Then(views.DeleteNotePost(app)))

	// Setup standard middleware
	standardMiddleware := alice.New(m.Recover, m.Logging, m.Headers)

	return standardMiddleware.Then(mux)
}
