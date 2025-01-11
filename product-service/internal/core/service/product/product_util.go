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
