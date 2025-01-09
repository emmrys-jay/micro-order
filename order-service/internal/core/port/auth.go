package port

import (
	"order-service/internal/core/domain"
)

// TokenService is an interface for interacting with token-related business logic
type TokenService interface {
	// VerifyToken verifies the token and returns the payload
	VerifyToken(tokenString string) (domain.Claims, error)
}
