package service

import (
	"context"
	"fmt"
	"product-service/internal/adapter/config"
	"product-service/internal/adapter/logger"
	"product-service/internal/core/domain"
	"product-service/internal/core/service/user"
	"product-service/internal/core/util"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newUserClient(conf *config.DiscoveryConfiguration) (*grpc.ClientConn, user.UserClient, error) {
	conn, err := grpc.NewClient(conf.OwnerUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create a new connection: %w", err)
	}

	userClient := user.NewUserClient(conn)
	return conn, userClient, nil
}

func (ps *ProductService) GetUser(ctx context.Context, userId primitive.ObjectID) (*user.UserResponse, domain.CError) {
	log := logger.FromCtx(ctx)
	userIdString := userId.Hex()

	log.Info("Checking if user exists in cache", zap.String("user_id", userIdString))

	cacheKey := util.GenerateCacheKey("user", userIdString)
	cachedUser, err := ps.cache.Get(ctx, cacheKey)
	if err == nil {
		log.Info("User with id found in cache", zap.String("user_id", userIdString))

		var user user.UserResponse
		err := util.Deserialize(cachedUser, &user)
		if err != nil {
			log.Error("Error deserializing found user in cache", zap.Error(err))
			return nil, domain.ErrInternal
		}

		return &user, nil
	}
	log.Info("User not found in cache")

	grpcConn, grpcClient, err := newUserClient(&config.GetConfig().Discovery)
	if err != nil {
		log.Error("Error creating user client", zap.Error(err))
		return nil, domain.ErrInternal
	}
	defer grpcConn.Close()

	log.Info("Created new grpc user client")
	log.Info("Making request to fetch user", zap.String("user_id", userIdString))

	user, err := grpcClient.Get(context.Background(), &user.UserRequest{UserId: userIdString})
	if err != nil {
		log.Error("Error fetching user", zap.Error(err))
		return nil, domain.ErrInternal
	}

	log.Info("Successfully fetched user")
	log.Info("Saving returned user to cache")

	serialUser, err := util.Serialize(user)
	if err != nil {
		log.Error("Error serializing user", zap.Error(err))
		return nil, domain.ErrInternal
	}

	err = ps.cache.Set(ctx, cacheKey, serialUser, ps.cacheTtl)
	if err != nil {
		log.Error("Error saving returned user to cache", zap.Error(err))
		return nil, domain.ErrInternal
	}

	log.Info("Successfully saved user to cache")

	return user, nil
}
