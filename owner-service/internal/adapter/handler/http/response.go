package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"owner-service/internal/core/domain"

	"github.com/go-playground/validator/v10"
)

// response represents a response body format
type response struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Success"`
	Data    any    `json:"data,omitempty"`
}

// newResponse is a helper function to create a response body
func newResponse(success bool, message string, data any) response {
	return response{
		Success: success,
		Message: message,
		Data:    data,
	}
}

// validationError sends an error response for some specific request validation error
func validationError(w http.ResponseWriter, err error) {
	errMsgs := parseError(err)
	errRsp := newErrorResponse(errMsgs)
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(errRsp)
}

// parseError parses error messages from the error object and returns a slice of error messages
func parseError(err error) []string {
	var errMsgs []string

	if errors.As(err, &validator.ValidationErrors{}) {
		for _, err := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, err.Error())
		}
	} else {
		errMsgs = append(errMsgs, err.Error())
	}

	return errMsgs
}

// errorResponse represents an error response body format
type errorResponse struct {
	Success  bool   `json:"success" example:"false"`
	Messages string `json:"messages" example:"Error message 1 - Error message 2"`
}

// newErrorResponse is a helper function to create an error response body
func newErrorResponse(errMsgs []string) errorResponse {
	msgs := ""
	for i, v := range errMsgs {
		if i == len(errMsgs)-1 {
			msgs += v
			continue
		}
		msgs += v + " - "
	}

	return errorResponse{
		Success:  false,
		Messages: msgs,
	}
}

// handleSuccess sends a success response with the specified status code and optional data
func handleSuccess(w http.ResponseWriter, code int, data any) {
	rsp := newResponse(true, "Success", data)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(rsp)
}

// handleSuccessWithMessage sends a success response with the specified status code, optional data and message
func handleSuccessWithMessage(w http.ResponseWriter, code int, data any, message string) {
	rsp := newResponse(true, message, data)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(rsp)
}

// handleError determines the status code of an error and returns a JSON response with the error message and status code
func handleError(w http.ResponseWriter, err domain.CError) {
	// TODO: Change the type of error received and the mech to get the code
	statusCode := err.Code()
	errMsg := parseError(err)
	errRsp := newErrorResponse(errMsg)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errRsp)
}
