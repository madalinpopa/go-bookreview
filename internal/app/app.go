package app

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
	"github.com/madalinpopa/go-bookreview/internal/models"
	"github.com/madalinpopa/go-bookreview/ui"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"
)

// contextKey is a custom type used for defining strongly typed keys in context to avoid collision and improve clarity.
type contextKey string

// IsAuthenticatedContextKey is a context key used to store and retrieve authentication status in a strongly typed manner.
const IsAuthenticatedContextKey = contextKey("isAuthenticated")

// functions is a template.FuncMap providing custom date formatting functions for use in HTML templates.
var functions = template.FuncMap{
	"humanDate": func(t time.Time) string {
		if t.IsZero() {
			return ""
		}
		return t.Format("02 Jan 2006 at 15:04")
	},
	"formatDate": func(t time.Time) string {
		if t.IsZero() {
			return ""
		}
		return t.Format("2006-01-02")
	},
	"iterate": func(count int) []int {
		var numbers []int
		for i := 0; i < count; i++ {
			numbers = append(numbers, i)
		}
		return numbers
	},
	"add": func(a, b int) int {
		return a + b
	},
	"sub": func(a, b int) int {
		return a - b
	},
	"mul": func(a, b int) int {
		return a * b
	},
	"min": func(a, b int) int {
		if a < b {
			return a
		}
		return b
	},
}

// TemplateData holds data passed to templates, including form state, page title, and CSRF token for security.
type TemplateData struct {

	// Form represents the form state or data passed to templates, which can be of any type.
	Form any

	// Title represents the title of the current page or template being rendered.
	Title string

	// CSRFToken is a string containing a CSRF (Cross-Site Request Forgery) token used to secure forms from CSRF attacks.
	CSRFToken string

	// Flash is a string used to store one-time messages, such as notifications or alerts,
	//intended for display to the user.
	Flash string

	// CurrentYear represents the current year, typically used for displaying dynamic year information in templates.
	CurrentYear int

	// IsAuthenticated indicates whether the user is authenticated, typically determined by session or request context data.
	IsAuthenticated bool

	// Username represents the authenticated user's name typically used in templates or session data.
	Username string

	// AuthenticatedUserId represents the ID of the currently authenticated user, typically retrieved from session data.
	AuthenticatedUserId int

	// Books is a slice of Book objects representing a collection of literary works that can be passed to templates.
	Books []models.Book

	// Book represents a single Book item to be passed to templates, containing details such as title, author, and ISBN.
	Book models.Book

	// Notes is a slice of Note objects representing user-created notes associated with books and specific pages.
	Notes []models.Note

	// Note represents a single user-created note containing information like text, page number, and associated identifiers.
	Note models.Note

	// Reviews is a slice of Review objects representing users' reviews of books, typically including ratings and text.
	Reviews []models.Review

	// Review holds a single review object, including its rating, user ID, book ID, and optional review text.
	Review models.Review

	// Page is the current page number for paginated data or content sections.
	Page int

	// PageSize specifies the number of items displayed per page in paginated data.
	PageSize int

	// Total represents the total number of records or items available for a specific query or list.
	Total int

	// TotalPages represents the total number of pages available based on the total items and current pagination settings.
	TotalPages int
}

// App represents the core application structure including database, configuration, and logging layout.
type App struct {

	// db represents the application's primary database connection pool.
	db *sql.DB

	// config holds the application's configuration settings including address, database connection string, and port details.
	Config *Config

	// Models provides centralized access to the application's aggregated data models, such as UserModel and others.
	Models *models.Models

	// Logger is used for capturing and managing application log messages and events.
	Logger *slog.Logger

	// serverError is a helper function for logging internal server errors and responding with a 500 status.

	// templates is a map for caching parsed templates, allowing efficient rendering and reuse within the application.
	templates map[string]*template.Template

	// SessionManager manages user sessions, providing methods for creating, retrieving, and deleting session data.
	SessionManager *scs.SessionManager

	// FormDecoder is used to decode form values into Go structs, supporting custom type decoding and validation.
	FormDecoder *form.Decoder
}

// NewApp initializes and returns a pointer to an App struct with the provided database connection and configuration.
func NewApp(config *Config, db *sql.DB) *App {

	// Initialize session manager with SQLite3-backed session store
	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// Register custom type function for decoding time.Time fields in forms
	formDecoder := form.NewDecoder()
	formDecoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		return time.Parse("2006-01-02", vals[0])
	}, time.Time{})

	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	return &App{
		db:             db,
		Config:         config,
		Logger:         logger,
		FormDecoder:    formDecoder,
		SessionManager: sessionManager,
		Models:         models.NewModels(db, logger),
	}
}

// IsAuthenticated checks
// if a user is authenticated by verifying the presence of session or authentication data in the request.
func (a *App) IsAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(IsAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

// GetAuthenticatedUserId retrieves the authenticated user ID from the session data using the HTTP request context.
func (a *App) GetAuthenticatedUserId(r *http.Request) int {
	userId := a.SessionManager.GetInt(r.Context(), "authenticatedUserID")
	return userId
}

// GetAuthenticatedUserName retrieves the authenticated user's name from the session data using the HTTP request context.
func (a *App) GetAuthenticatedUserName(r *http.Request) string {
	userName := a.SessionManager.GetString(r.Context(), "authenticatedUsername")
	return userName
}

// LoadTemplates parses HTML templates and stores them in a cache for efficient rendering during runtime.
func (a *App) LoadTemplates() error {
	cache := map[string]*template.Template{}

	baseTmpl, err := template.New("base").Funcs(functions).ParseFS(ui.Files,
		"html/base.tmpl",
		"html/partials/layout/*.tmpl",
		"html/partials/components/*.tmpl",
	)
	if err != nil {
		return fmt.Errorf("parsing base templates: %w", err)
	}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return fmt.Errorf("reading page templates: %w", err)
	}

	for _, page := range pages {
		name := filepath.Base(page) // Extract the filename

		tmpl, err := baseTmpl.Clone()
		if err != nil {
			return fmt.Errorf("cloning base template: %w", err)
		}

		tmpl, err = tmpl.ParseFS(ui.Files, page)
		if err != nil {
			return fmt.Errorf("parsing %s template: %w", name, err)
		}

		for _, t := range tmpl.Templates() {
			if t.Name() != "" {
				cache[t.Name()] = tmpl
			}
		}
	}

	a.templates = cache
	return nil
}

// Render processes and renders an HTML template with the passed data and HTTP status code to the response writer.
func (a *App) Render(w http.ResponseWriter, r *http.Request, name string, data interface{}, status int) {
	t, ok := a.templates[name]
	if !ok {
		err := fmt.Errorf("template %s not found", name)
		a.ServerError(w, r, err)
		return
	}
	buf := new(bytes.Buffer)
	err := t.ExecuteTemplate(buf, name, data)
	if err != nil {
		a.ServerError(w, r, err)
		return
	}

	w.WriteHeader(status)

	_, err = buf.WriteTo(w)
	if err != nil {
		a.ServerError(w, r, err)
		return
	}
}

// GetTemplateData generates a TemplateData struct populated with a CSRF token for the provided HTTP request.
func (a *App) GetTemplateData(r *http.Request) TemplateData {
	return TemplateData{
		AuthenticatedUserId: a.GetAuthenticatedUserId(r),
		IsAuthenticated:     a.IsAuthenticated(r),
		CurrentYear:         time.Now().Year(),
		CSRFToken:           nosurf.Token(r),
		Username:            a.GetAuthenticatedUserName(r),
		Flash:               a.SessionManager.PopString(r.Context(), "flash"),
	}
}

// SetFlashMessage sets a flash message in the session for the given HTTP request using the "flash" key.
func (a *App) SetFlashMessage(r *http.Request, msg string) {
	a.SessionManager.Put(r.Context(), "flash", msg)
}

// ServerError logs the provided server-side error, including the request method, URL, and stack trace, then sends a 500 response.
func (a *App) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		url    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	a.Logger.Error(err.Error(), "method", method, "url", url, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// ClientError logs a client-side error with method, URL, and status, then sends an HTTP response with the provided status.
func (a *App) ClientError(w http.ResponseWriter, r *http.Request, status int, err error) {
	var (
		method = r.Method
		url    = r.URL.RequestURI()
	)
	a.Logger.Error(err.Error(), "method", method, "url", url, "status", status)
	http.Error(w, http.StatusText(status), status)
}

// IsHtmxRequest checks if the incoming HTTP request is an HTMX request by inspecting the "HX-Request" header.
func (a *App) IsHtmxRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

// HtmxRedirect sets the "HX-Redirect" header to the specified URL and sends an HTTP 302 Found status response.
func (a *App) HtmxRedirect(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Redirect", url)
	w.WriteHeader(http.StatusFound)
}

// HtmxLocation sets the "HX-Location" response header with JSON data containing path, target, and swap values.
func (a *App) HtmxLocation(w http.ResponseWriter, r *http.Request, url string, target string, swap string) {
	locationObj := map[string]string{
		"path":   url,
		"target": target,
		"swap":   swap,
	}

	locationJSON, err := json.Marshal(locationObj)
	if err != nil {
		a.ServerError(w, r, err)
		return
	}
	w.Header().Set("HX-Location", string(locationJSON))
	w.WriteHeader(http.StatusOK)
}
