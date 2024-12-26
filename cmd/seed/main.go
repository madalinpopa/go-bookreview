package main

import (
	"github.com/madalinpopa/go-bookreview/internal/app"
	"github.com/madalinpopa/go-bookreview/internal/models"
	"log/slog"
	"os"
)

// Logger is a global variable that holds a pointer to an instance of slog.Logger for logging application messages.
var Logger *slog.Logger

// seedAdminUser seeds the database with an admin user if it does not already exist.
// Returns an error if the existence check or user creation fails.
func seedAdminUser(m *models.Models) error {
	exist, err := m.Users.Exists(models.Username, "admin")
	if err != nil {
		return err
	}
	Logger.Info("Checking if admin user exists", "exist", exist)
	if !exist {
		Logger.Info("Creating admin user")
		err = m.Users.Create("admin", "admin@test.com", "secret123")
		if err != nil {
			return err
		}
	} else {
		Logger.Info("Admin user already exists")
	}

	return nil
}

// main is the entry point of the application; it initializes configuration, sets up the database, and seeds initial data.
func main() {
	Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	config := app.NewConfig()
	db, err := app.CreateDatabaseConnection(config.Dsn)
	if err != nil {
		Logger.Error(err.Error())
		os.Exit(1)
	}
	Logger.Info("Database connection established")
	m := models.NewModels(db, Logger)

	err = seedAdminUser(m)
	if err != nil {
		Logger.Error(err.Error())
	}

}
