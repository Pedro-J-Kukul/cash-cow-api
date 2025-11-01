// File: internal/data/domain/users/services.go

package users

import (
	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/errors"
	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/filters"
	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/validator"
)

// Service handles user business logic
type Service struct {
	repo *UserRepository
}

// NewService creates a new user service
func NewService(repo *UserRepository) *Service {
	return &Service{repo: repo}
}

// CreateUser creates a new user with validation
func (s *Service) CreateUser(user *User) error {
	v := validator.New()
	ValidateUser(v, user)
	if !v.Valid() {
		return errors.ValidationFailed(v.Errors)
	}
	return (*s.repo).Insert(user)
}

// UpdateUser updates an existing user with validation
func (s *Service) UpdateUser(user *User) error {
	v := validator.New()
	ValidateUser(v, user)
	if !v.Valid() {
		return errors.ValidationFailed(v.Errors)
	}
	return (*s.repo).Update(user)
}

// ChangeUserPassword changes a user's password with validation
func (s *Service) ChangeUserPassword(user *User, newPassword string) error {
	v := validator.New()
	ValidatePasswordPlaintext(v, newPassword)
	if !v.Valid() {
		return errors.ValidationFailed(v.Errors)
	}
	err := user.Password.Set(newPassword)
	if err != nil {
		return err
	}
	return (*s.repo).UpdatePassword(user)
}

// DeleteUserSoft performs a soft delete on a user
func (s *Service) DeleteUserSoft(user *User) error {
	return (*s.repo).DeleteSoft(user)
}

// DeleteUserHard performs a hard delete on a user by ID
func (s *Service) DeleteUserHard(userID int64) error {
	return (*s.repo).DeleteHard(userID)
}

// GetUserByID retrieves a user by their ID
func (s *Service) GetUserByID(userID int64) (*User, error) {
	return (*s.repo).GetByID(userID)
}

// GetUserByEmail retrieves a user by their email
func (s *Service) GetUserByEmail(email string) (*User, error) {
	return (*s.repo).GetByEmail(email)
}

// GetUserByFarmerID retrieves a user by their farmer ID
func (s *Service) GetUserByFarmerID(farmerID string) (*User, error) {
	return (*s.repo).GetByFarmerID(farmerID)
}

// GetAllUsers retrieves all users with filtering and pagination
func (s *Service) GetAllUsers(u UserFilters) ([]*User, filters.MetaData, error) {
	return (*s.repo).GetAll(u)
}
