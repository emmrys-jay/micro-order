package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"product-service/internal/adapter/config"
	"product-service/internal/adapter/logger"
	"product-service/internal/core/domain"
	"product-service/internal/core/port"
	"product-service/internal/core/service/product"
	"product-service/internal/core/util"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

/**
 * ProductService implements port.ProductService interface
 */
type ProductService struct {
	repo     port.ProductRepository
	cache    port.CacheRepository
	producer port.MessageQueueRepository
	cacheTtl time.Duration
}

// NewProductService creates a new product service instance
func NewProductService(repo port.ProductRepository, cache port.CacheRepository, producer port.MessageQueueRepository) *ProductService {
	cacheTtl, err := time.ParseDuration(config.GetConfig().Redis.Ttl)
	if err != nil {
		zap.L().Info("Error parsing cache ttl, defaulting to 24h", zap.Error(err))
		cacheTtl = 24 * time.Hour
	}

	return &ProductService{
		repo,
		cache,
		producer,
		cacheTtl,
	}
}

func (ps *ProductService) CreateProduct(ctx context.Context, prod *domain.CreateProductRequest, userID primitive.ObjectID) (*domain.Product, domain.CError) {
	log := logger.FromCtx(ctx)
	prodToCreate := domain.Product{
		Name:        prod.Name,
		Description: prod.Description,
		Price:       prod.Price,
		Quantity:    prod.Quantity,
		Status:      domain.ProductStatusActive,
	}

	retUser, err := ps.GetUser(context.Background(), userID)
	if err != nil {
		log.Error("Error fetching user", zap.Error(err))
		return nil, domain.ErrInternal
	}

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

	var nameIsUpdated bool
	if req.Name != retProd.Name {
		nameIsUpdated = true
	}

	retProd.Name = req.Name
	retProd.Description = req.Description
	retProd.Price = req.Price
	retProd.Quantity = req.Quantity

	if status, ok := domain.StringToProductStatus[req.Status]; ok {
		retProd.Status = status
	}

	productResponse, cerr := ps.repo.UpdateProduct(ctx, retProd)
	if cerr != nil {
		if cerr.Code() == 500 {
			log.Error("Error updating product", zap.Error(cerr))
			return nil, domain.ErrInternal
		}
		return nil, cerr
	}

	// Start the process of publishing update to queue
	// Error in this process does not affect the request
	productToProduce := product.ProductUpdateForQueue{
		Id:            productResponse.ID.Hex(),
		Name:          productResponse.Name,
		Description:   productResponse.Description,
		Price:         productResponse.Price,
		Quantity:      productResponse.Quantity,
		OwnerId:       productResponse.OwnerID.Hex(),
		OwnerEmail:    productResponse.OwnerEmail,
		OwnerName:     productResponse.OwnerName,
		OwnerPhone:    productResponse.OwnerPhone,
		CreatedAt:     productResponse.CreatedAt.String(),
		UpdatedAt:     productResponse.UpdatedAt.String(),
		NameIsUpdated: nameIsUpdated,
	}

	status, _ := product.StringToProductStatus(productResponse.Status.String())
	productToProduce.Status = status

	sProd, err := util.Serialize(productToProduce)
	if err != nil {
		log.Error("Error serializing product to publish to queue", zap.Error(cerr))
		return productResponse, nil
	}

	queue := "product-updates"
	log.Info("Publishing the updated product to message queue", zap.String("queue", queue))
	correlationId := ctx.Value(domain.CorrelationIDCtxKey)

	headers := map[string]any{
		"name_is_updated":                  nameIsUpdated,
		string(domain.CorrelationIDCtxKey): correlationId,
	}

	err = ps.producer.Publish(ctx, queue, sProd, headers)
	if err != nil {
		log.Error("Could not publish product update to the queue", zap.Error(err))
		return productResponse, nil
	}

	log.Info("Successfully published message about update to queue")

	return productResponse, nil
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

func (ps *ProductService) UpdateProductsFromQueue(log *zap.Logger, msg []byte) error {
	log.Info("Received a new message", zap.String("update", string(msg)))

	var update domain.UserUpdateForQueue
	err := util.Deserialize(msg, &update)
	if err != nil {
		log.Error("Could not deserialize message for updating product owner", zap.Error(err))
		return err
	}

	updatedProducts, err := ps.repo.UpdateProductOwner(context.Background(), &update)
	if err != nil {
		log.Error("Could not update products owner", zap.Error(err))
		return err
	}

	log.Info(fmt.Sprintf("Successfully updated %v product(s)", updatedProducts))
	return nil
}
