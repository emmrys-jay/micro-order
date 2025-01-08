package port

import (
	"context"

	"owner-service/internal/core/domain"
)

// PingRepository is an interface for interacting with ping-related data
type PingRepository interface {
	// CreatePing inserts a new user into the database
	CreatePing(ctx context.Context, user *domain.Ping) error
}

// PingService is an interface for interacting with ping-related business logic
type PingService interface {
	// CreateCategory creates a new category
	Ping(ctx context.Context, ping *domain.Ping) (domain.Ping, domain.CError)
}
