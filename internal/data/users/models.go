// File: internal/data/domain/users/models.go
package users

import (
	"errors"
	"time"

	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/filters"
	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/validator"
	"golang.org/x/crypto/bcrypt"
)

/************************************************************************************************************/
// User declarations
/************************************************************************************************************/
type UserModels interface {
	ValidateUser(v *validator.Validator, user *User)
	ValidateEmail(v *validator.Validator, email string)
	ValidatePasswordPlaintext(v *validator.Validator, password string)
}

// Password stores the hashed password and optional plaintext (used for validation during write operations).
type Password struct {
	hash      []byte
	plaintext *string
}

// User represents an application user.
type User struct {
	ID          int64     `json:"id"`
	FarmerID    string    `json:"farmer_id"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	MiddleName  string    `json:"middle_name"`
	Password    Password  `json:"-"`
	IsActivated bool      `json:"is_activated"`
	IsDeleted   bool      `json:"is_deleted"`
	IsVerified  bool      `json:"is_verified"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AnonymousUser is a sentinel anonymous user instance.
var AnonymousUser = &User{}

/************************************************************************************************************/
// Password helpers
/************************************************************************************************************/

// Set hashes the supplied plaintext password.
func (p *Password) Set(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), 12)
	if err != nil {
		return err
	}
	p.hash = hash
	p.plaintext = &plaintext
	return nil
}

// Matches verifies that the supplied plaintext password matches the stored hash.
func (p *Password) Matches(plaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintext))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

/************************************************************************************************************/
// User Validation
/************************************************************************************************************/

// IsAnonymous checks if the user is anonymous
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

// ValidateEmail checks if the email is valid
func validateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(len(email) <= 254, "email", "must not be more than 254 bytes long")
	v.Check(v.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

// ValidatePasswordPlaintext checks if the plaintext password is valid
func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 characters long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 characters long")
	v.Check(v.Matches(password, validator.PasswordNumberRX), "password", "must contain at least one number")
	v.Check(v.Matches(password, validator.PasswordUpperRX), "password", "must contain at least one uppercase letter")
	v.Check(v.Matches(password, validator.PasswordLowerRX), "password", "must contain at least one lowercase letter")
	v.Check(v.Matches(password, validator.PasswordSpecialRX), "password", "must contain at least one special character")
}

// ValidateUser checks if the user struct is valid
func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.FirstName != "", "first_name", "must be provided")
	v.Check(len(user.FirstName) <= 50, "first_name", "must not be more than 50 characters long")
	v.Check(user.LastName != "", "last_name", "must be provided")
	v.Check(len(user.LastName) <= 50, "last_name", "must not be more than 50 characters long")
	v.Check(len(user.MiddleName) <= 50, "middle_name", "must not be more than 50 characters long")
	v.Check(len(user.FarmerID) <= 50, "farmer_id", "must not be more than 50 characters long")
	v.Check(len(user.PhoneNumber) <= 15, "phone_number", "must not be more than 15 characters long")
	validateEmail(v, user.Email)
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

/************************************************************************************************************
 * Filters
 ************************************************************************************************************/
type UserFilters struct {
	FarmerID    string
	Email       string
	PhoneNumber string
	Name        string
	IsActivated *bool // nil = don't filter, otherwise filter by value
	IsDeleted   *bool // nil = don't filter, otherwise filter by value
	IsVerified  *bool // nil = don't filter, otherwise filter by value
	Filters     filters.Filters
}
