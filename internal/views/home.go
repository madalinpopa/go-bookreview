package views

import (
	"github.com/madalinpopa/go-bookreview/internal/app"
	"net/http"
)

// HomePage returns an HTTP handler that renders the home page template with data provided by the App instance.
func HomePage(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := app.GetTemplateData(r)
		app.Render(w, r, "index.tmpl", data, http.StatusOK)
	}
}
