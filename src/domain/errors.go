package domain

import (
	"fmt"
	"strings"
)

// ValidationError represents an error during domain entity validation
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("validation failed: ")

	for i, err := range e {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(err.Error())
	}

	return sb.String()
}

// Add adds a validation error to the collection
func (e *ValidationErrors) Add(field, message string) {
	*e = append(*e, ValidationError{Field: field, Message: message})
}

// HasErrors returns true if there are any validation errors
func (e ValidationErrors) HasErrors() bool {
	return len(e) > 0
}

// DomainError represents a general domain error
type DomainError struct {
	Code    string
	Message string
	Cause   error
}

func (e DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewDomainError creates a new domain error
func NewDomainError(code, message string, cause error) DomainError {
	return DomainError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// Predefined error codes
const (
	ErrInvalidInput        = "INVALID_INPUT"
	ErrInvalidState        = "INVALID_STATE"
	ErrInvalidTransition   = "INVALID_TRANSITION"
	ErrResourceNotFound    = "RESOURCE_NOT_FOUND"
	ErrConcurrencyConflict = "CONCURRENCY_CONFLICT"
)
