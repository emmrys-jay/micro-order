package http

import (
	"encoding/json"
	"net/http"

	"owner-service/internal/adapter/logger"
	"owner-service/internal/core/domain"
	"owner-service/internal/core/port"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// UserHandler represents the HTTP handler for user-related requests
type UserHandler struct {
	svc      port.UserService
	validate *validator.Validate
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(svc port.UserService, vld *validator.Validate) *UserHandler {
	return &UserHandler{
		svc,
		vld,
	}
}

// RegisterUser godoc
//
//	@Summary		Register a new user
//	@Description	register a new user with all required details
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			domain.CreateUserRequest	body		domain.CreateUserRequest	true	"User"
//	@Success		201							{object}	response					"User created successfully"
//	@Failure		400							{object}	errorResponse				"Validation error"
//	@Failure		409							{object}	errorResponse				"Conflict error"
//	@Failure		500							{object}	errorResponse				"Internal server error"
//	@Router			/user [post]
func (ch *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.FromCtx(r.Context()).Error("Error decoding json body", zap.Error(err))
		handleError(w, domain.ErrInternal)
		return
	}

	if err := ch.validate.Struct(&req); err != nil {
		validationError(w, err)
		return
	}

	result, cerr := ch.svc.RegisterUser(r.Context(), &req)
	if cerr != nil {
		handleError(w, cerr)
		return
	}

	handleSuccessWithMessage(w, http.StatusCreated, result, "User created successfully")
}

// GetUser godoc
//
//	@Summary		Get a user by id
//	@Description	fetch a user through id
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string			true	"User id"
//	@Success		200	{object}	response		"Success"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/user/{id} [get]
//	@Security		BearerAuth
func (ch *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		logger.FromCtx(r.Context()).Error("Error parsing url uid", zap.Error(err))
		handleError(w, domain.NewBadRequestCError("Invalid user id"))
		return
	}

	result, cerr := ch.svc.GetUser(r.Context(), id)
	if cerr != nil {
		handleError(w, cerr)
		return
	}

	handleSuccess(w, http.StatusOK, result)
}

// ListUsers godoc
//
//	@Summary		List all users
//	@Description	list all registered active users
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response		"Success"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/users [get]
//	@Security		BearerAuth
func (ch *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	results, cerr := ch.svc.ListUsers(r.Context())
	if cerr != nil {
		handleError(w, cerr)
		return
	}

	handleSuccess(w, http.StatusOK, results)
}

// UpdateUser godoc
//
//	@Summary		Update a user
//	@Description	update a user
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id							path		string						true	"User id"
//	@Param			domain.UpdateUserRequest	body		domain.UpdateUserRequest	true	"User"
//	@Success		200							{object}	response					"Success"
//	@Failure		400							{object}	errorResponse				"Validation error"
//	@Failure		500							{object}	errorResponse				"Internal server error"
//	@Router			/user/{id} [patch]
//	@Security		BearerAuth
func (ch *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		logger.FromCtx(r.Context()).Error("Error parsing url uid", zap.Error(err))
		handleError(w, domain.NewBadRequestCError("Invalid user id"))
		return
	}

	var req domain.UpdateUserRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.FromCtx(r.Context()).Error("Error decoding json body", zap.Error(err))
		handleError(w, domain.ErrInternal)
		return
	}

	if err := ch.validate.Struct(&req); err != nil {
		validationError(w, err)
		return
	}

	result, cerr := ch.svc.UpdateUser(r.Context(), id, &req)
	if cerr != nil {
		handleError(w, cerr)
		return
	}

	handleSuccess(w, http.StatusOK, result)
}

// DeleteUser godoc
//
//	@Summary		Delete a user by id
//	@Description	delete a user through id
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string			true	"User id"
//	@Success		200	{object}	response		"Success"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		404	{object}	errorResponse	"Not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/user/{id} [delete]
//	@Security		BearerAuth
func (ch *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		logger.FromCtx(r.Context()).Error("Error parsing url uid", zap.Error(err))
		handleError(w, domain.NewBadRequestCError("Invalid user id"))
		return
	}

	cerr := ch.svc.DeleteUser(r.Context(), id)
	if cerr != nil {
		handleError(w, cerr)
		return
	}

	handleSuccessWithMessage(w, http.StatusOK, nil, "Deleted user successfully")
}
