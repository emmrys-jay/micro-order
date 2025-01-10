package service

import (
	"fmt"
	"product-service/internal/adapter/config"
	"product-service/internal/core/service/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserClient(conf *config.DiscoveryConfiguration) (*grpc.ClientConn, user.UserClient, error) {
	conn, err := grpc.NewClient(conf.OwnerUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create a new connection: %w", err)
	}

	userClient := user.NewUserClient(conn)
	return conn, userClient, nil
}
