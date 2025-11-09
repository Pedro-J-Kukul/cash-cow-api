// File: routes.go
package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// routes sets up the HTTP routes for the application
func (app *App) routes() http.Handler {
	router := httprouter.New() // create a new router instance

	return router
}
