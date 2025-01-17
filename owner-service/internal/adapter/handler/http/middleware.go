package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"owner-service/internal/adapter/logger"
	"owner-service/internal/core/domain"
	"owner-service/internal/core/port"

	"github.com/rs/xid"
	"go.uber.org/zap"
)

const (
	// authorizationHeaderKey is the key for authorization header in the request
	authorizationHeaderKey = "authorization"
	// authorizationType is the accepted authorization type
	authorizationType = "bearer"
	// authorizationPayloadKey is the key for authorization payload in the context
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(next http.Handler, token port.TokenService, logger *zap.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := getToken(r, authorizationHeaderKey)
		if tokenString == "" {
			handleError(w, domain.ErrEmptyAuthorizationHeader)
			return
		}

		fields := strings.Fields(tokenString)
		isValid := len(fields) == 2
		if !isValid {
			handleError(w, domain.ErrInvalidAuthorizationType)
			return
		}

		claims, err := token.VerifyToken(fields[1])
		if err != nil {
			logger.Error("error verifying token", zap.Error(err))
			handleError(w, err)
			return
		}

		// Set details from token in context
		ctx := context.WithValue(r.Context(), domain.AuthContextKey, contextInfo{
			ID:    claims.ID,
			Email: claims.Email,
		})

		// call the next handler in the chain, passing the response writer and
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type contextInfo struct {
	ID    string
	Role  string
	Email string
}

func getToken(r *http.Request, header string) string {
	return r.Header.Get(header)
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := logger.Get()

		correlationID := xid.New().String()

		ctx := context.WithValue(
			r.Context(),
			domain.CorrelationIDCtxKey,
			correlationID,
		)

		r = r.WithContext(ctx)

		ctx = logger.WithCtx(ctx, l.With(
			zap.String(string(domain.CorrelationIDCtxKey), correlationID),
		))
		w.Header().Add("X-Correlation-ID", correlationID)

		lrw := newLoggingResponseWriter(w)

		r = r.WithContext(ctx)

		defer func(start time.Time) {
			l.Info(
				fmt.Sprintf(
					"%s request to %s completed",
					r.Method,
					r.RequestURI,
				),
				zap.String("method", r.Method),
				zap.String("url", r.RequestURI),
				// zap.String("user_agent", r.UserAgent()),
				zap.Int("status_code", lrw.statusCode),
				zap.Duration("elapsed_ms", time.Since(start)),
			)
		}(time.Now())

		lrw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(lrw, r)
	})
}

func adminMiddleware(next http.Handler, token port.TokenService, logger *zap.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := getToken(r, authorizationHeaderKey)
		if tokenString == "" {
			handleError(w, domain.ErrEmptyAuthorizationHeader)
			return
		}

		fields := strings.Fields(tokenString)
		isValid := len(fields) == 2
		if !isValid {
			handleError(w, domain.ErrInvalidAuthorizationHeader)
			return
		}

		claims, err := token.VerifyToken(fields[1])
		if err != nil {
			logger.Error("error verifying token", zap.Error(err))
			handleError(w, err)
			return
		}

		// claims.Email is of the form <email,role>
		identifier := strings.Split(claims.Email, ",")
		if len(identifier) != 2 {
			handleError(w, domain.ErrInvalidToken)
			return
		}

		email, role := identifier[0], identifier[len(identifier)-1]
		if role != domain.RAdmin.String() {
			handleError(w, domain.ErrUnauthorized)
		}

		// Set details from token in context
		ctx := context.WithValue(r.Context(), domain.AuthContextKey, contextInfo{
			ID:    claims.ID,
			Email: email,
		})

		// call the next handler in the chain, passing the response writer and
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type contextKey string

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
