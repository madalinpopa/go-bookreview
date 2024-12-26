package main

import (
	"fmt"
	"github.com/madalinpopa/go-bookreview/internal/app"
	"github.com/madalinpopa/go-bookreview/internal/routes"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// main initializes and starts the HTTP app,
// logging the app address and handling potential errors during execution.
func main() {

	// Initialize the application configuration and create a new app instance.
	config := app.NewConfig()

	// Open database connection using the DSN from the configuration.
	db, err := app.CreateDatabaseConnection(config.Dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Run database migrations
	if err := app.MakeMigrations(db); err != nil {
		log.Fatal(err)
	}

	// Create a new app instance using the configuration.
	a := app.NewApp(config, db)

	// Create file uploads directory
	if err := os.MkdirAll(config.UploadDir, 0755); err != nil {
		a.Logger.Error(err.Error())
		os.Exit(1)
	}

	// Load templates into the app to set up rendering for HTML views.
	if err := a.LoadTemplates(); err != nil {
		a.Logger.Error(err.Error())
		os.Exit(1)
	}

	// Create and configure the HTTP server.
	s := http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Addr, config.Port),
		Handler:      routes.UrlPatterns(a),
		ErrorLog:     slog.NewLogLogger(a.Logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
	}

	a.Logger.Info("Starting app", "addr", s.Addr)

	if err := s.ListenAndServe(); err != nil {
		a.Logger.Error(err.Error())
		os.Exit(1)
	}
}
