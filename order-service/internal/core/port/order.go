package port

import (
	"context"

	"order-service/internal/core/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrderRepository is an interface for interacting with order-related data
type OrderRepository interface {
	// CreateOrder creates an order and retuns the created order
	CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, domain.CError)
	// GetOrder gets an order through its id and returns it
	GetOrder(ctx context.Context, id primitive.ObjectID) (*domain.Order, domain.CError)
	// ListOrders returns all of a user's orders
	ListOrders(ctx context.Context, userId primitive.ObjectID) ([]domain.Order, domain.CError)
	// UpdateOrder updates a single order and retuns the updated order
	UpdateOrder(ctx context.Context, order *domain.Order) (*domain.Order, domain.CError)
	// UpdateOrderProductsupdates the name of all order items with the prod
	UpdateOrderProducts(ctx context.Context, prod *domain.ProductUpdateFromQueue) (int64, domain.CError)
}

// OrderService is an interface for interacting with order-related business logic
type OrderService interface {
	// PlaceOrder places a new order using the provided details of products and user id. It returns the order
	PlaceOrder(ctx context.Context, userId primitive.ObjectID, req *domain.CreateOrderRequest) (*domain.Order, domain.CError)
	// GetOrder fetches and returns a new order by its id
	GetOrder(ctx context.Context, id primitive.ObjectID) (*domain.Order, domain.CError)
	// ListUserOrders returns all orders placed by a user
	ListUserOrders(ctx context.Context, userId primitive.ObjectID) ([]domain.Order, domain.CError)
	// UpdateOrderStatus updates the status of an order specified by the order id
	UpdateOrderStatus(ctx context.Context, orderId primitive.ObjectID, req *domain.UpdateOrderRequest) (*domain.Order, domain.CError)
	// CancelOrder cancels an order specified by its id
	CancelOrder(ctx context.Context, id primitive.ObjectID) (*domain.Order, domain.CError)
}
