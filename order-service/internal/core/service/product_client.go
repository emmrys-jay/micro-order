package service

import (
	"context"
	"fmt"
	"order-service/internal/adapter/config"
	"order-service/internal/adapter/logger"
	"order-service/internal/core/domain"
	"order-service/internal/core/service/product"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newProductClient(conf *config.DiscoveryConfiguration) (*grpc.ClientConn, product.ProductClient, error) {
	conn, err := grpc.NewClient(conf.ProductUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create a new connection: %w", err)
	}

	productClient := product.NewProductClient(conn)
	return conn, productClient, nil
}

func GetProduct(ctx context.Context, productId primitive.ObjectID) (exists bool, cerr domain.CError) {
	log := logger.FromCtx(ctx)
	grpcConn, grpcClient, err := newProductClient(&config.GetConfig().Discovery)
	if err != nil {
		log.Error("Error creating product client", zap.Error(err))
		return false, domain.ErrInternal
	}
	defer grpcConn.Close()

	log.Info("Created new grpc product client")
	log.Info("Making request to fetch product", zap.String("product_id", productId.Hex()))

	_, err = grpcClient.Get(context.Background(), &product.ProductRequest{ProductId: productId.Hex()})
	if err != nil {
		return false, domain.ErrInternal
	}

	log.Info("Successfully fetched user")

	return true, nil
}

func GetProductsByIDs(ctx context.Context, validProductIDs []string) (products []*product.ProductResponse, cerr domain.CError) {
	log := logger.FromCtx(ctx)
	pConn, pClient, err := newProductClient(&config.GetConfig().Discovery)
	if err != nil {
		log.Error("Error creating product client", zap.Error(err))
		return nil, domain.ErrInternal
	}
	defer pConn.Close()

	log.Info("Created new grpc product client")
	log.Info("Making request to fetch many products", zap.Any("product_ids", validProductIDs))

	prods, err := pClient.GetMany(context.Background(), &product.ProductsRequest{ProductIds: validProductIDs})
	if err != nil {
		log.Error("Error fetching products", zap.Error(err))
		return nil, domain.ErrInternal
	}

	log.Info("Successfully fetched products")

	return prods.Products, nil
}
