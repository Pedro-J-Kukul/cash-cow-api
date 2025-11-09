// File: cmd/api/context.go
package main

import (
	"context"
	"net/http"

	data "github.com/Pedro-J-Kukul/cash-cow-api/internal/data/users"
)

// contextKey is a custom type for context keys defined in this package.
type contextKey string

// Define a key for storing the authenticated user in the request context.
const userContextKey = contextKey("user")

// contextSetUser adds the authenticated user to the request context.
func (a *App) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// contextGetUser retrieves the authenticated user from the request context.
func contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in context")
	}
	return user
}
