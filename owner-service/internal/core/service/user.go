package service

import (
	"context"
	"net/http"

	"owner-service/internal/adapter/logger"
	"owner-service/internal/core/domain"
	"owner-service/internal/core/port"
	"owner-service/internal/core/util"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

/**
 * UserService implements port.UserService interface
 */
type UserService struct {
	repo  port.UserRepository
	cache port.CacheRepository
}

// NewUserService creates a new auth service instance
func NewUserService(repo port.UserRepository, cache port.CacheRepository) *UserService {
	return &UserService{
		repo,
		cache,
	}
}

func (us *UserService) RegisterUser(ctx context.Context, user *domain.CreateUserRequest) (*domain.User, domain.CError) {
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {

		logger.FromCtx(ctx).Error("Error hashing user password", zap.Error(err))
		return nil, domain.ErrInternal
	}

	userToCreate := domain.User{
		Email:     user.Email,
		Password:  hashedPassword,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		IsActive:  true,
	}

	if user.Role.String() == "" {
		userToCreate.Role = domain.RUser
	}

	userResponse, cerr := us.repo.CreateUser(ctx, &userToCreate)
	if cerr != nil {
		if cerr.Code() == 409 { // conflict
			return nil, domain.NewCError(cerr.Code(), "email already exists")
		}

		logger.FromCtx(ctx).Error("Error creating user", zap.Error(cerr))
		return nil, domain.ErrInternal
	}
	userResponse.Password = ""

	return userResponse, nil
}

func (us *UserService) GetUser(ctx context.Context, id primitive.ObjectID) (*domain.User, domain.CError) {
	user, cerr := us.repo.GetUserByID(ctx, id)
	if cerr != nil {
		if cerr.Code() == 500 {

			logger.FromCtx(ctx).Error("Error getting user", zap.Error(cerr))
			return nil, domain.ErrInternal
		}
		return nil, cerr
	}
	user.Password = ""

	return user, nil
}

func (us *UserService) ListUsers(ctx context.Context) ([]domain.User, domain.CError) {
	users, cerr := us.repo.ListUsers(ctx)
	if cerr != nil {

		logger.FromCtx(ctx).Error("Error listing user", zap.Error(cerr))
		return nil, domain.ErrInternal
	}

	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}

func (us *UserService) UpdateUser(ctx context.Context, id primitive.ObjectID, req *domain.UpdateUserRequest) (*domain.User, domain.CError) {
	retUser, cerr := us.GetUser(ctx, id)
	if cerr != nil {
		return nil, cerr
	}

	if req.FirstName == retUser.FirstName && req.LastName == retUser.LastName && req.Role == retUser.Role.String() {
		return nil, domain.NewCError(http.StatusBadRequest, "There are no changes to update")
	}

	retUser.FirstName = req.FirstName
	retUser.LastName = req.LastName

	if role, ok := domain.StringToUserRole[req.Role]; ok {
		retUser.Role = role
	}

	userResponse, cerr := us.repo.UpdateUser(ctx, retUser)
	if cerr != nil {
		if cerr.Code() == 500 {

			logger.FromCtx(ctx).Error("Error updating user", zap.Error(cerr))
			return nil, domain.ErrInternal
		}
		return nil, cerr
	}
	userResponse.Password = ""

	return userResponse, nil
}

func (us *UserService) DeleteUser(ctx context.Context, id primitive.ObjectID) domain.CError {
	cerr := us.repo.DeleteUser(ctx, id)
	if cerr != nil {
		if cerr.Code() == 500 {

			logger.FromCtx(ctx).Error("Error deleting user", zap.Error(cerr))
			return domain.ErrInternal
		}
		return cerr
	}

	return nil
}

func (us *UserService) CreateAdminUser(ctx context.Context, email, password string) domain.CError {
	user := domain.CreateUserRequest{
		Email:    email,
		Password: password,
		Role:     domain.RAdmin,
	}

	_, cerr := us.RegisterUser(ctx, &user)
	if cerr != nil {
		return cerr
	}

	return nil
}
