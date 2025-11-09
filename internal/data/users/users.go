// File: internal/data/domain/users/models.go
package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	internalErrors "github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/errors"
	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/filters"
	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/validator"
	"golang.org/x/crypto/bcrypt"
)

/************************************************************************************************************
 * User declarations
 ************************************************************************************************************/

// Password stores the password hash and kept plaintext for validation during writes.
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
	Password    Password  `json:"-"`
	IsActivated *bool     `json:"is_activated"`
	IsDeleted   *bool     `json:"is_deleted"`
	IsVerified  *bool     `json:"is_verified"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AnonymousUser is a sentinel anonymous user instance.
var AnonymousUser = &User{}

// User Filters struct for filtering user queries
type UserFilters struct {
	FarmerID    string
	Email       string
	PhoneNumber string
	Name        string
	IsActivated *bool
	IsDeleted   *bool
	IsVerified  *bool
	Filters     filters.Filters
}

// UserModels struct for database operations
type UserModel struct {
	DB *sql.DB
}

/************************************************************************************************************
 * Password helpers
 ************************************************************************************************************/

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

/************************************************************************************************************
 * User Validation
 ************************************************************************************************************/

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
 * Database Operations
 ************************************************************************************************************/

func (m *UserModel) Insert(user *User) error {
	// Query
	query := `
		INSERT INTO users (farmer_id, email, first_name, last_name, password_hash, is_activated, is_deleted, is_verified, phone_number)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at, version`

	// Arguments for Query
	args := []any{
		user.FarmerID,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Password.hash,
		user.IsActivated,
		user.IsDeleted,
		user.IsVerified,
		user.PhoneNumber,
	}

	// Get Context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Version)
	if err != nil {
		switch {
		case internalErrors.IsUniqueViolation(err, "email"):
			return internalErrors.ErrDuplicateValue("email")
		case internalErrors.IsUniqueViolation(err, "farmer_id"):
			return internalErrors.ErrDuplicateValue("farmer_id")
		case internalErrors.IsUniqueViolation(err, "phone_number"):
			return internalErrors.ErrDuplicateValue("phone_number")
		default:
			return internalErrors.WrapInsertError(err, "Users")
		}
	}
	return nil
}

/*
***************************************************************************************

	Update Methods

***************************************************************************************
*/
func (m *UserModel) Update(user *User) error {
	// Query
	query := `
		UPDATE users
		SET farmer_id = $1, email = $2, phone_number = $3, first_name = $4, last_name = $5, password_hash = $6, is_activated = $7, is_deleted = $8, is_verified = $9, updated_at = now(), version = version + 1
		WHERE id = $10 AND version = $11
		RETURNING updated_at, version  `

	// Arguments for Query
	args := []any{
		user.FarmerID,
		user.Email,
		user.PhoneNumber,
		user.FirstName,
		user.LastName,
		user.Password.hash,
		user.IsActivated,
		user.IsDeleted,
		user.IsVerified,
		user.ID,
		user.Version,
	}

	// Get Context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute Query
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.UpdatedAt, &user.Version)
	if err != nil {
		switch {
		case internalErrors.IsUniqueViolation(err, "email"):
			return internalErrors.ErrDuplicateValue("email")
		case internalErrors.IsUniqueViolation(err, "farmer_id"):
			return internalErrors.ErrDuplicateValue("farmer_id")
		case internalErrors.IsUniqueViolation(err, "phone_number"):
			return internalErrors.ErrDuplicateValue("phone_number")
		case internalErrors.IsEditConflict(err):
			return internalErrors.ErrEditConflict
		default:
			return internalErrors.WrapUpdateError(err, "Users")
		}
	}
	return nil
}

// Update Password Method
func (m *UserModel) UpdatePassword(user *User) error {
	// Query
	query := `
		UPDATE users
		SET password_hash = $1, updated_at = now(), version = version + 1
		WHERE id = $2 AND version = $3
		RETURNING updated_at, version`

	// Arguments for Query
	args := []any{
		user.Password.hash,
		user.ID,
		user.Version,
	}

	// Get Context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute Query
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.UpdatedAt, &user.Version)
	if err != nil {
		switch {
		case internalErrors.IsEditConflict(err):
			return internalErrors.ErrEditConflict
		default:
			return internalErrors.WrapUpdatePasswordError(err)
		}
	}
	return nil
}

/*
***************************************************************************************

	Delete Methods

***************************************************************************************
*/

// Soft Delete Method
func (m *UserModel) DeleteSoft(user *User) error {
	// Query
	query := `
		UPDATE users
		SET is_deleted = true, updated_at = now(), version = version + 1
		WHERE id = $1 AND is_deleted = false
		RETURNING updated_at, version`

	// Get Context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute Query
	err := m.DB.QueryRowContext(ctx, query, user.ID).Scan(&user.UpdatedAt, &user.Version)
	if err != nil {
		switch {
		case internalErrors.IsEditConflict(err):
			return internalErrors.ErrAlreadyDeleted
		default:
			return internalErrors.WrapDeleteError(err, "Users")
		}
	}
	return nil
}

// Hard Delete Method
func (m *UserModel) DeleteHard(userID int64) error {
	// Query
	query := `
		DELETE FROM users
		WHERE id = $1`

	// Get Context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute Query
	result, err := m.DB.ExecContext(ctx, query, userID)
	if err != nil {
		return internalErrors.WrapDeleteError(err, "Users")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return internalErrors.WrapDeleteError(err, "Users")
	}
	if rows == 0 {
		return internalErrors.ErrRecordNotFound
	}
	return nil
}

/*
***************************************************************************************

	Get Methods

***************************************************************************************
*/

// GetByID Method
func (m *UserModel) GetByID(userID int64) (*User, error) {
	// Query
	query := `
		SELECT id, farmer_id, email, phone_number, first_name, last_name, password_hash, is_activated, is_deleted, is_verified, version, created_at, updated_at
		FROM users
		WHERE id = $1`

	// Declare User variable
	var user User

	// Get Context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	scan := []any{
		&user.ID,
		&user.FarmerID,
		&user.Email,
		&user.PhoneNumber,
		&user.FirstName,
		&user.LastName,
		&user.Password.hash,
		&user.IsActivated,
		&user.IsDeleted,
		&user.IsVerified,
		&user.Version,
		&user.CreatedAt,
		&user.UpdatedAt,
	}

	// Execute Query
	err := m.DB.QueryRowContext(ctx, query, userID).Scan(scan...)
	if err != nil {
		switch {
		case internalErrors.IsRecordNotFound(err):
			return nil, internalErrors.ErrRecordNotFound
		default:
			return nil, internalErrors.WrapGetError(err, "Users")
		}
	}
	return &user, nil
}

// GetByEmail Method
func (m *UserModel) GetByEmail(email string) (*User, error) {
	// Query
	query := `
		SELECT id, farmer_id, email, phone_number, first_name, last_name, password_hash, is_activated, is_deleted, is_verified, version, created_at, updated_at
		FROM users
		WHERE email = $1`

	// Declare User variable
	var user User

	// Get Context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	scan := []any{
		&user.ID,
		&user.FarmerID,
		&user.Email,
		&user.PhoneNumber,
		&user.FirstName,
		&user.LastName,
		&user.Password.hash,
		&user.IsActivated,
		&user.IsDeleted,
		&user.IsVerified,
		&user.Version,
		&user.CreatedAt,
		&user.UpdatedAt,
	}

	// Execute Query
	err := m.DB.QueryRowContext(ctx, query, email).Scan(scan...)
	if err != nil {
		switch {
		case internalErrors.IsRecordNotFound(err):
			return nil, internalErrors.ErrRecordNotFound
		default:
			return nil, internalErrors.WrapGetError(err, "Users")
		}
	}
	return &user, nil
}

// GetByFarmerID Method
func (m *UserModel) GetByFarmerID(farmerID string) (*User, error) {
	// Query
	query := `
		SELECT id, farmer_id, email, phone_number, first_name, last_name, password_hash, is_activated, is_deleted, is_verified, version, created_at, updated_at
		FROM users
		WHERE farmer_id = $1`

	// Declare User variable
	var user User

	// Get Context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	scan := []any{
		&user.ID,
		&user.FarmerID,
		&user.Email,
		&user.PhoneNumber,
		&user.FirstName,
		&user.LastName,
		&user.Password.hash,
		&user.IsActivated,
		&user.IsDeleted,
		&user.IsVerified,
		&user.Version,
		&user.CreatedAt,
		&user.UpdatedAt,
	}

	// Execute Query
	err := m.DB.QueryRowContext(ctx, query, farmerID).Scan(scan...)
	if err != nil {
		switch {
		case internalErrors.IsRecordNotFound(err):
			return nil, internalErrors.ErrRecordNotFound
		default:
			return nil, internalErrors.WrapGetError(err, "Users")
		}
	}
	return &user, nil
}

// GetAll Method
func (m *UserModel) GetAll(u *UserFilters) ([]*User, filters.MetaData, error) {
	// Base Query
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, farmer_id, email, phone_number, first_name, last_name, password_hash, is_activated, is_deleted, is_verified, version, created_at, updated_at
		FROM users
		WHER E (to_tsvector('simple', farmer_id) @@ plainto_tsquery('simple', $1) OR $1 = '')
        AND (to_tsvector('simple', first_name || ' ' || last_name || ' ') @@ plainto_tsquery('simple', $2) OR $2 = '')
        AND (to_tsvector('simple', email) @@ plainto_tsquery('simple', $3) OR $3 = '')
		AND (to_tsvector('simple', phone_number) @@ plainto_tsquery('simple', $4) OR $4 = '')
        AND ($4::boolean IS NULL OR is_deleted = $4)
        AND ($5::boolean IS NULL OR is_activated = $5)
        AND ($6::boolean IS NULL OR is_verified = $6)
		ORDER BY %s %s, id ASC
		LIMIT $6 OFFSET $7`, u.Filters.SortColumn(), u.Filters.SortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		u.FarmerID,
		u.Name,
		u.Email,
		u.PhoneNumber,
		u.IsDeleted,
		u.IsActivated,
		u.IsVerified,
		u.Filters.Limit(),
		u.Filters.Offset(),
	}

	// declare MetaData variable
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, filters.MetaData{}, internalErrors.WrapGetAllError(err, "Users")
	}
	defer rows.Close()

	totalRecords := 0
	users := []*User{}
	for rows.Next() {
		var user User
		scan := []any{
			&totalRecords,
			&user.ID,
			&user.FarmerID,
			&user.Email,
			&user.PhoneNumber,
			&user.FirstName,
			&user.LastName,
			&user.Password.hash,
			&user.IsActivated,
			&user.IsDeleted,
			&user.IsVerified,
			&user.Version,
			&user.CreatedAt,
			&user.UpdatedAt,
		}
		err := rows.Scan(scan...)
		if err != nil {
			return nil, filters.MetaData{}, internalErrors.WrapGetAllError(err, "Users")
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, filters.MetaData{}, err // Return any error encountered while iterating over the rows
	}
	meta := filters.CalculateMetaData(totalRecords, u.Filters.Page, u.Filters.PageSize)
	return users, meta, nil
}
