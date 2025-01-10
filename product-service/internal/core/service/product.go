package service

import (
	"context"
	"net/http"

	"product-service/internal/adapter/config"
	"product-service/internal/adapter/logger"
	"product-service/internal/core/domain"
	"product-service/internal/core/port"
	"product-service/internal/core/service/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

/**
 * ProductService implements port.ProductService interface
 */
type ProductService struct {
	repo  port.ProductRepository
	cache port.CacheRepository
}

// NewProductService creates a new product service instance
func NewProductService(repo port.ProductRepository, cache port.CacheRepository) *ProductService {
	return &ProductService{
		repo,
		cache,
	}
}

func (ps *ProductService) CreateProduct(ctx context.Context, prod *domain.CreateProductRequest, userID primitive.ObjectID) (*domain.Product, domain.CError) {
	prodToCreate := domain.Product{
		Name:        prod.Name,
		Description: prod.Description,
		Price:       prod.Price,
		Quantity:    prod.Quantity,
		Status:      domain.ProductStatusActive,
	}

	log := logger.FromCtx(ctx)
	grpcConn, grpcClient, err := NewUserClient(&config.GetConfig().Discovery)
	if err != nil {
		log.Error("Error creating user client", zap.Error(err))
		return nil, domain.ErrInternal
	}
	defer grpcConn.Close()

	log.Info("Created new grpc user client")
	log.Info("Making request to fetch user", zap.String("user_id", userID.Hex()))

	retUser, err := grpcClient.Get(context.Background(), &user.UserRequest{UserId: userID.Hex()})
	if err != nil {
		log.Error("Error fetching user", zap.Error(err))
		return nil, domain.ErrInternal
	}

	log.Info("Successfully fetched user")

	prodToCreate.OwnerID, _ = primitive.ObjectIDFromHex(retUser.Id)
	prodToCreate.OwnerPhone = retUser.Phone
	prodToCreate.OwnerName = retUser.FirstName + " " + retUser.LastName
	prodToCreate.OwnerEmail = retUser.Email

	prodResponse, cerr := ps.repo.CreateProduct(ctx, &prodToCreate)
	if cerr != nil {
		if cerr.Code() == 500 {

			log.Error("Error creating product", zap.Error(cerr))
			return nil, domain.ErrInternal
		}
		return nil, cerr
	}

	return prodResponse, nil
}

func (ps *ProductService) GetProduct(ctx context.Context, id primitive.ObjectID) (*domain.Product, domain.CError) {
	log := logger.FromCtx(ctx)
	product, cerr := ps.repo.GetProductByID(ctx, id)
	if cerr != nil {
		if cerr.Code() == 500 {
			log.Error("Error getting product", zap.Error(cerr))
			return nil, domain.ErrInternal
		}
		return nil, cerr
	}

	return product, nil
}

func (ps *ProductService) ListProducts(ctx context.Context) ([]domain.Product, domain.CError) {
	log := logger.FromCtx(ctx)
	users, cerr := ps.repo.ListProducts(ctx)
	if cerr != nil {
		log.Error("Error listing products", zap.Error(cerr))
		return nil, domain.ErrInternal
	}

	return users, nil
}

func (ps *ProductService) UpdateProduct(ctx context.Context, id primitive.ObjectID, req *domain.UpdateProductRequest) (*domain.Product, domain.CError) {
	log := logger.FromCtx(ctx)
	retProd, cerr := ps.GetProduct(ctx, id)
	if cerr != nil {
		return nil, cerr
	}

	if req.Name == retProd.Name && req.Description == retProd.Description && req.Status == retProd.Status.String() &&
		req.Price == retProd.Price && req.Quantity == retProd.Quantity {
		return nil, domain.NewCError(http.StatusBadRequest, "There are no changes to update")
	}

	retProd.Name = req.Name
	retProd.Description = req.Description
	retProd.Price = req.Price
	retProd.Quantity = req.Quantity

	if status, ok := domain.StringToProductStatus[req.Status]; ok {
		retProd.Status = status
	}

	userResponse, cerr := ps.repo.UpdateProduct(ctx, retProd)
	if cerr != nil {
		if cerr.Code() == 500 {
			log.Error("Error updating product", zap.Error(cerr))
			return nil, domain.ErrInternal
		}
		return nil, cerr
	}

	return userResponse, nil
}

func (ps *ProductService) DeleteProduct(ctx context.Context, id primitive.ObjectID) domain.CError {
	log := logger.FromCtx(ctx)
	cerr := ps.repo.DeleteProduct(ctx, id)
	if cerr != nil {
		if cerr.Code() == 500 {
			log.Error("Error deleting product", zap.Error(cerr))
			return domain.ErrInternal
		}
		return cerr
	}

	return nil
}
