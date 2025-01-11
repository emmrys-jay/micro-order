package service

import (
	"context"
	"fmt"
	"product-service/internal/adapter/config"
	"product-service/internal/adapter/logger"
	"product-service/internal/core/domain"
	"product-service/internal/core/service/user"

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

func GetUser(ctx context.Context, userId primitive.ObjectID) (resp *user.UserResponse, cerr domain.CError) {
	log := logger.FromCtx(ctx)

	grpcConn, grpcClient, err := newUserClient(&config.GetConfig().Discovery)
	if err != nil {
		log.Error("Error creating user client", zap.Error(err))
		return nil, domain.ErrInternal
	}
	defer grpcConn.Close()

	log.Info("Created new grpc user client")
	log.Info("Making request to fetch user", zap.String("user_id", userId.Hex()))

	resp, err = grpcClient.Get(context.Background(), &user.UserRequest{UserId: userId.Hex()})
	if err != nil {
		log.Error("Error fetching user", zap.Error(err))
		return nil, domain.ErrInternal
	}

	log.Info("Successfully fetched user")
	return resp, nil
}
