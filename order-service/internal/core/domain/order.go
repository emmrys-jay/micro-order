package domain

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "Pending"
	OrderStatusProcessing OrderStatus = "Processing"
	OrderStatusShipped    OrderStatus = "Shipped"
	OrderStatusDelivered  OrderStatus = "Delivered"
	OrderStatusCancelled  OrderStatus = "Cancelled"
)

var StringToOrderStatus = map[string]OrderStatus{
	"Pending":    OrderStatusPending,
	"Processing": OrderStatusProcessing,
	"Shipped":    OrderStatusShipped,
	"Delivered":  OrderStatusDelivered,
	"Cancelled":  OrderStatusCancelled,
}

func (e *OrderStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = OrderStatus(s)
	case string:
		*e = OrderStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for OrderStatus: %T", src)
	}
	return nil
}

type Order struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Status      OrderStatus        `json:"status" bson:"status"`
	TotalAmount float64            `json:"total_amount" bson:"total_amount"`
	OrderItems  []OrderItem        `json:"order_items,omitempty" bson:"order_items,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type OrderItem struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	OrderID     primitive.ObjectID `json:"order_id" bson:"order_id"`
	ProductID   primitive.ObjectID `json:"product_id" bson:"product_id"`
	ProductName string             `json:"product_name" bson:"product_name"`
	Quantity    int32              `json:"quantity" bson:"quantity"`
	UnitPrice   float64            `json:"unit_price" bson:"unit_price"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

type ProductInfo struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,gte=1"`
}

type CreateOrderRequest struct {
	Products []ProductInfo `json:"products"`
}

type UpdateOrderRequest struct {
	Status string `json:"status"`
}

func (ur OrderStatus) String() string {
	return string(ur)
}

// User roles
const (
	RAdmin string = "admin"
	RUser  string = "user"
)
