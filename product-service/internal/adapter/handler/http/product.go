package http

import (
	"encoding/json"
	"net/http"

	"product-service/internal/adapter/logger"
	"product-service/internal/core/domain"
	"product-service/internal/core/port"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// ProductHandler represents the HTTP handler for product-related requests
type ProductHandler struct {
	svc      port.ProductService
	validate *validator.Validate
}

// NewProductHandler creates a new ProductHandler instance
func NewProductHandler(svc port.ProductService, vld *validator.Validate) *ProductHandler {
	return &ProductHandler{
		svc,
		vld,
	}
}

// CreateProduct godoc
//
//	@Summary		Create a new product
//	@Description	create a new product with all required details
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			domain.CreateProductRequest	body		domain.CreateProductRequest	true	"Product"
//	@Success		201							{object}	response					"Product created successfully"
//	@Failure		400							{object}	errorResponse				"Validation error"
//	@Failure		500							{object}	errorResponse				"Internal server error"
//	@Router			/product [post]
//	@Security		BearerAuth
func (ch *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.FromCtx(r.Context()).Error("Error decoding json body", zap.Error(err))
		handleError(w, domain.ErrInternal)
		return
	}

	if err := ch.validate.Struct(&req); err != nil {
		validationError(w, err)
		return
	}

	result, cerr := ch.svc.CreateProduct(r.Context(), &req)
	if cerr != nil {
		handleError(w, cerr)
		return
	}

	handleSuccessWithMessage(w, http.StatusCreated, result, "Product created successfully")
}

// GetProduct godoc
//
//	@Summary		Get a product by id
//	@Description	fetch a product through id
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string			true	"Product id"
//	@Success		200	{object}	response		"Success"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/product/{id} [get]
//	@Security		BearerAuth
func (ch *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		logger.FromCtx(r.Context()).Error("Error parsing url product id", zap.Error(err))
		handleError(w, domain.NewBadRequestCError("Invalid product id"))
		return
	}

	result, cerr := ch.svc.GetProduct(r.Context(), id)
	if cerr != nil {
		handleError(w, cerr)
		return
	}

	handleSuccess(w, http.StatusOK, result)
}

// ListProducts godoc
//
//	@Summary		List all products
//	@Description	list all active products
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response		"Success"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/products [get]
//	@Security		BearerAuth
func (ch *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	result, cerr := ch.svc.ListProducts(r.Context())
	if cerr != nil {
		handleError(w, cerr)
		return
	}

	handleSuccess(w, http.StatusOK, result)
}

// UpdateProduct godoc
//
//	@Summary		Update a product
//	@Description	update a product
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			id							path		string						true	"Product id"
//	@Param			domain.UpdateProductRequest	body		domain.UpdateProductRequest	true	"Product"
//	@Success		200							{object}	response					"Success"
//	@Failure		400							{object}	errorResponse				"Validation error"
//	@Failure		500							{object}	errorResponse				"Internal server error"
//	@Router			/product/{id} [patch]
//	@Security		BearerAuth
func (ch *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		logger.FromCtx(r.Context()).Error("Error parsing url product id", zap.Error(err))
		handleError(w, domain.NewBadRequestCError("Invalid product id"))
		return
	}

	var req domain.UpdateProductRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.FromCtx(r.Context()).Error("Error decoding json body", zap.Error(err))
		handleError(w, domain.ErrInternal)
		return
	}

	if err := ch.validate.Struct(&req); err != nil {
		validationError(w, err)
		return
	}

	result, cerr := ch.svc.UpdateProduct(r.Context(), id, &req)
	if cerr != nil {
		handleError(w, cerr)
		return
	}

	handleSuccess(w, http.StatusOK, result)
}

// DeleteProduct godoc
//
//	@Summary		Delete a product by id
//	@Description	delete a product through id
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string			true	"Product id"
//	@Success		200	{object}	response		"Success"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		404	{object}	errorResponse	"Not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/product/{id} [delete]
//	@Security		BearerAuth
func (ch *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		logger.FromCtx(r.Context()).Error("Error parsing url product id", zap.Error(err))
		handleError(w, domain.NewBadRequestCError("Invalid product id"))
		return
	}

	cerr := ch.svc.DeleteProduct(r.Context(), id)
	if cerr != nil {
		handleError(w, cerr)
		return
	}

	handleSuccessWithMessage(w, http.StatusOK, nil, "Deleted product successfully")
}
