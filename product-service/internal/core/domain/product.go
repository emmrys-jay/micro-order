package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductStatus string

const (
	ProductStatusActive     ProductStatus = "active"
	ProductStatusInactive   ProductStatus = "inactive"
	ProductStatusOutOfStock ProductStatus = "out_of_stock"
)

var StringToProductStatus = map[string]ProductStatus{
	"active":       ProductStatusActive,
	"inactive":     ProductStatusInactive,
	"out_of_stock": ProductStatusOutOfStock,
}

func (ur ProductStatus) String() string {
	return string(ur)
}

type Product struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Price       float64            `json:"price" bson:"price"`
	Quantity    int32              `json:"quantity" bson:"quantity"`
	Status      ProductStatus      `json:"status" bson:"status"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type CreateProductRequest struct {
	Name        string        `json:"name" validate:"required"`
	Description string        `json:"description" validate:"required"`
	Price       float64       `json:"price" validate:"required,gte=0"`
	Quantity    int32         `json:"quantity" validate:"required,gte=1"`
	Status      ProductStatus `json:"-"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int32   `json:"quantity"`
	Status      string  `json:"status"`
}

// User roles
const (
	RAdmin string = "admin"
	RUser  string = "user"
)
