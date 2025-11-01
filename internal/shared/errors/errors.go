// File: internal/shared/errors/errors.go

package errors

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

// Predefined errors for common scenarios
var (
	ErrRecordNotFound      = errors.New("record not found")
	ErrEditConflict        = errors.New("edit conflict detected")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrAlreadyDeleted      = errors.New("user already deleted")
	ErrNoMatch             = errors.New("no matching records found")
	ErrForeignKeyViolation = errors.New("constraint violation")
	ErrInvalidUpdateData   = errors.New("invalid update data")
	ErrPermissionDenied    = errors.New("permission denied for update")
	ErrStaleData           = errors.New("stale data: record has been modified elsewhere")
	ErrLockedResource      = errors.New("resource is locked, cannot update")
	ErrInsertFailed        = errors.New("insert failed: ")
	ErrUpdateFailed        = errors.New("update failed: ")
	ErrDeleteFailed        = errors.New("delete failed: ")
)

// isUniqueViolation checks where the error is a unique constraint violation
func IsUniqueViolation(err error, column string) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505" && strings.Contains(pqErr.Detail, "("+column+")")
}

// response for duplicate value errors
func ErrDuplicateValue(column string) error {
	return errors.New("duplicate value for column: " + column)
}

// for optimistic locking
func IsEditConflict(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

// for foreign key violations
func IsForeignKeyViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23503"
}

// WrapInsertError wraps an insert error with additional context
func WrapInsertError(err error, model string) error {
	return fmt.Errorf("%s insert failed: %w", model, err)
}

// WrapUpdateError wraps an update error with additional context
func WrapUpdateError(err error, model string) error {
	return fmt.Errorf("%s update failed: %w", model, err)
}

// wrapUpdatePasswordError wraps a password update error with additional context
func WrapUpdatePasswordError(err error) error {
	return fmt.Errorf("User password update failed: %w", err)
}

// WrapDeleteError wraps a delete error with additional context
func WrapDeleteError(err error, model string) error {
	return fmt.Errorf("%s delete failed: %w", model, err)
}

func WrapPasswordSetError(err error) error {
	return fmt.Errorf("password set failed: %w", err)
}

func WrapGetError(err error, model string) error {
	return fmt.Errorf("%s get failed: %w", model, err)
}

func WrapGetAllError(err error, model string) error {
	return fmt.Errorf("%s get all failed: %w", model, err)
}

func IsRecordNotFound(err error) bool {
	return sql.ErrNoRows == err
}

func ValidationFailed(details map[string]string) error {
	return fmt.Errorf("validation failed: %v", details)
}
