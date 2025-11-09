// File: cmd/api/server.go
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// serve starts the HTTP server and listens for incoming requests
func (app *App) serve() error {

	// Define the server configuration
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),                      // server address with port
		Handler:      app.routes(),                                             // HTTP handler for routing requests
		IdleTimeout:  1 * time.Minute,                                          // maximum idle time for connections
		ReadTimeout:  5 * time.Second,                                          // maximum duration for reading the request
		WriteTimeout: 10 * time.Second,                                         // maximum duration for writing the response
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError), // custom error logger
	}

	shutdown := make(chan error) // channel for shutdown errors

	// Start a goroutine to listen for shutdown signals
	go func() {
		quit := make(chan os.Signal, 1)                                              // channel for OS signals
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)                         // listen for interrupt and terminate signals
		sig := <-quit                                                                // block until a signal is received
		app.logger.Info("shutting down server", slog.String("signal", sig.String())) // log the shutdown signal

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // context with timeout for shutdown
		defer cancel()                                                           // ensure the context is cancelled to free resources

		err := srv.Shutdown(ctx) // attempt to gracefully shutdown the server
		if err != nil {
			shutdown <- err // send any shutdown error to the channel
		}

		app.logger.Info("completing background tasks") // log completion of background tasks
		app.wg.Wait()                                  // wait for all background tasks to complete
		shutdown <- nil                                // signal that shutdown is complete
	}()

	app.logger.Info("starting server", slog.String("env", app.config.env), slog.Int("port", app.config.port)) // log server start

	err := srv.ListenAndServe()                // start the server and listen for requests
	if !errors.Is(err, http.ErrServerClosed) { // check if the error is not due to server shutdown
		return err // return any unexpected error
	}

	err = <-shutdown // wait for shutdown to complete
	if err != nil {
		return err // return any shutdown error
	}

	app.logger.Info("server stopped") // log that the server has stopped
	return nil                        // return nil indicating successful shutdown
}
