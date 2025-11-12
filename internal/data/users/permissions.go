package users

import (
	"context"
	"database/sql"
	"slices"
	"time"

	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/validator"
	"github.com/lib/pq"
)

/************************************************************************************************************
 * Permission declarations
 ************************************************************************************************************/

// Permission for action
type Permission struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

// User Permission Relationship
type UserPermission struct {
	UserID       int64 `json:"user_id"`
	PermissionID int64 `json:"permission_id"`
}

// Array of permissions
type Permissions []string

type PermissionModel struct {
	DB *sql.DB
}

// Helper to check if string is included in the list of permissions.
func (p Permissions) Includes(code string) bool {
	return slices.Contains(p, code)
}

/************************************************************************************************************
 * Permission Validation
 ************<************************************************************************************************/

func ValidatePermission(v *validator.Validator, permission *Permission) {
	// Code must be provided
	v.Check(permission.Code != "", "code", "must be provided")
	// Code must not exceed 100 characters
	v.Check(len(permission.Code) <= 100, "code", "must not exceed 100 characters")
	// Description must not exceed 500 characters
	v.Check(len(permission.Description) <= 500, "description", "must not exceed 500 characters")
}

/************************************************************************************************************
 * Permission Database Operations
 ************<************************************************************************************************/

// GetAllForUser retrieves all permissions for a specific user.
func (m *UserModel) GetAllPermissionsForUser(userID int64) (Permissions, error) {
	query := `
		SELECT p.code
		FROM permissions AS p
		INNER JOIN user_permissions AS up ON p.id = up.permission_id
		WHERE up.user_id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions Permissions
	for rows.Next() {
		var code string
		err := rows.Scan(&code)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, code)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

// AssignPermissionsToUser assigns a list of permissions to a user.
func (m *UserModel) AssignPermissionsToUser(userID int64, permissions ...string) error {
	query := `
		INSERT INTO user_permissions (user_id, permission_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, permission_id) DO NOTHING
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, userID, pq.Array(permissions))
	if err != nil {
		return err
	}
	return nil
}
