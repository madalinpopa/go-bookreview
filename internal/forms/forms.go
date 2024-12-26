package forms

import (
	"errors"
	"fmt"
	"github.com/madalinpopa/go-bookreview/internal/app"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"time"
)

// UserLoginForm represents the structure for capturing user login data submitted via a form.
type UserLoginForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
	Base     `form:"-"`
}

// Validate validates the UserLoginForm by ensuring both the username and password fields are not blank.
func (ul *UserLoginForm) Validate() {
	ul.CheckField(NotBlank(ul.Username), "username", "This field is required")
	ul.CheckField(NotBlank(ul.Password), "password", "This field is required")
}

// RegisterForm represents the structure for user registration form input data.
type RegisterForm struct {
	Email    string `form:"email"`
	Username string `form:"username"`
	Password string `form:"password"`
	Base     `form:"-"`
}

// Validate performs validation on the RegisterForm fields and adds errors for invalid or missing input.
func (r *RegisterForm) Validate() {
	r.CheckField(NotBlank(r.Username), "username", "This field is required")
	r.CheckField(NotBlank(r.Email), "email", "This field is required")
	r.CheckField(Matches(r.Email, EmailRX), "email", "The email address is not valid.")
	r.CheckField(NotBlank(r.Password), "password", "This field is required")
	r.CheckField(MinChars(r.Password, 8), "password", "Password must be at least 8 characters.")
}

// BookForm represents the form data for creating a book with fields for title, author, ISBN, and publication year.
// Includes embedded Base for validation-related functionalities.
type BookForm struct {
	Id              int    `form:"id"`
	Title           string `form:"title"`
	Author          string `form:"author"`
	ISBN            string `form:"isbn"`
	PublicationYear int    `form:"publication_year"`
	Status          string `form:"status"`
	ImageURL        string `form:"-"`
	CurrentImageURL string `form:"-"`
	Base            `form:"-"`
}

// Validate checks the BookForm fields for compliance with required rules and adds errors for blank fields.
func (cb *BookForm) Validate() {
	cb.CheckField(NotBlank(cb.Title), "title", "Title is required")
	cb.CheckField(NotBlank(cb.Author), "author", "Author is required")
	cb.CheckField(NotBlank(cb.ISBN), "isbn", "ISBN is required")
}

// HandleFileUpload processes an uploaded image file, validates its type, and saves it to the server's upload directory.
// Returns an error for invalid file types, missing files, or file handling issues.
func (cb *BookForm) HandleFileUpload(app *app.App, r *http.Request) error {
	existingImageUrl := r.FormValue("current_image_url")

	// Check if any error and if the file is not provided.
	file, fileHeader, err := r.FormFile("image_upload")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			cb.ImageURL = existingImageUrl
			return nil
		}
		app.Logger.Error(err.Error())
		return fmt.Errorf("%w: %v", ErrFormBadRequest, err)
	}
	defer func() {
		if file != nil {
			if err := file.Close(); err != nil {
				app.Logger.Error(err.Error())
			}
		}
	}()

	if file != nil && fileHeader != nil {
		acceptedTypes := []string{"image/jpeg", "image/png"}
		fileType := fileHeader.Header.Get("Content-Type")
		if !slices.Contains(acceptedTypes, fileType) {
			cb.AddNonFieldError("Invalid file type")
			return ErrInvalidFileType
		}

		filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), fileHeader.Filename)
		imageUrl := fmt.Sprintf("/uploads/%s", filename)
		fullPath := filepath.Join(app.Config.UploadDir, filename)

		dst, err := os.Create(fullPath)
		if err != nil {
			return err
		}
		defer func() {
			if err := dst.Close(); err != nil {
				app.Logger.Error(err.Error())
			}
		}()

		if _, err := io.Copy(dst, file); err != nil {
			return err
		}

		cb.ImageURL = imageUrl

		// delete old image
		if existingImageUrl != "" {
			oldPath := filepath.Join(app.Config.UploadDir, filepath.Base(existingImageUrl))
			if err := os.Remove(oldPath); err != nil {
				app.Logger.Error("Failed to remove old image: " + err.Error())
			}
		}
	}

	return nil
}

// BookReviewForm represents a form structure for submitting a book review with a rating and review text.
type BookReviewForm struct {
	Id         int    `form:"id"`
	BookId     int    `form:"book_id"`
	Rating     int    `form:"rating"`
	ReviewText string `form:"review_text"`
	Base       `form:"-"`
}

// Validate performs validation on the BookReviewForm fields to ensure required fields are not blank.
func (br *BookReviewForm) Validate() {
	br.CheckField(MaxNumber(br.Rating, 5), "rating", "Rating must be a number between 1 and 5")
	br.CheckField(NotBlank(br.ReviewText), "review_text", "Review text is required")
}

// BookNoteForm represents the form for submitting a note related to a book, including validation and associated details.
type BookNoteForm struct {
	Id         int    `form:"id"`
	BookId     int    `form:"book_id"`
	NoteText   string `form:"note_text"`
	PageNumber int    `form:"page_number"`
	Base       `form:"-"`
}

// Validate checks if the NoteText field is not blank and adds a validation error if the field is empty or whitespace-only.
func (bn *BookNoteForm) Validate() {
	bn.CheckField(NotBlank(bn.NoteText), "note_text", "Note text is required")
}
