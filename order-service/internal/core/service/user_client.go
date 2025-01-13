package service

import (
	"context"
	"fmt"
	"order-service/internal/adapter/config"
	"order-service/internal/adapter/logger"
	"order-service/internal/core/domain"
	"order-service/internal/core/service/user"
	"order-service/internal/core/util"

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

func (os *OrderService) GetUser(ctx context.Context, userId primitive.ObjectID) (exists bool, cerr domain.CError) {
	log := logger.FromCtx(ctx)

	log.Info("Checking if user exists in cache", zap.String("product_id", userId.Hex()))
	cacheKey := util.GenerateCacheKey("user", userId.Hex())
	_, err := os.cache.Get(ctx, cacheKey)
	if err == nil {
		log.Info("User with id found in cache", zap.String("user_id", userId.Hex()))
		return true, nil
	}

	grpcConn, grpcClient, err := newUserClient(&config.GetConfig().Discovery)
	if err != nil {
		log.Error("Error creating user client", zap.Error(err))
		return false, domain.ErrInternal
	}
	defer grpcConn.Close()

	log.Info("Created new grpc user client")
	log.Info("Making request to fetch user", zap.String("user_id", userId.Hex()))

	user, err := grpcClient.Get(context.Background(), &user.UserRequest{UserId: userId.Hex()})
	if err != nil {
		log.Error("Error fetching user", zap.Error(err))
		return false, domain.ErrInternal
	}
	log.Info("Successfully fetched user")

	log.Info("Saving user to cache")
	sUser, err := util.Serialize(user)
	if err != nil {
		log.Error("Error serializing user", zap.Error(err))
		return false, domain.ErrInternal
	}

	err = os.cache.Set(ctx, cacheKey, sUser, os.cacheTtl)
	if err != nil {
		log.Error("Error saving returned user to cache", zap.Error(err))
		return false, domain.ErrInternal
	}

	return true, nil
}
