package service

import (
	"context"

	"owner-service/internal/adapter/logger"
	"owner-service/internal/core/domain"
	"owner-service/internal/core/port"
	"owner-service/internal/core/util"

	"go.uber.org/zap"
)

/**
 * AuthService implements port.AuthService interface
 * and provides an access to the user repository
 * and token service
 */
type AuthService struct {
	repo  port.UserRepository
	ts    port.TokenService
	cache port.CacheRepository
}

// NewAuthService creates a new auth service instance
func NewAuthService(repo port.UserRepository, ts port.TokenService, cache port.CacheRepository) *AuthService {
	return &AuthService{
		repo,
		ts,
		cache,
	}
}

// Login gives a registered user an access token if the credentials are valid
func (as *AuthService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, domain.CError) {
	user, cerr := as.repo.GetUserByEmail(ctx, req.Email)
	if cerr != nil {
		if cerr == domain.ErrDataNotFound {
			return nil, domain.ErrInvalidCredentials
		}

		logger.FromCtx(ctx).Error("Error fetching user by email", zap.Error(cerr))
		return nil, domain.ErrInternal
	}

	err := util.ComparePassword(req.Password, user.Password)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, err := as.ts.CreateToken(user.ID.Hex(), req.Email, user.Role.String())
	if err != nil {

		logger.FromCtx(ctx).Error("Error creating token", zap.Error(cerr))
		return nil, domain.ErrTokenCreation
	}

	logger.FromCtx(ctx).Info("created token for user using",
		zap.String("email", req.Email),
		zap.String("role", user.Role.String()))

	return &domain.LoginResponse{
		Token: accessToken,
		User:  *user,
	}, nil
}
