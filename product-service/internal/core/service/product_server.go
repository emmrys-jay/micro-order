package service

import (
	"context"
	"product-service/internal/core/port"
	"product-service/internal/core/service/product"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

var _ product.ProductServer = (*grpcServer)(nil)

type Config struct {
}

type grpcServer struct {
	product.UnimplementedProductServer
	config      *Config
	productRepo port.ProductRepository
}

func NewGRPCServer(config *Config, userRepo port.ProductRepository, opts ...grpc.ServerOption) (*grpc.Server, error) {

	logger := zap.L().Named("grpc_server")
	zapOpts := []grpc_zap.Option{
		grpc_zap.WithDurationField(
			func(duration time.Duration) zapcore.Field {
				return zap.Int64(
					"grpc.time_ns",
					duration.Nanoseconds(),
				)
			},
		),
	}

	opts = append(opts,
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(logger, zapOpts...),
		)),
	)

	gsrv := grpc.NewServer(opts...)
	srv, err := newgrpcServer(config, userRepo)
	if err != nil {
		return nil, err
	}

	product.RegisterProductServer(gsrv, srv)
	return gsrv, nil
}

func newgrpcServer(config *Config, productRepo port.ProductRepository) (srv *grpcServer, err error) {
	srv = &grpcServer{
		config:      config,
		productRepo: productRepo,
	}
	return srv, nil
}

func (s *grpcServer) Get(ctx context.Context, req *product.ProductRequest) (*product.ProductResponse, error) {
	logger := zap.L().Named("grpc_server")
	logger.Info("Received Get request", zap.String("product_id", req.ProductId))

	// TODO: authorize request through a mechanism

	productId, err := primitive.ObjectIDFromHex(req.ProductId)
	if err != nil {
		logger.Error("Failed to parse product ID", zap.Error(err))
		return nil, err
	}

	// Get user from repository
	retProd, err := s.productRepo.GetProductByID(ctx, productId)
	if err != nil {
		logger.Error("Failed to get product", zap.Error(err))
		return nil, err
	}

	prod := product.ProductResponse{
		Id:          retProd.ID.Hex(),
		Name:        retProd.Name,
		Description: retProd.Description,
		Price:       retProd.Price,
		Quantity:    retProd.Quantity,
		OwnerId:     retProd.OwnerID.Hex(),
		OwnerName:   retProd.OwnerName,
		OwnerPhone:  retProd.OwnerPhone,
		OwnerEmail:  retProd.OwnerEmail,
		CreatedAt:   retProd.CreatedAt.String(),
		UpdatedAt:   retProd.CreatedAt.String(),
	}

	if status, err := product.StringToProductStatus(retProd.Status.String()); err == nil {
		prod.Status = status
	}

	return &prod, nil
}

func (s *grpcServer) GetMany(ctx context.Context, req *product.ProductsRequest) (*product.ProductsResponse, error) {
	logger := zap.L().Named("grpc_server")
	logger.Info("Received GetMany request", zap.Any("request", req))

	// TODO: authorize request through a mechanism

	productIds := make([]primitive.ObjectID, 0, len(req.ProductIds))
	for _, id := range req.ProductIds {
		oId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			logger.Error("Failed to parse product ID", zap.Error(err))
			return nil, err
		}
		productIds = append(productIds, oId)
	}

	products, err := s.productRepo.GetProductsByIDs(ctx, productIds)
	if err != nil {
		logger.Error("Failed to get products", zap.Error(err))
		return nil, err
	}

	prodResp := make([]*product.ProductResponse, 0, len(products))
	for _, prod := range products {
		p := &product.ProductResponse{
			Id:          prod.ID.Hex(),
			Name:        prod.Name,
			Description: prod.Description,
			Price:       prod.Price,
			Quantity:    prod.Quantity,
			OwnerId:     prod.OwnerID.Hex(),
			OwnerName:   prod.OwnerName,
			OwnerPhone:  prod.OwnerPhone,
			OwnerEmail:  prod.OwnerEmail,
			CreatedAt:   prod.CreatedAt.String(),
			UpdatedAt:   prod.UpdatedAt.String(),
		}

		if status, err := product.StringToProductStatus(prod.Status.String()); err == nil {
			p.Status = status
		}
		prodResp = append(prodResp, p)
	}

	return &product.ProductsResponse{Products: prodResp}, nil
}
