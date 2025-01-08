package http

import (
	"encoding/json"
	"net/http"

	"owner-service/internal/adapter/logger"
	"owner-service/internal/core/domain"
	"owner-service/internal/core/port"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// AuthHandler represents the HTTP handler for authentication-related requests
type AuthHandler struct {
	svc      port.AuthService
	validate *validator.Validate
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(svc port.AuthService, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{
		svc,
		validate,
	}
}

// Login godoc
//
//	@Summary		Login and get an access token
//	@Description	Logs in a registered user and returns an access token if the credentials are valid.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		domain.LoginRequest	true	"Login request body"
//	@Success		200		{object}	response			"Succesfully logged in"
//	@Failure		400		{object}	errorResponse		"Validation error"
//	@Failure		401		{object}	errorResponse		"Unauthorized error"
//	@Failure		500		{object}	errorResponse		"Internal server error"
//	@Router			/login [post]
func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.FromCtx(r.Context()).Error("Error decoding json body", zap.Error(err))
		validationError(w, err)
		return
	}

	result, err := ah.svc.Login(r.Context(), &req)
	if err != nil {
		handleError(w, err)
		return
	}

	handleSuccess(w, http.StatusOK, result)
}
