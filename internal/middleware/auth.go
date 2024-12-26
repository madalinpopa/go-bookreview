package middleware

import (
	"context"
	"github.com/madalinpopa/go-bookreview/internal/app"
	"github.com/madalinpopa/go-bookreview/internal/models"
	"net/http"
)

// LoginRequired is middleware that enforces user authentication by redirecting unauthenticated requests to the login page.
func (m *Middleware) LoginRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the user is not authenticated, redirect them to the login page and return
		// from the middleware chain so that no subsequent handlers in the chain are
		// executed.
		if !m.app.IsAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Otherwise set the "Cache-Control: no-store" header so that pages
		// require authentication are not stored in the users browser cache (or
		// other intermediary cache).
		w.Header().Add("Cache-Control", "no-store")

		// And call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

// RedirectAuthenticatedUsers is middleware that redirects authenticated users away from login or register pages to /events.
func (m *Middleware) RedirectAuthenticatedUsers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the user is authenticated using app's existing method
		if m.app.IsAuthenticated(r) {
			// Check if the path is either the login or register page
			if r.URL.Path == "/login" || r.URL.Path == "/register" {
				// Redirect authenticated users to the events page
				http.Redirect(w, r, "/events", http.StatusSeeOther)
				return
			}
		}
		// If not authenticated or not targeting the restricted pages, continue
		next.ServeHTTP(w, r)
	})
}

// Authenticate is middleware that ensures incoming requests are associated with a valid authenticated user.
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the authenticatedUserID value from the session using the GetInt()
		// method.
		//This will return the zero value for an int (0) if no
		// "authenticatedUserID" value is in the session -- in which case we call the
		// next handler in the chain as normal and return.
		id := m.app.SessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// Otherwise, we check to see if a user with that ID exists in our
		// database.
		exists, err := m.app.Models.Users.Exists(models.ID, id)
		if err != nil {
			m.app.ServerError(w, r, err)
			return
		}

		// If a matching user is found, we know that the request is coming from an
		// authenticated user who exists in our database.
		//We create a new copy of the
		// request (with an isAuthenticatedContextKey value of true in the request data)
		// and assign it to r.
		if exists {
			ctx := context.WithValue(r.Context(), app.IsAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}
