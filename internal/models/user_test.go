package models

import (
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/madalinpopa/go-bookreview/internal/testutil"
)

// TestUserModel_Create tests the Create method of UserModel for various scenarios, including valid and invalid user creation.
func TestUserModel_Create(t *testing.T) {
	db := testutil.NewTestDB(t)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	model := UserModel{DB: db, Logger: logger}

	tests := []struct {
		name     string
		username string
		email    string
		password string
		wantErr  error
	}{
		{
			name:     "valid user",
			username: "testuser",
			email:    "test@example.com",
			password: "password123",
			wantErr:  nil,
		},
		{
			name:     "duplicate username",
			username: "testuser",
			email:    "different@example.com",
			password: "password123",
			wantErr:  ErrDuplicateUsername,
		},
		{
			name:     "duplicate email",
			username: "differentuser",
			email:    "test@example.com",
			password: "password123",
			wantErr:  ErrDuplicateEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := model.Create(tt.username, tt.email, tt.password)
			if !errors.Is(err, tt.wantErr) {
				if tt.wantErr == nil {
					t.Errorf("unexpected error: %v", err)
				} else if !errors.Is(err, tt.wantErr) {
					t.Errorf("got error %v; want %v", err, tt.wantErr)
				}
			}
		})
	}
}
