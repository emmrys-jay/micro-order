package port

import (
	"context"

	"owner-service/internal/core/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRepository is an interface for interacting with User-related data
type UserRepository interface {
	// CreateUser inserts a new user into the database
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, domain.CError)
	// GetUserByID fetches a new user from the database using it's id
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*domain.User, domain.CError)
	// GetUserByEmail fetches a new user from the database using it's email
	GetUserByEmail(ctx context.Context, email string) (*domain.User, domain.CError)
	// ListUsers fetches and returns all users in the database
	ListUsers(ctx context.Context) ([]domain.User, domain.CError)
	// UpdateUser updates a user in the database and returns the updated user
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, domain.CError)
	// DeleteUser performs a soft delete on a user specified by its id
	DeleteUser(ctx context.Context, id primitive.ObjectID) domain.CError
}

// UserService is an interface for interacting with User-related business logic
type UserService interface {
	// RegisterUser is used to register a new user. It returns the new user after saving it
	RegisterUser(ctx context.Context, user *domain.CreateUserRequest) (*domain.User, domain.CError)
	// GetUser returns a user specified by its id
	GetUser(ctx context.Context, id primitive.ObjectID) (*domain.User, domain.CError)
	// ListUsers returns all users in the system
	ListUsers(ctx context.Context) ([]domain.User, domain.CError)
	// UpdateUser updates a user with the specified details and returns the updated user
	UpdateUser(ctx context.Context, id primitive.ObjectID, user *domain.UpdateUserRequest) (*domain.User, domain.CError)
	// DeleteUser deletes a user specified by id
	DeleteUser(ctx context.Context, id primitive.ObjectID) domain.CError
	// CreateAdminUser is an admin-only function used to create an admin user
	CreateAdminUser(ctx context.Context, email, password string) domain.CError
}
