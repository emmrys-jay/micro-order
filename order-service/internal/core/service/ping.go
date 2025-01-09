package service

import (
	"context"

	"order-service/internal/core/domain"
	"order-service/internal/core/port"
)

/**
 * PingService implements port.PingService interface
 */
type PingService struct {
	repo  port.PingRepository
	cache port.CacheRepository
}

// NewAuthService creates a new auth service instance
func NewPingService(repo port.PingRepository, cache port.CacheRepository) *PingService {
	return &PingService{
		repo,
		cache,
	}
}

// Login gives a registered user an access token if the credentials are valid
func (ps *PingService) Ping(ctx context.Context, ping *domain.Ping) (domain.Ping, domain.CError) {
	_ = ps.repo.CreatePing(ctx, ping)
	return *ping, nil
}
