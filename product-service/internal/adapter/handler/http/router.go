package http

import (
	"net/http"
	"strings"

	"product-service/internal/adapter/config"
	"product-service/internal/core/port"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// Router is a wrapper for HTTP router
type Router struct {
	chi.Router
}

// NewRouter creates a new HTTP router
func NewRouter(
	config *config.ServerConfiguration,
	token port.TokenService,
	logger *zap.Logger,
	pingHandler PingHandler,
	productHandler ProductHandler,
) (*Router, error) {

	// CORS
	corsConfig := cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}

	allowedOrigins := config.HttpAllowedOrigins
	if allowedOrigins != "" {
		originsList := strings.Split(config.HttpAllowedOrigins, ",")
		corsConfig.AllowedOrigins = originsList
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(corsConfig))

	// Logger
	router.Use(requestLogger)
	router.Use(middleware.Recoverer)

	// Swagger
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("0.0.0.0:"+config.HttpPort+"/swagger/doc.json"), //The url pointing to API definition
	))

	// v1
	router.Route("/api/v1", func(r chi.Router) {

		// Ping
		r.Route("/health", func(r chi.Router) {
			r.Get("/", pingHandler.PingGet)
			r.Post("/", pingHandler.PingPost)
		})

		// Product
		r.Route("/product", func(r chi.Router) {
			r.Post("/", adminMiddleware(http.HandlerFunc(productHandler.CreateProduct), token, logger))
			r.Patch("/{id}", adminMiddleware(http.HandlerFunc(productHandler.UpdateProduct), token, logger))
			r.Delete("/{id}", adminMiddleware(http.HandlerFunc(productHandler.DeleteProduct), token, logger))

			r.Get("/{id}", authMiddleware(http.HandlerFunc(productHandler.GetProduct), token, logger))
		})
		r.Get("/products", authMiddleware(http.HandlerFunc(productHandler.ListProducts), token, logger))
	})

	return &Router{
		router,
	}, nil
}
