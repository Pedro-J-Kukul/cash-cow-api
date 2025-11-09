// File: cmd/api/helpers.go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/filters"
	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/validator"
	"github.com/julienschmidt/httprouter"
)

// envelope is a generic map for wrapping data in JSON responses
type envelope map[string]any

// errorResponse represents the structure of error responses
type errorResponse struct {
	Error string `json:"error"`
}

// writes it to the response writer
func (a *App) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// convert the data to JSON
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n') // add a newline for readability

	// add any headers provided in the headers parameter
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json") // set the content type header
	w.WriteHeader(status)                              // write the status code
	_, err = w.Write(js)                               // write the JSON response body
	if err != nil {
		return err // return any error encountered while writing the response
	}

	return nil // we're done, return nil error
}

// readJSON decodes a JSON request body into the provided destination struct
func (a *App) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576                                    // limit request body to 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes)) // enforce the limit

	dec := json.NewDecoder(r.Body) // create a new JSON decoder
	dec.DisallowUnknownFields()    // disallow unknown fields in the JSON

	err := dec.Decode(dst) // decode the JSON into the destination struct
	if err != nil {
		var (
			syntaxError        *json.SyntaxError           // JSON syntax error
			unmarshalTypeError *json.UnmarshalTypeError    // JSON type error
			invalidUnmarshal   *json.InvalidUnmarshalError // invalid unmarshal error
			maxBytesError      *http.MaxBytesError         // request body too large error
		)

		switch {
		case errors.As(err, &syntaxError): // catch syntax errors
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF): // catch unexpected EOF errors
			return fmt.Errorf("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError): // catch type errors
			if unmarshalTypeError.Field != "" { // if the error is related to a specific field
				return fmt.Errorf("body contains badly-formed JSON (field %q)", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF): // catch empty body error
			return fmt.Errorf("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "): // catch unknown field errors
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case errors.As(err, &maxBytesError): // catch request body too large error
			return fmt.Errorf("body exceeds maximum size of 1MB")
		case errors.As(err, &invalidUnmarshal):
			panic(err) // this is a programmer error, so panic
		default:
			return err // return the original error for any other cases
		}
	}

	// err = dec.Decode(&struct{}{}) // check for multiple JSON objects in the body
	// if errors.Is(err, io.EOF) {
	// 	return fmt.Errorf("body must only contain a single JSON value") // return error if there's more than one JSON object
	// }

	return nil
}

/************************************************************************************************************/
// Helper functions for reading URL parameters
/************************************************************************************************************/

// readIDParameter extracts and validates an "id" parameter from the URL
func (a *App) readIDParameter(r *http.Request) (int64, error) {

	params := httprouter.ParamsFromContext(r.Context()) // get the URL parameters from the request context

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64) // parse the "id" parameter as a base-10 int64
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter") // return an error if parsing fails or id is less than 1
	}

	return id, nil // return the valid id
}

// getSingleQueryParameter retrieves a single query parameter from the URL, returning a default value if not found
func (a *App) getSingleQueryParameter(params url.Values, key string, defaultValue string) string {
	result := params.Get(key) // get the value of the specified query parameter
	if result == "" {
		return defaultValue // return the default value if the parameter is not found
	}
	return result // return the parameter value
}

// getMultipleQueryParameter retrieves multiple values for a query parameter from the URL, returning a default slice if not found
func (a *App) getMultipleQueryParameter(params url.Values, key string, defaultValue []string) []string {
	result := params.Get(key) // get the values of the specified query parameter
	if result == "" {
		return defaultValue // return the default slice if the parameter is not found
	}
	return strings.Split(result, ",") // split the parameter value by commas and return the resulting slice
}

// getSingleIntQueryParameter retrieves and validates a single integer query parameter from the URL, returning a default value if not found or invalid
func (a *App) getSingleIntQueryParameter(params url.Values, key string, defaultValue int, v *validator.Validator) int {
	result := params.Get(key) // get the value of the specified query parameter
	if result == "" {
		return defaultValue // return the default value if the parameter is not found
	}

	i, err := strconv.Atoi(result) // attempt to convert the parameter value to an integer
	if err != nil {
		v.AddError(key, "must be an integer value") // add a validation error if conversion fails
		return defaultValue                         // return the default value in case of error
	}

	return i // return the valid integer value
}

// getOptionalBoolQueryParameter retrieves a boolean query parameter returning a pointer if present.
func (a *App) getOptionalBoolQueryParameter(params url.Values, key string, v *validator.Validator) *bool {
	value := params.Get(key)
	if value == "" {
		return nil
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		v.AddError(key, "must be true or false")
		return nil
	}

	return &b
}

// getOptionalInt64QueryParameter retrieves an int64 query parameter returning a pointer if present.
func (a *App) getOptionalInt64QueryParameter(params url.Values, key string, v *validator.Validator) *int64 {
	value := params.Get(key)
	if value == "" {
		return nil
	}

	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return nil
	}

	return &i
}

// readFilters constructs a Filters struct using standard query parameters and validates it.
func (a *App) readFilters(query url.Values, defaultSort string, defaultPageSize int, safelist []string, v *validator.Validator) filters.Filters {
	f := filters.Filters{
		Page:         a.getSingleIntQueryParameter(query, "page", 1, v),
		PageSize:     a.getSingleIntQueryParameter(query, "page_size", defaultPageSize, v),
		Sort:         a.getSingleQueryParameter(query, "sort", defaultSort),
		SortSafelist: safelist,
	}

	filters.ValidateFilters(v, f)
	return f
}

/************************************************************************************************************/
// Go routine helper functions
/************************************************************************************************************/
// background runs a function in the background as a goroutine, recovering from any panics and logging them
func (a *App) background(fn func()) {
	a.wg.Add(1) // increment the wait group counter

	go func() {
		defer a.wg.Done() // decrement the wait group counter when the goroutine completes

		// recover from any panics and log the error
		defer func() {
			if err := recover(); err != nil {
				a.logger.Error("panic recovered in background goroutine", slog.Any("error", err)) // log the panic error
			}
		}()

		fn() // execute the provided function
	}()
}
