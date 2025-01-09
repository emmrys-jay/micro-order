package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"order-service/internal/adapter/logger"
	"order-service/internal/core/domain"
	"order-service/internal/core/port"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// OrderHandler represents the HTTP handler for order-related requests
type OrderHandler struct {
	svc      port.OrderService
	validate *validator.Validate
}

// NewOrderHandler creates a new OrderHandler instance
func NewOrderHandler(svc port.OrderService, vld *validator.Validate) *OrderHandler {
	return &OrderHandler{
		svc,
		vld,
	}
}

// CreateOrder godoc
//
//	@Summary		Create a new order
//	@Description	create a new order with all required details
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Param			domain.CreateOrderRequest	body		domain.CreateOrderRequest	true	"Order"
//	@Success		201							{object}	response					"Order created successfully"
//	@Failure		400							{object}	errorResponse				"Validation error"
//	@Failure		500							{object}	errorResponse				"Internal server error"
//	@Router			/order [post]
//	@Security		BearerAuth
func (ch *OrderHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.FromCtx(r.Context()).Error("Error decoding json body", zap.Error(err))
		handleError(w, domain.ErrInternal)
		return
	}

	if err := ch.validate.Struct(&req); err != nil {
		validationError(w, err)
		return
	}

	uInfo := r.Context().Value(authContextKey).(contextInfo)
	userId, err := primitive.ObjectIDFromHex(uInfo.ID)
	if err != nil {
		handleError(w, domain.NewBadRequestCError("invalid user id in request"))
		return
	}

	result, cerr := ch.svc.PlaceOrder(r.Context(), userId, &req)
	if cerr != nil {
		handleError(w, cerr)
		return
	}

	handleSuccessWithMessage(w, http.StatusCreated, result, "Order placed successfully for "+fmt.Sprint(len(result.OrderItems))+" product(s)")
}

// GetOrder godoc
//
//	@Summary		Get an order by id
//	@Description	fetch an order through id
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string			true	"Order id"
//	@Success		200	{object}	response		"Success"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/order/{id} [get]
//	@Security		BearerAuth
func (ch *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		logger.FromCtx(r.Context()).Error("Error parsing url order id", zap.Error(err))
		handleError(w, domain.NewBadRequestCError("Invalid order id"))
		return
	}

	result, cerr := ch.svc.GetOrder(r.Context(), id)
	if cerr != nil {
		handleError(w, cerr)
		return
	}

	handleSuccess(w, http.StatusOK, result)
}

// ListUserOrders godoc
//
//	@Summary		List all orders by a user
//	@Description	list all orders by a user
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path		string			true	"User id"
//	@Success		200		{object}	response		"Success"
//	@Failure		500		{object}	errorResponse	"Internal server error"
//	@Router			/user/{user_id}/orders [get]
//	@Security		BearerAuth
func (ch *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "user_id")
	userId, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		logger.FromCtx(r.Context()).Error("Error parsing url user id", zap.Error(err))
		handleError(w, domain.NewBadRequestCError("Invalid user id in url path"))
		return
	}

	result, cerr := ch.svc.ListUserOrders(r.Context(), userId)
	if cerr != nil {
		handleError(w, cerr)
		return
	}

	handleSuccess(w, http.StatusOK, result)
}

// UpdateOrderStatus godoc
//
//	@Summary		Update an order
//	@Description	update an order
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Param			id							path		string						true	"Order id"
//	@Param			domain.UpdateOrderRequest	body		domain.UpdateOrderRequest	true	"Order Status"
//	@Success		200							{object}	response					"Success"
//	@Failure		400							{object}	errorResponse				"Validation error"
//	@Failure		500							{object}	errorResponse				"Internal server error"
//	@Router			/order/{id} [patch]
//	@Security		BearerAuth
func (ch *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		logger.FromCtx(r.Context()).Error("Error parsing url order id", zap.Error(err))
		handleError(w, domain.NewBadRequestCError("Invalid order id"))
		return
	}

	var req domain.UpdateOrderRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.FromCtx(r.Context()).Error("Error decoding json body", zap.Error(err))
		handleError(w, domain.ErrInternal)
		return
	}

	if err := ch.validate.Struct(&req); err != nil {
		validationError(w, err)
		return
	}

	result, cerr := ch.svc.UpdateOrderStatus(r.Context(), id, &req)
	if cerr != nil {
		handleError(w, cerr)
		return
	}

	handleSuccess(w, http.StatusOK, result)
}

// CancelOrder godoc
//
//	@Summary		Cancel an order
//	@Description	cancel an order
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string			true	"Order id"
//	@Success		200	{object}	response		"Success"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/order/{id}/cancel [patch]
//	@Security		BearerAuth
func (ch *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		logger.FromCtx(r.Context()).Error("Error parsing url order id", zap.Error(err))
		handleError(w, domain.NewBadRequestCError("Invalid order id"))
		return
	}

	result, cerr := ch.svc.CancelOrder(r.Context(), id)
	if cerr != nil {
		handleError(w, cerr)
		return
	}

	handleSuccess(w, http.StatusOK, result)
}
