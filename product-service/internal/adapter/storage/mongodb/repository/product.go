package repository

import (
	"context"
	"time"

	"product-service/internal/adapter/config"
	"product-service/internal/adapter/storage/mongodb"
	"product-service/internal/core/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/**
 * ProductRepository implements port.ProductRepository interface
 * and provides an access to the mongodb database
 */
type ProductRepository struct {
	collection *mongo.Collection
}

// NewProductRepository creates a new product repository instance
func NewProductRepository(db *mongodb.DB) *ProductRepository {
	return &ProductRepository{
		collection: db.Client.Database(config.GetConfig().Database.Name).Collection("users"),
	}
}

func (ur *ProductRepository) CreateProduct(ctx context.Context, prod *domain.Product) (*domain.Product, domain.CError) {
	prod.ID = primitive.NewObjectID()
	prod.CreatedAt = time.Now()
	prod.UpdatedAt = time.Now()

	_, err := ur.collection.InsertOne(ctx, prod)
	if err != nil {
		return nil, domain.NewInternalCError(err.Error())
	}

	return prod, nil
}

// GetProductByID gets a product by its ID from the database
func (ur *ProductRepository) GetProductByID(ctx context.Context, id primitive.ObjectID) (*domain.Product, domain.CError) {
	var prod domain.Product

	err := ur.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&prod)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrDataNotFound
		}
		return nil, domain.NewInternalCError(err.Error())
	}

	return &prod, nil
}

// ListProducts lists all products in the database
func (ur *ProductRepository) ListProducts(ctx context.Context) ([]domain.Product, domain.CError) {
	var prods = make([]domain.Product, 0)

	cursor, err := ur.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, domain.NewInternalCError(err.Error())
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var prod domain.Product
		err := cursor.Decode(&prod)
		if err != nil {
			return nil, domain.NewInternalCError(err.Error())
		}
		prods = append(prods, prod)
	}

	if err := cursor.Err(); err != nil {
		return nil, domain.NewInternalCError(err.Error())
	}

	return prods, nil
}

// UpdateProduct updates a product by ID in the database
func (ur *ProductRepository) UpdateProduct(ctx context.Context, prod *domain.Product) (*domain.Product, domain.CError) {
	prod.UpdatedAt = time.Now()

	filter := bson.M{"_id": prod.ID, "deleted_at": bson.M{"$exists": false}}
	update := bson.M{
		"$set": bson.M{
			"name":        prod.Name,
			"description": prod.Description,
			"price":       prod.Price,
			"quantity":    prod.Quantity,
			"status":      prod.Status,
			"updated_at":  prod.UpdatedAt,
		},
	}

	err := ur.collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&prod)
	if err != nil {
		return nil, domain.NewInternalCError(err.Error())
	}

	return prod, nil
}

// DeleteProduct deletes a product by ID from the database
func (ur *ProductRepository) DeleteProduct(ctx context.Context, id primitive.ObjectID) domain.CError {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
			"status":     domain.ProductStatusInactive,
		},
	}

	_, err := ur.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.ErrDataNotFound
		}
		return domain.NewInternalCError(err.Error())
	}

	return nil
}

// GetProductsByIDs gets a number of products by their ids
func (ur *ProductRepository) GetProductsByIDs(ctx context.Context, productIds []primitive.ObjectID) ([]domain.Product, domain.CError) {
	var prods = make([]domain.Product, 0)

	filter := bson.M{
		"_id":        bson.M{"$in": productIds},
		"deleted_at": bson.M{"$exists": false},
		"status":     domain.ProductStatusActive,
	}

	cursor, err := ur.collection.Find(ctx, filter, options.Find().SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, domain.NewInternalCError(err.Error())
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var prod domain.Product
		err := cursor.Decode(&prod)
		if err != nil {
			return nil, domain.NewInternalCError(err.Error())
		}
		prods = append(prods, prod)
	}

	if err := cursor.Err(); err != nil {
		return nil, domain.NewInternalCError(err.Error())
	}

	return prods, nil
}
