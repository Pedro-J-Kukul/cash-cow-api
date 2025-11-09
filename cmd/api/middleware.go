// FileName: internal/data/middleware.go
package main

import (
	"expvar"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

/***********************************************************************************************
 * handling Panics
 ************************************************************************************************/

// recoverPanic is a middleware that recovers from panics and returns a 500 Internal Server Error response.
func (app *App) recoverPanic(next http.Handler) http.Handler {
	// Return a handler function that wraps the next handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil { // Recover from panic
				w.Header().Set("Connection", "close")                // Close the connection after the response is sent
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err)) // Send a 500 Internal Server Error response
			}
		}()
		next.ServeHTTP(w, r) // Call the next handler in the chain
	})
}

/***********************************************************************************************
 * rate limiting
 ************************************************************************************************/

// rateLimit is a middleware that limits the rate of incoming requests.
func (app *App) rateLimit(next http.Handler) http.Handler {
	// client is a struct to hold information about each client
	type client struct {
		limiter  *rate.Limiter // Rate limiter for the client
		lastSeen time.Time     // Last time the client was seen
	}

	var (
		mu      sync.Mutex                 // Mutex to protect access to the clients map
		clients = make(map[string]*client) // Map to hold clients by their IP address
	)

	// Start a background goroutine to clean up old clients every minute
	go func() {
		for {
			time.Sleep(time.Minute) // Sleep for one minute
			mu.Lock()               // Lock the mutex to safely access the clients map
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute { // If the client hasn't been seen for over 3 minutes
					delete(clients, ip) // Remove the client from the map
				}
			}
			mu.Unlock() // Unlock the mutex
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.limiter.enabled { // Check if rate limiting is enabled
			ip := r.RemoteAddr // Get the client's IP address

			mu.Lock()                            // Lock the mutex to safely access the clients map
			if _, found := clients[ip]; !found { // If the client is not already in the map
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst), // Create a new rate limiter for the client
				}
			}
			clients[ip].lastSeen = time.Now() // Update the last seen time for the client
			if !clients[ip].limiter.Allow() { // Check if the client is allowed to make a request
				mu.Unlock()                         // Unlock the mutex before returning
				app.rateLimitExceededResponse(w, r) // Send a 429 Too Many Requests response
				return
			}
			mu.Unlock() // Unlock the mutex
		}
		next.ServeHTTP(w, r) // Call the next handler in the chain
	})
}

/***********************************************************************************************
 * Enabling CORS
 ************************************************************************************************/
// enableCORS is a middleware that adds CORS headers to the response.
func (app *App) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Origin")                        // Indicate that the response varies based on the Origin header
		w.Header().Set("Vary", "Access-Control-Request-Method") // Indicate that the response varies based on the Access-Control-Request-Method header

		origin := r.Header.Get("Origin") // Get the Origin header from the request

		if origin != "" {
			for i := range app.config.cors.trustedOrigins {
				// Check if the origin is in the trusted origins list
				if origin == app.config.cors.trustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin) // Allow the specific origin
					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						// Handle preflight request
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE") // Allowed methods
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type") // Allowed headers
						w.WriteHeader(http.StatusOK)                                                  // Respond with 200 OK
						return
					}
				}
			}
		}
		next.ServeHTTP(w, r) // Call the next handler in the chain
	})
}

// /***********************************************************************************************
//  * Authentication and Authorization
//  ************************************************************************************************/

// // authenticate is a middleware that checks for a valid authentication token in the request.
// func (app *App) authenticate(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Add("Vary", "Authorization") // Indicate that the response varies based on the Authorization header

// 		authorizationHeader := r.Header.Get("Authorization") // Get the Authorization header

// 		// If the Authorization header is empty, set the user in the context to anonymous and call the next handler
// 		if authorizationHeader == "" {
// 			r = app.contextSetUser(r, data.AnonymousUser) // Set the user in the context to anonymous
// 			next.ServeHTTP(w, r)                          // Call the next handler in the chain
// 			return                                        // Return to avoid further processing
// 		}

// 		// Split the Authorization header into parts
// 		headerParts := strings.Split(authorizationHeader, " ")   // Split the header by spaces
// 		if len(headerParts) != 2 || headerParts[0] != "Bearer" { // Check if the header is in the correct format
// 			app.invalidAuthenticationTokenResponse(w, r) // Send a 401 Unauthorized response
// 			return                                       // Return to avoid further processing
// 		}

// 		tokenPlaintext := headerParts[1] // Get the token part of the header

// 		// Validate the token plaintext
// 		v := validator.New()                                            // Create a new validator instance
// 		if data.ValidateTokenPlaintext(v, tokenPlaintext); !v.Valid() { // Validate the token format
// 			app.invalidAuthenticationTokenResponse(w, r) // Send a 401 Unauthorized response
// 			return                                       // Return to avoid further processing
// 		}

// 		// Get the user associated with the token
// 		user, err := app.models.User.GetForToken(data.ScopeAuthentication, tokenPlaintext) // Get the user for the token
// 		if err != nil {
// 			switch {
// 			case errors.Is(err, data.ErrRecordNotFound):
// 				app.invalidAuthenticationTokenResponse(w, r) // Send a 401 Unauthorized response if the token is not found
// 			default:
// 				app.serverErrorResponse(w, r, err) // Send a 500 Internal Server Error response for other errors
// 			}
// 			return // Return to avoid further processing
// 		}

// 		// Set the user in the request context
// 		r = app.contextSetUser(r, user) // Set the authenticated user in the context

// 		next.ServeHTTP(w, r) // Call the next handler in the chain
// 	})
// }

// // requireAuthenticatedUser is a middleware that ensures the user is authenticated.
// func (app *App) requireAuthenticatedUser(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		user := app.contextGetUser(r) // Get the user from the context

// 		if user.IsAnonymous() { // Check if the user is anonymous
// 			app.authenticationRequiredResponse(w, r) // Send a 401 Unauthorized response
// 			return                                   // Return to avoid further processing
// 		}

// 		next.ServeHTTP(w, r) // Call the next handler in the chain
// 	})
// }

// // requireActivatedUser is a middleware that ensures the user is activated.
// func (app *App) requireActivatedUser(next http.Handler) http.Handler {
// 	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		user := app.contextGetUser(r) // Get the user from the context

// 		if !user.IsActivated { // Check if the user is not activated
// 			app.inactiveAccountResponse(w, r) // Send a 403 Forbidden response
// 			return                            // Return to avoid further processing
// 		}

// 		next.ServeHTTP(w, r) // Call the next handler in the chain
// 	})
// 	return app.requireAuthenticatedUser(fn)
// }

// /************************************************************************************************************/
// // Permissions
// /************************************************************************************************************/
// // requireRole is a middleware that ensures the user has a specific role.
// func (app *App) requirePermissions(requiredPermissions string) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			user := app.contextGetUser(r) // Get the user from the context

// 			// Check if the user has the required permission
// 			hasPermissions, err := app.models.Role.HasPermission(user.ID, requiredPermissions)
// 			if err != nil {
// 				app.serverErrorResponse(w, r, err) // Send a 500 Internal Server Error response for errors
// 				return
// 			}

// 			if !hasPermissions {
// 				app.notPermittedResponse(w, r) // Send a 403 Forbidden response if not permitted
// 				return
// 			}

// 			next.ServeHTTP(w, r) // Call the next handler in the chain
// 		})
// 		return app.requireActivatedUser(fn) // Ensure the user is activated before checking permissions
// 	}
// }

/************************************************************************************************************/
//  Metrics
/************************************************************************************************************/

// metricsResponseWriter is a custom http.ResponseWriter that captures the status code of the response.
type metricsResponseWriter struct {
	wrapped       http.ResponseWriter // The original ResponseWriter
	statusCode    int                 // The status code of the response
	headerWritten bool                // Flag to indicate if the header has been written
}

// newMetricsResponseWriter creates a new metricsResponseWriter that wraps the given http.ResponseWriter.
func newMetricsResponseWriter(w http.ResponseWriter) *metricsResponseWriter {
	return &metricsResponseWriter{wrapped: w, statusCode: http.StatusOK} // Default status code to 200 OK
}

// Header returns the header map that will be sent by WriteHeader.
func (mw *metricsResponseWriter) Header() http.Header {
	return mw.wrapped.Header() // Delegate to the wrapped ResponseWriter
}

// WriterHeader sends an HTTP response header with the provided status code.
func (mw *metricsResponseWriter) WriteHeader(code int) {
	mw.wrapped.WriteHeader(code) // Delegate to the wrapped ResponseWriter
	if !mw.headerWritten {
		mw.statusCode = code    // Capture the status code
		mw.headerWritten = true // Mark that the header has been written
	}
}

// Write writes the data to the connection as part of an HTTP reply.
func (mw *metricsResponseWriter) Write(b []byte) (int, error) {
	mw.headerWritten = true    // Mark that the header has been written
	return mw.wrapped.Write(b) // Delegate to the wrapped ResponseWriter
}

// Unwrap returns the original http.ResponseWriter.
func (mw *metricsResponseWriter) Unwrap() http.ResponseWriter {
	return mw.wrapped // Return the wrapped ResponseWriter
}

// metrics is a middleware that collects and exposes various metrics about the HTTP requests.
func (app *App) metrics(next http.Handler) http.Handler {
	// Define variables to hold the metrics
	var (
		totalResponsesSentByStatus      = expvar.NewMap("total_responses_sent_by_status")     // Map to hold the count of responses by status code
		totalRequestsReceived           = expvar.NewInt("total_requests_received")            // Counter for total requests received
		totalResponsesSent              = expvar.NewInt("total_responses_sent")               // Counter for total responses sent
		totalProcessingTimeMicroseconds = expvar.NewInt("total_processing_time_microseconds") // Counter for total processing time in microseconds
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()                                            // Record the start time of the request
		totalRequestsReceived.Add(1)                                   // Increment the total requests received counter
		mw := newMetricsResponseWriter(w)                              // Create a new metrics response writer
		next.ServeHTTP(mw, r)                                          // Call the next handler in the chain
		totalResponsesSent.Add(1)                                      // Increment the total responses sent counter
		totalResponsesSentByStatus.Add(strconv.Itoa(mw.statusCode), 1) // Increment the count for the specific status code
		duration := time.Since(start).Microseconds()                   // Calculate the processing time in microseconds
		totalProcessingTimeMicroseconds.Add(duration)                  // Add the processing time to the total
	})
}
