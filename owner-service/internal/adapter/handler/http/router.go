package http

import (
	"net/http"
	"strings"

	"owner-service/internal/adapter/config"
	"owner-service/internal/core/port"

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
	userHandler UserHandler,
	authHandler AuthHandler,
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

		// Auth
		r.Post("/login", authHandler.Login)

		// User
		r.Route("/user", func(r chi.Router) {
			r.Post("/", userHandler.RegisterUser)
			r.Get("/{id}", adminMiddleware(http.HandlerFunc(userHandler.GetUser), token, logger))
			r.Patch("/{id}", authMiddleware(http.HandlerFunc(userHandler.UpdateUser), token, logger))
			r.Delete("/{id}", adminMiddleware(http.HandlerFunc(userHandler.DeleteUser), token, logger))
		})
		r.Get("/users", adminMiddleware(http.HandlerFunc(userHandler.ListUsers), token, logger))
	})

	return &Router{
		router,
	}, nil
}
