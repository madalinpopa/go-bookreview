package middleware

import (
	"fmt"
	"github.com/madalinpopa/go-bookreview/internal/app"
	"net/http"
	"strings"
)

// responseWriter is a custom wrapper for http.ResponseWriter to capture HTTP status codes within handlers.
type responseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader sets the HTTP status code and forwards it to the embedded ResponseWriter.
func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

type Middleware struct {
	app *app.App
}

// NewMiddleware initializes and returns a new Middleware instance, associating it with the provided app.App instance.
func NewMiddleware(app *app.App) *Middleware {
	return &Middleware{app: app}
}

// Headers is an HTTP middleware that sets various security-related headers for incoming requests.
func (m *Middleware) Headers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set Content-Security-Policy to restrict the sources of content such as scripts, styles, and images
		w.Header().Set("Content-Security-Policy", strings.TrimSpace(`
            default-src 'self';
            script-src 'self' 'unsafe-inline' 'unsafe-eval' cdn.jsdelivr.net *.iconify.design *.simplesvg.com *.unisvg.com;
            style-src 'self' 'unsafe-inline' cdn.jsdelivr.net;
            img-src 'self' data: cdn.jsdelivr.net *.iconify.design *.simplesvg.com *.unisvg.com;
            font-src 'self' cdn.jsdelivr.net;
            connect-src 'self' cdn.jsdelivr.net *.iconify.design *.simplesvg.com *.unisvg.com;
        `))

		// Set Referrer-Policy to control the amount of referrer information sent with requests
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Set X-Content-Type-Options to prevent browsers from interpreting files as a different MIME type
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Set X-Frame-Options to prevent clickjacking attacks by disallowing the page from being framed
		w.Header().Set("X-Frame-Options", "deny")

		// Set X-XSS-Protection to disable the browser's XSS protection, preventing unintended behavior
		w.Header().Set("X-XSS-Protection", "0")

		// Set Server header to specify the server software being used (can also be customized or omitted)
		w.Header().Set("Server", "Go")

		next.ServeHTTP(w, r)
	})
}

// Logging is an HTTP middleware that logs details of incoming requests, including IP, protocol, method, and URL.
func (m *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			url    = r.URL.RequestURI()
			status = rw.status
		)

		m.app.Logger.Info("request", "ip", ip, "proto", proto, "method", method, "url", url, "status", status)
	})
}

// Recover is a middleware that intercepts panics, logs the error via ServerError, and responds with a 500 status.
func (m *Middleware) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				m.app.ServerError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
