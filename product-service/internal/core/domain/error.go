package domain

import (
	"net/http"
)

// CError is a custom error interface that also wraps a status code
type CError interface {
	Code() int
	Error() string
}

type err struct {
	code    int
	message string
}

func (e err) Code() int {
	return e.code
}

func (e err) Error() string {
	return e.message
}

// NewCError returns a new custom error from code and message
func NewCError(code int, message string) CError {
	return err{code, message}
}

// NewUnauthorizedError returns a new custom unauthorized error from message
func NewUnauthorizedCError(message string) CError {
	return err{http.StatusUnauthorized, message}
}

// NewInternalError returns a new custom internal error from message
func NewInternalCError(message string) CError {
	return err{http.StatusInternalServerError, message}
}

// NewBadRequestCError returns a new custom bad request error from message
func NewBadRequestCError(message string) CError {
	return err{http.StatusBadRequest, message}
}

var (
	// ErrInternal is an error for when an internal service fails to process the request
	ErrInternal = NewCError(http.StatusInternalServerError, "internal server error")
	// ErrDataNotFound is an error for when requested data is not found
	ErrDataNotFound = NewCError(http.StatusNotFound, "data not found")
	// ErrConflictingData is an error for when data conflicts with existing data
	ErrConflictingData = NewCError(http.StatusConflict, "data conflicts with existing data in unique column")
	// ErrForeignKeyViolation is an error for when there is a foreign key violation
	ErrForeignKeyViolation = NewCError(http.StatusConflict, "some of the specified ids were not found")
	// ErrInsufficientPayment is an error for when total paid is less than total price
	ErrTokenDuration = NewUnauthorizedCError("invalid token duration format")
	// ErrTokenCreation is an error for when the token creation fails
	ErrTokenCreation = NewUnauthorizedCError("error creating token")
	// ErrExpiredToken is an error for when the access token is expired
	ErrExpiredToken = NewUnauthorizedCError("access token has expired")
	// ErrInvalidToken is an error for when the access token is invalid
	ErrInvalidToken = NewUnauthorizedCError("access token is invalid")
	// ErrEmptyAuthorizationHeader is an error for when the authorization header is empty
	ErrEmptyAuthorizationHeader = NewUnauthorizedCError("authorization header is not provided")
	// ErrInvalidAuthorizationHeader is an error for when the authorization header is invalid
	ErrInvalidAuthorizationHeader = NewUnauthorizedCError("authorization header format is invalid")
	// ErrInvalidAuthorizationType is an error for when the authorization type is invalid
	ErrInvalidAuthorizationType = NewUnauthorizedCError("authorization type is not supported")
	// ErrUnauthorized is an error for when the user is unauthorized
	ErrUnauthorized = NewUnauthorizedCError("user is unauthorized to access the resource")
	// ErrInvalidCredentials is an error for when the credentials are invalid
	ErrInvalidCredentials = NewUnauthorizedCError("invalid email or password")
)
