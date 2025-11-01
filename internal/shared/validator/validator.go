// File: internal/shared/validator/valitator.go

package validator

import (
	"regexp"
	"slices"
)

// Validator Structure to hold validation errors
type Validator struct {
	Errors map[string]string
}

// New creates a new Validator instance
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// emailRX is a regular expression for validating email addresses
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Password complexity regular expressions
var (
	PasswordNumberRX  = regexp.MustCompile("[0-9]")
	PasswordUpperRX   = regexp.MustCompile("[A-Z]")
	PasswordLowerRX   = regexp.MustCompile("[a-z]")
	PasswordSpecialRX = regexp.MustCompile("[!@#~$%^&*()+|_]")
	PasswordMinLength = 8
	PasswordMaxLength = 72
)

// Valid checks if there are any validation entries
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message for a specific field if it doesn't already exist
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check adds an error message for a field if the condition is false
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// Matches checks if a value matches a given regular expression
func (v *Validator) Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Permitted checks if a value is in a list of permitted values
func (v *Validator) IsPermitted(value string, permittedValues ...string) bool {
	return slices.Contains(permittedValues, value)
}
