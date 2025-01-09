package repository

import (
	"context"

	"product-service/internal/adapter/storage/mongodb"
	"product-service/internal/core/domain"
)

/**
 * CategoryRepository implements port.CategoryRepository interface
 * and provides an access to the postgres database
 */
type PingRepository struct {
	db *mongodb.DB
}

// NewCategoryRepository creates a new category repository instance
func NewPingRepository(db *mongodb.DB) *PingRepository {
	return &PingRepository{
		db,
	}
}

func (pr *PingRepository) CreatePing(ctx context.Context, category *domain.Ping) error {
	return nil
}
