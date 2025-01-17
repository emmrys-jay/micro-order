package product

import "fmt"

func ProductStatusToString(ps ProductStatus) string {
	switch ps {
	case ProductStatus_PRODUCT_STATUS_ACTIVE:
		return "ACTIVE"
	case ProductStatus_PRODUCT_STATUS_INACTIVE:
		return "INACTIVE"
	case ProductStatus_PRODUCT_STATUS_OUT_OF_STOCK:
		return "OUT_OF_STOCK"
	default:
		return "UNKNOWN"
	}
}

func StringToProductStatus(status string) (ProductStatus, error) {
	switch status {
	case "ACTIVE":
		return ProductStatus_PRODUCT_STATUS_ACTIVE, nil
	case "INACTIVE":
		return ProductStatus_PRODUCT_STATUS_INACTIVE, nil
	case "OUT_OF_STOCK":
		return ProductStatus_PRODUCT_STATUS_OUT_OF_STOCK, nil
	default:
		return -1, fmt.Errorf("invalid product status: %s", status)
	}
}

type ProductUpdateForQueue struct {
	Id            string        `json:"id,omitempty"`
	Name          string        `json:"name,omitempty"`
	Description   string        `json:"description,omitempty"`
	Price         float64       `json:"price,omitempty"`
	Quantity      int32         `json:"quantity,omitempty"`
	Status        ProductStatus `json:"status,omitempty"`
	OwnerId       string        `json:"owner_id,omitempty"`
	OwnerName     string        `json:"owner_name,omitempty"`
	OwnerPhone    string        `json:"owner_phone,omitempty"`
	OwnerEmail    string        `json:"owner_email,omitempty"`
	CreatedAt     string        `json:"created_at,omitempty"`
	UpdatedAt     string        `json:"updated_at,omitempty"`
	NameIsUpdated bool          `json:"name_is_updated"`
}
