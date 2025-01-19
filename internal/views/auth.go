package views

import (
	"errors"
	"net/http"

	"github.com/madalinpopa/go-bookreview/internal/app"
	"github.com/madalinpopa/go-bookreview/internal/forms"
	"github.com/madalinpopa/go-bookreview/internal/models"
)

// LoginPage handles the HTTP request for rendering the user login page with a preloaded form and template data.
func LoginPage(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form forms.UserLoginForm

		data := app.GetTemplateData(r)
		data.Form = form
		app.Render(w, r, "login.tmpl", data, http.StatusOK)
	}
}

// LoginPost handles HTTP POST requests for user login authentication.
// It validates submitted form data, verifies user credentials, and initiates a session upon successful login.
func LoginPost(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form forms.UserLoginForm
		data := app.GetTemplateData(r)

		err := app.FormDecoder.Decode(&form, r.PostForm)
		if err != nil {
			app.ClientError(w, r, http.StatusBadRequest, err)
			return
		}

		form.Validate()
		if !form.Valid() {
			data.Form = form
			app.Render(w, r, "login.tmpl", data, http.StatusUnprocessableEntity)
			return
		}

		id, err := app.Models.Users.Authenticate(form.Username, form.Password)
		if err != nil {
			if errors.Is(err, models.ErrInvalidCredentials) {
				form.AddNonFieldError("Invalid email or password.")
				data.Form = form
				app.Render(w, r, "login.tmpl", data, http.StatusUnprocessableEntity)
			} else {
				app.ServerError(w, r, err)
			}
			return
		}

		err = app.SessionManager.RenewToken(r.Context())
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		app.SessionManager.Put(r.Context(), "authenticatedUserID", id)
		app.SessionManager.Put(r.Context(), "authenticatedUsername", form.Username)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// LogoutPost handles POST requests for user logout,
// renews the session token, and removes authentication session data.
func LogoutPost(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Renew session token to prevent session fixation attacks
		err := app.SessionManager.RenewToken(r.Context())
		if err != nil {
			app.ServerError(w, r, err)
			return
		}
		// Remove authenticated user id
		app.SessionManager.Remove(r.Context(), "authenticatedUserID")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// RegisterPage handles the registration page rendering by executing the "auth/register.tmpl"
// template with necessary data.
func RegisterPage(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form forms.RegisterForm
		data := app.GetTemplateData(r)
		data.Form = form
		app.Render(w, r, "register.tmpl", data, http.StatusOK)
	}
}

// RegisterPost handles the HTTP POST request for user registration by validating form input and rendering templates.
// Displays validation errors or continues with further processing if input is valid.
// Returns an HTTP response with appropriate status codes based on the validation outcome.
func RegisterPost(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := app.GetTemplateData(r)
		var form forms.RegisterForm

		err := app.FormDecoder.Decode(&form, r.PostForm)
		if err != nil {
			app.ClientError(w, r, http.StatusBadRequest, err)
			return
		}

		form.Validate()
		if !form.Valid() {
			data.Form = form
			app.Render(w, r, "register.tmpl", data, http.StatusUnprocessableEntity)
			return
		}

		err = app.Models.Users.Create(form.Username, form.Email, form.Password)
		if err != nil {

			if errors.Is(err, models.ErrDuplicateEmail) {
				form.AddFieldError("email", "This email address is already registered.")
				data.Form = form
				app.Render(w, r, "register.tmpl", data, http.StatusUnprocessableEntity)
			} else if errors.Is(err, models.ErrDuplicateUsername) {
				form.AddFieldError("username", "This username is already registered.")
				data.Form = form
				app.Render(w, r, "register.tmpl", data, http.StatusUnprocessableEntity)
			} else {
				app.ServerError(w, r, err)
			}
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
