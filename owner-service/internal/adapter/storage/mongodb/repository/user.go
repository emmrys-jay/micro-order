package repository

import (
	"context"
	"time"

	"owner-service/internal/adapter/config"
	"owner-service/internal/adapter/storage/mongodb"
	"owner-service/internal/core/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/**
 * UserRepository implements port.UserRepository interface
 * and provides an access to the mongo database
 */
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *mongodb.DB) *UserRepository {
	return &UserRepository{
		collection: db.Client.Database(config.GetConfig().Database.Name).Collection("users"),
	}
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, domain.CError) {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := ur.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, domain.NewInternalCError(err.Error())
	}

	return user, nil
}

// GetUserByID gets a user by ID from the database
func (ur *UserRepository) GetUserByID(ctx context.Context, id primitive.ObjectID) (*domain.User, domain.CError) {
	var user domain.User

	err := ur.collection.FindOne(ctx, bson.M{"_id": id, "deleted_at": bson.M{"$exists": false}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrDataNotFound
		}
		return nil, domain.NewInternalCError(err.Error())
	}

	return &user, nil
}

// GetUserByEmail gets a user by email from the database
func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, domain.CError) {
	var user domain.User

	err := ur.collection.FindOne(ctx, bson.M{"email": email, "deleted_at": bson.M{"$exists": false}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrDataNotFound
		}
		return nil, domain.NewInternalCError(err.Error())
	}

	return &user, nil
}

// ListUsers lists all users from the database
func (ur *UserRepository) ListUsers(ctx context.Context) ([]domain.User, domain.CError) {
	var users []domain.User

	cursor, err := ur.collection.Find(ctx, bson.M{"deleted_at": bson.M{"$exists": false}}, options.Find().SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, domain.NewInternalCError(err.Error())
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			return nil, domain.NewInternalCError(err.Error())
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, domain.NewInternalCError(err.Error())
	}

	return users, nil
}

// UpdateUser updates a user by ID in the database
func (ur *UserRepository) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, domain.CError) {
	user.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"updated_at": user.UpdatedAt,
		},
	}

	result := ur.collection.FindOneAndUpdate(ctx, bson.M{"_id": user.ID, "deleted_at": bson.M{"$exists": false}}, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, domain.ErrDataNotFound
		}
		return nil, domain.NewInternalCError(result.Err().Error())
	}

	if err := result.Decode(&user); err != nil {
		return nil, domain.NewInternalCError(err.Error())
	}

	return user, nil
}

// DeleteUser deletes a user by ID from the database
func (ur *UserRepository) DeleteUser(ctx context.Context, id primitive.ObjectID) domain.CError {
	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
			"is_active":  false,
		},
	}

	_, err := ur.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return domain.NewInternalCError(err.Error())
	}

	return nil
}
