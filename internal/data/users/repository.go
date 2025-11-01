// internal/data/users/repository.go
package users

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/errors"
	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/filters"
)

type UserRepository interface {
	Insert(user *User) error
	Update(user *User) error
	UpdatePassword(user *User) error
	DeleteSoft(user *User) error
	DeleteHard(userID int64) error
	GetByID(userID int64) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByFarmerID(farmerID string) (*User, error)
	GetAll(u UserFilters) ([]*User, filters.MetaData, error)
}

type userRepository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) UserRepository {
	return &userRepository{DB: db}

}

/*
***************************************************************************************

	Insert Method

***************************************************************************************
*/
func (r *userRepository) Insert(user *User) error {
	// Query
	query := `
		INSERT INTO users (farmer_id, email, first_name, last_name, middle_name, password_hash, is_activated, is_deleted, is_verified, phone_number)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at, version`

	// Arguments for Query
	args := []any{
		user.FarmerID,
		user.Email,
		user.FirstName,
		user.LastName,
		user.MiddleName,
		user.Password.hash,
		user.IsActivated,
		user.IsDeleted,
		user.IsVerified,
		user.PhoneNumber,
	}

	// Get Context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Version)
	if err != nil {
		switch {
		case errors.IsUniqueViolation(err, "email"):
			return errors.ErrDuplicateValue("email")
		case errors.IsUniqueViolation(err, "farmer_id"):
			return errors.ErrDuplicateValue("farmer_id")
		case errors.IsUniqueViolation(err, "phone_number"):
			return errors.ErrDuplicateValue("phone_number")
		default:
			return errors.WrapInsertError(err, "Users")
		}
	}
	return nil
}

/*
***************************************************************************************

	Update Methods

***************************************************************************************
*/
func (r *userRepository) Update(user *User) error {
	// Query
	query := `
		UPDATE users
		SET farmer_id = $1, email = $2, phone_number = $3, first_name = $4, last_name = $5, middle_name = $6, password_hash = $7, is_activated = $8, is_deleted = $9, is_verified = $10, updated_at = now(), version = version + 1
		WHERE id = $11 AND version = $12
		RETURNING updated_at, version  `

	// Arguments for Query
	args := []any{
		user.FarmerID,
		user.Email,
		user.PhoneNumber,
		user.FirstName,
		user.LastName,
		user.MiddleName,
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
	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&user.UpdatedAt, &user.Version)
	if err != nil {
		switch {
		case errors.IsUniqueViolation(err, "email"):
			return errors.ErrDuplicateValue("email")
		case errors.IsUniqueViolation(err, "farmer_id"):
			return errors.ErrDuplicateValue("farmer_id")
		case errors.IsUniqueViolation(err, "phone_number"):
			return errors.ErrDuplicateValue("phone_number")
		case errors.IsEditConflict(err):
			return errors.ErrEditConflict
		default:
			return errors.WrapUpdateError(err, "Users")
		}
	}
	return nil
}

// Update Password Method
func (r *userRepository) UpdatePassword(user *User) error {
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
	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&user.UpdatedAt, &user.Version)
	if err != nil {
		switch {
		case errors.IsEditConflict(err):
			return errors.ErrEditConflict
		default:
			return errors.WrapUpdatePasswordError(err)
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
func (r *userRepository) DeleteSoft(user *User) error {
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
	err := r.DB.QueryRowContext(ctx, query, user.ID).Scan(&user.UpdatedAt, &user.Version)
	if err != nil {
		switch {
		case errors.IsEditConflict(err):
			return errors.ErrAlreadyDeleted
		default:
			return errors.WrapDeleteError(err, "Users")
		}
	}
	return nil
}

// Hard Delete Method
func (r *userRepository) DeleteHard(userID int64) error {
	// Query
	query := `
		DELETE FROM users
		WHERE id = $1`

	// Get Context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute Query
	result, err := r.DB.ExecContext(ctx, query, userID)
	if err != nil {
		return errors.WrapDeleteError(err, "Users")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return errors.WrapDeleteError(err, "Users")
	}
	if rows == 0 {
		return errors.ErrRecordNotFound
	}
	return nil
}

/*
***************************************************************************************

	Get Methods

***************************************************************************************
*/

// GetByID Method
func (r *userRepository) GetByID(userID int64) (*User, error) {
	// Query
	query := `
		SELECT id, farmer_id, email, phone_number, first_name, last_name, middle_name, password_hash, is_activated, is_deleted, is_verified, version, created_at, updated_at
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
		&user.MiddleName,
		&user.Password.hash,
		&user.IsActivated,
		&user.IsDeleted,
		&user.IsVerified,
		&user.Version,
		&user.CreatedAt,
		&user.UpdatedAt,
	}

	// Execute Query
	err := r.DB.QueryRowContext(ctx, query, userID).Scan(scan...)
	if err != nil {
		switch {
		case errors.IsRecordNotFound(err):
			return nil, errors.ErrRecordNotFound
		default:
			return nil, errors.WrapGetError(err, "Users")
		}
	}
	return &user, nil
}

// GetByEmail Method
func (r *userRepository) GetByEmail(email string) (*User, error) {
	// Query
	query := `
		SELECT id, farmer_id, email, phone_number, first_name, last_name, middle_name, password_hash, is_activated, is_deleted, is_verified, version, created_at, updated_at
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
		&user.MiddleName,
		&user.Password.hash,
		&user.IsActivated,
		&user.IsDeleted,
		&user.IsVerified,
		&user.Version,
		&user.CreatedAt,
		&user.UpdatedAt,
	}

	// Execute Query
	err := r.DB.QueryRowContext(ctx, query, email).Scan(scan...)
	if err != nil {
		switch {
		case errors.IsRecordNotFound(err):
			return nil, errors.ErrRecordNotFound
		default:
			return nil, errors.WrapGetError(err, "Users")
		}
	}
	return &user, nil
}

// GetByFarmerID Method
func (r *userRepository) GetByFarmerID(farmerID string) (*User, error) {
	// Query
	query := `
		SELECT id, farmer_id, email, phone_number, first_name, last_name, middle_name, password_hash, is_activated, is_deleted, is_verified, version, created_at, updated_at
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
		&user.MiddleName,
		&user.Password.hash,
		&user.IsActivated,
		&user.IsDeleted,
		&user.IsVerified,
		&user.Version,
		&user.CreatedAt,
		&user.UpdatedAt,
	}

	// Execute Query
	err := r.DB.QueryRowContext(ctx, query, farmerID).Scan(scan...)
	if err != nil {
		switch {
		case errors.IsRecordNotFound(err):
			return nil, errors.ErrRecordNotFound
		default:
			return nil, errors.WrapGetError(err, "Users")
		}
	}
	return &user, nil
}

// GetAll Method
func (r *userRepository) GetAll(u UserFilters) ([]*User, filters.MetaData, error) {
	// Base Query
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, farmer_id, email, phone_number, first_name, last_name, middle_name, password_hash, is_activated, is_deleted, is_verified, version, created_at, updated_at
		FROM users
		WHERE (to_tsvector('simple', farmer_id) @@ plainto_tsquery('simple', $1) OR $1 = '')
        AND (to_tsvector('simple', first_name || ' ' || last_name || ' ' || coalesce(middle_name, '')) @@ plainto_tsquery('simple', $2) OR $2 = '')
        AND (to_tsvector('simple', email) @@ plainto_tsquery('simple', $3) OR $3 = '')
		ANDD (to_tsvector('simple', phone_number) @@ plainto_tsquery('simple', $4) OR $4 = '')
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
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, filters.MetaData{}, errors.WrapGetAllError(err, "Users")
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
			&user.MiddleName,
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
			return nil, filters.MetaData{}, errors.WrapGetAllError(err, "Users")
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, filters.MetaData{}, err // Return any error encountered while iterating over the rows
	}
	meta := filters.CalculateMetaData(totalRecords, u.Filters.Page, u.Filters.PageSize)
	return users, meta, nil
}
