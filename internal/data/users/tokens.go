// file: internal/data/tokens.go
package users

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"time"

	internalErrors "github.com/Pedro-J-Kukul/cash-cow-api/internal/data/errors"
	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/validator"
)

/************************************************************************************************************
 * Token Declarations
 ************************************************************************************************************/

// Scope constants for different token types
const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
	ScopePasswordReset  = "password_reset"
)

// Token represents a user token for various purposes like activation, authentication, or password reset
type Token struct {
	Plaintext string    `json:"plaintext"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

type TokenModel struct {
	DB *sql.DB
}

/************************************************************************************************************
 * Token Helpers
 ************************************************************************************************************/

// generateToken generates a new token for a user with a specific scope and expiry duration
func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	// Create a new token with the provided userID, scope, and calculated expiry time
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	// Generate a random plaintext token
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, err // Return error if random byte generation fails
	}
	token.Plaintext = base64.RawURLEncoding.EncodeToString(randomBytes) // Encode to URL-safe base64
	hash := sha256.Sum256([]byte(token.Plaintext))                      // Hash the plaintext token

	token.Hash = hash[:] // Set the hash field to the hashed value
	return token, nil    // Return the generated token
}

// ValidateTokenPlaintext checks if a token is valid based on its plaintext value
func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")                // Token must be provided
	v.Check(len(tokenPlaintext) == 22, "token", "must be 22 characters long") // Token must match expected length
}

/************************************************************************************************************
 * Database Operations
 ************************************************************************************************************/

// TokenRepository is an interface for token-related database operations
type TokenRepository interface {
	GetUserToken(tokenScope, tokenPlaintext string) (*User, error)
	DeleteAllForUser(scope string, userID int64) error
	Insert(token *Token) error
	New(userID int64, ttl time.Duration, scope string) (*Token, error)
}

// GetUserToken Method
func (m *TokenModel) GetUserToken(tokenScope, tokenPlaintext string) (*User, error) {
	// Query
	query := `
		SELECT u.id, u.farmer_id, u.email, u.phone_number, u.first_name, u.last_name, u.password_hash, u.is_activated, u.is_deleted, u.is_verified, u.version, u.created_at, u.updated_at
		FROM users AS u
		INNER JOIN tokens AS t ON u.id = t.user_id
		WHERE t.scope = $1 AND t.plaintext = $2 AND t.expiry > $3`

	// Hash the token plaintext
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	// Arguments for Query
	args := []any{
		tokenScope,
		tokenHash[:],
		time.Now(),
	}

	// Declare User variable
	var user User

	// Declare Scan variables
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

	// Get Context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Execute Query
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(scan...)
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

// DeleteAllForUser removes all tokens for a specific user and scope from the database
func (m *TokenModel) DeleteAllForUser(scope string, userID int64) error {
	// SQL query to delete tokens for a specific user and scope
	query := `
		DELETE FROM tokens
		WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // Create a context with a 3-second timeout
	defer cancel()                                                          // Ensure the context is cancelled to free resources

	_, err := m.DB.ExecContext(ctx, query, scope, userID) // Execute the delete query
	return err                                            // Return any error that occurred during execution
}

// Insert adds a new token to the database
func (m *TokenModel) Insert(token *Token) error {
	// SQL query to insert a new token into the tokens table
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4)`

	// Prepare the arguments for the query
	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // Create a context with a 3-second timeout
	defer cancel()                                                          // Ensure the context is cancelled to free resources

	_, err := m.DB.ExecContext(ctx, query, args...) // Execute the insert query
	return err                                      // Return any error that occurred during execution
}

// New creates a new token for a user and stores it in the database
func (m *TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	// Generate a new token
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err // Return error if token generation fails
	}

	err = m.Insert(token) // Insert the token into the database
	return token, err     // Return the token and any insertion error
}
