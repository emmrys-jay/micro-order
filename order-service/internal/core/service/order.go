package service

import (
	"context"
	"fmt"
	"time"

	"order-service/internal/adapter/config"
	"order-service/internal/adapter/logger"
	"order-service/internal/core/domain"
	"order-service/internal/core/port"
	"order-service/internal/core/util"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

/**
 * OrderService implements port.OrderService interface
 */
type OrderService struct {
	repo     port.OrderRepository
	cache    port.CacheRepository
	cacheTtl time.Duration
}

// NewOrderService creates a new order service instance
func NewOrderService(
	repo port.OrderRepository,
	cache port.CacheRepository,
) *OrderService {

	cacheTtl, err := time.ParseDuration(config.GetConfig().Redis.Ttl)
	if err != nil {
		zap.L().Info("Error parsing cache ttl, defaulting to 24h", zap.Error(err))
		cacheTtl = 24 * time.Hour
	}

	return &OrderService{
		repo,
		cache,
		cacheTtl,
	}
}

func (os *OrderService) PlaceOrder(ctx context.Context, userId primitive.ObjectID, req *domain.CreateOrderRequest) (*domain.Order, domain.CError) {
	log := logger.FromCtx(ctx)
	_, err := os.GetUser(ctx, userId)
	if err != nil {
		log.Error("Error fetching user", zap.Error(err))
		return nil, domain.ErrInternal
	}

	var productIDMap = make(map[string]int)
	validProductIDs := make([]string, 0)
	for _, v := range req.Products {
		if _, err := primitive.ObjectIDFromHex(v.ProductID); err == nil {
			validProductIDs = append(validProductIDs, v.ProductID)
			productIDMap[v.ProductID] = v.Quantity
		}
	}

	products, err := os.GetProductsByIDs(ctx, validProductIDs)
	if err != nil {
		log.Error("Error fetching products", zap.Error(err))
		return nil, domain.ErrInternal
	}

	if len(products) == 0 {
		return nil, domain.NewBadRequestCError("none of the products specified was found")
	}

	// Check for the integrity of order quantity with quantity in stock
	// Calculate total order amount
	// Populate order items
	var totalAmount float64
	var orderItems = make([]domain.OrderItem, 0)

	for _, v := range products {
		quantityOrdered := productIDMap[v.Id]

		if int(v.Quantity)-quantityOrdered < 0 {
			errMsg := fmt.Sprintf("The quantity specified for '%s' is more than the quantity in stock: %v (specified) for %v (in stock)",
				v.Name, productIDMap[v.Id], v.Quantity)
			return nil, domain.NewBadRequestCError(errMsg)
		}

		totalAmount += (v.Price * float64(quantityOrdered))
		productId, _ := primitive.ObjectIDFromHex(v.Id)
		orderItems = append(orderItems, domain.OrderItem{
			ProductID:   productId,
			ProductName: v.Name,
			Quantity:    int32(quantityOrdered),
			UnitPrice:   v.Price,
		})
	}

	order := domain.Order{
		UserID:      userId,
		TotalAmount: totalAmount,
		OrderItems:  orderItems,
		Status:      domain.OrderStatusPending,
	}

	retOrder, cerr := os.repo.CreateOrder(ctx, &order)
	if cerr != nil {
		if cerr.Code() == 500 {
			logger.FromCtx(ctx).Error("Error placing orders", zap.Error(cerr))
			return nil, domain.ErrInternal
		}
		return nil, cerr
	}

	return retOrder, nil
}

func (os *OrderService) GetOrder(ctx context.Context, id primitive.ObjectID) (*domain.Order, domain.CError) {
	order, cerr := os.repo.GetOrder(ctx, id)
	if cerr != nil {
		if cerr.Code() == 500 {
			logger.FromCtx(ctx).Error("Error getting single order", zap.Error(cerr))
			return nil, domain.ErrInternal
		}
		return nil, cerr
	}

	return order, nil
}

func (os *OrderService) ListUserOrders(ctx context.Context, userId primitive.ObjectID) ([]domain.Order, domain.CError) {
	orders, cerr := os.repo.ListOrders(ctx, userId)
	if cerr != nil {
		if cerr.Code() == 500 {

			logger.FromCtx(ctx).Error("Error fetching user orders", zap.Error(cerr))
			return nil, domain.ErrInternal
		}
		return nil, cerr
	}

	return orders, nil
}

func (os *OrderService) UpdateOrderStatus(ctx context.Context, orderId primitive.ObjectID, req *domain.UpdateOrderRequest) (*domain.Order, domain.CError) {
	retOrder, cerr := os.GetOrder(ctx, orderId)
	if cerr != nil {
		return nil, cerr
	}

	if _, ok := domain.StringToOrderStatus[req.Status]; !ok {
		return nil, domain.NewBadRequestCError("invalid status specified: " + req.Status)
	}

	order := domain.Order{
		ID:     orderId,
		Status: domain.StringToOrderStatus[req.Status],
	}

	_, cerr = os.repo.UpdateOrder(ctx, &order)
	if cerr != nil {
		if cerr.Code() == 500 {

			logger.FromCtx(ctx).Error("Error updating order", zap.Error(cerr))
			return nil, domain.ErrInternal
		}
		return nil, cerr
	}

	retOrder.Status = order.Status // Add the updated status to the order struct to be returned
	return retOrder, nil
}

func (os *OrderService) CancelOrder(ctx context.Context, id primitive.ObjectID) (*domain.Order, domain.CError) {
	retOrder, cerr := os.GetOrder(ctx, id)
	if cerr != nil {
		return nil, cerr
	}

	if retOrder.Status != domain.OrderStatusPending {
		return nil, domain.NewBadRequestCError("You cannot cancel this order again since it has already been processed. Please contact admin")
	}

	order := domain.Order{
		ID:     id,
		Status: domain.OrderStatusCancelled,
	}

	_, cerr = os.repo.UpdateOrder(ctx, &order)
	if cerr != nil {
		if cerr.Code() == 500 {

			logger.FromCtx(ctx).Error("Error canceling order", zap.Error(cerr))
			return nil, domain.ErrInternal
		}
		return nil, cerr
	}

	retOrder.Status = order.Status // Add the updated status to the order struct to be returned
	return retOrder, nil
}

func (os *OrderService) UpdateOrdersFromQueue(log *zap.Logger, msg []byte) error {
	log.Info("Received a new message", zap.String("update", string(msg)))
	ctx := context.Background()

	var update domain.ProductUpdateFromQueue
	err := util.Deserialize(msg, &update)
	if err != nil {
		log.Error("Could not deserialize message for updating orders", zap.Error(err))
		return err
	}
	if update.NameIsUpdated {

		updatedOrderItems, err := os.repo.UpdateOrderProducts(ctx, &update)
		if err != nil {
			log.Error("Could not update products in orders", zap.Error(err))
			return err
		}

		log.Info(fmt.Sprintf("Successfully updated %v products", updatedOrderItems))
	}

	// Save to cache
	log.Info("Saving the received product to cache...")
	cacheKey := util.GenerateCacheKey("product", update.Id)

	err = os.cache.Set(ctx, cacheKey, msg, os.cacheTtl)
	if err != nil {
		log.Error("Could not save product update in cache", zap.Error(err))
	}
	log.Info("Successfully saved the product update to cache")

	return nil
}
