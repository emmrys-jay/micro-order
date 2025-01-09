package http

import (
	"net/http"

	"order-service/internal/adapter/logger"
	"order-service/internal/core/domain"
	"order-service/internal/core/port"

	"github.com/go-playground/validator/v10"
)

// PingHandler represents the HTTP handler for ping-related requests
type PingHandler struct {
	svc      port.PingService
	validate *validator.Validate
}

// NewCategoryHandler creates a new CategoryHandler instance
func NewPingHandler(svc port.PingService, vld *validator.Validate) *PingHandler {
	return &PingHandler{
		svc,
		vld,
	}
}

// PingPost godoc
//
//	@Summary		Create a new ping object
//	@Description	create a new ping object with name
//	@Tags			Ping
//	@Accept			json
//	@Produce		json
//	@Param			domain.Ping	body		domain.Ping		true	"Create ping request"
//	@Success		201			{object}	response		"Ping created"
//	@Failure		400			{object}	errorResponse	"Validation error"
//	@Failure		500			{object}	errorResponse	"Internal server error"
//	@Router			/health [post]
func (ch *PingHandler) PingPost(w http.ResponseWriter, r *http.Request) {
	var req domain.Ping

	if err := ch.validate.Struct(&req); err != nil {
		validationError(w, err)
		return
	}

	ping, err := ch.svc.Ping(r.Context(), &req)
	if err != nil {
		handleError(w, err)
		return
	}

	handleSuccess(w, http.StatusCreated, ping)
}

// PingGet godoc
//
//	@Summary		Check server status
//	@Description	check server status
//	@Tags			Ping
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response	"Ping created"
//	@Router			/health [get]
func (ch *PingHandler) PingGet(w http.ResponseWriter, r *http.Request) {
	logger.FromCtx(r.Context()).Info("Alive!")
	handleSuccessWithMessage(w, 200, nil, "Server OK")
}
