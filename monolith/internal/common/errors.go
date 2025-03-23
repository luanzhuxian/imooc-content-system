package common

import (
	"errors"
	"fmt"
)

// Standard domain errors - these should be used by services
var (
	// ErrUserAlreadyExists is returned when attempting to register a user with an ID that already exists
	ErrUserAlreadyExists = errors.New("user already exists")
	
	// ErrUserNotFound is returned when a user cannot be found
	ErrUserNotFound = errors.New("user not found")
	
	// ErrInvalidCredentials is returned when user authentication fails
	ErrInvalidCredentials = errors.New("invalid credentials")
	
	// ErrInternalServer is returned for unexpected server errors
	ErrInternalServer = errors.New("internal server error")
)

// DomainError wraps a standard error with additional context information
type DomainError struct {
	Err     error  // Original error
	Message string // Human-readable message
	Code    string // Error code for client-side error handling
}

// Error implements the error interface
func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

// Unwrap returns the wrapped error for errors.Is and errors.As support
func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError creates a new domain error
func NewDomainError(err error, message, code string) *DomainError {
	return &DomainError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

// NewUserAlreadyExistsError creates a specific error for user already exists
func NewUserAlreadyExistsError(userID string) *DomainError {
	return &DomainError{
		Err:     ErrUserAlreadyExists,
		Message: fmt.Sprintf("User with ID '%s' already exists", userID),
		Code:    "USER_ALREADY_EXISTS",
	}
} 