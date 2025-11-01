// File: internal/data/domain/users/services.go

package users

import (
	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/errors"
	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/filters"
	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/validator"
)

type UserService struct {
	repo UserRepository
}

func NewUserUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser creates a new user with validation
func (s *UserService) CreateUser(user *User) error {
	v := validator.New()
	ValidateUser(v, user)
	if !v.Valid() {
		return errors.ValidationFailed(v.Errors)
	}
	return s.repo.Insert(user)
}

// UpdateUser updates an existing user with validation
func (s *UserService) UpdateUser(user *User) error {
	v := validator.New()
	ValidateUser(v, user)
	if !v.Valid() {
		return errors.ValidationFailed(v.Errors)
	}
	return s.repo.Update(user)
}

// ChangeUserPassword changes a user's password with validation
func (s *UserService) ChangeUserPassword(user *User, newPassword string) error {
	v := validator.New()
	ValidatePasswordPlaintext(v, newPassword)
	if !v.Valid() {
		return errors.ValidationFailed(v.Errors)
	}
	err := user.Password.Set(newPassword)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(user)
}

// DeleteUserSoft performs a soft delete on a user
func (s *UserService) DeleteUserSoft(user *User) error {
	return s.repo.DeleteSoft(user)
}

// DeleteUserHard performs a hard delete on a user by ID
func (s *UserService) DeleteUserHard(userID int64) error {
	return s.repo.DeleteHard(userID)
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(userID int64) (*User, error) {
	return s.repo.GetByID(userID)
}

// GetUserByEmail retrieves a user by their email
func (s *UserService) GetUserByEmail(email string) (*User, error) {
	return s.repo.GetByEmail(email)
}

// GetUserByFarmerID retrieves a user by their farmer ID
func (s *UserService) GetUserByFarmerID(farmerID string) (*User, error) {
	return s.repo.GetByFarmerID(farmerID)
}

// GetAllUsers retrieves all users with filtering and pagination
func (s *UserService) GetAllUsers(u UserFilters) ([]*User, filters.MetaData, error) {
	return s.repo.GetAll(u)
}
