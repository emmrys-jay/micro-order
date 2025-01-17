package port

import (
	"context"

	"owner-service/internal/core/domain"
)

// TokenService is an interface for interacting with token-related business logic
type TokenService interface {
	// CreateToken creates a new token for a given user
	CreateToken(id, email string, role string) (string, error)
	// VerifyToken verifies the token and returns the payload
	VerifyToken(tokenString string) (domain.Claims, domain.CError)
}

// UserService is an interface for interacting with user authentication-related business logic
type AuthService interface {
	// Login authenticates a user by email and password and returns a token
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, domain.CError)
}
