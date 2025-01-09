package mongodb

import (
	"context"
	"fmt"
	"time"

	"owner-service/internal/adapter/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

/**
 * DB is a wrapper for MongoDB database connection
 */
type DB struct {
	Client *mongo.Client
	url    string
}

// dsn constructs the MongoDB connection string
func dsn(config *config.DatabaseConfiguration) string {
	url := fmt.Sprintf("%s://%s:%s@%s:%s/%s?authSource=admin",
		config.Protocol,
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)

	return url
}

// New creates a new MongoDB database instance
func New(ctx context.Context, config *config.DatabaseConfiguration) (*DB, error) {
	url := dsn(config)
	zap.L().Info("Connecting to the database", zap.String("url", url))

	clientOptions := options.Client().ApplyURI(url)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return &DB{
		Client: client,
		url:    url,
	}, nil
}

// Close closes the database connection
func (db *DB) Close() {
	if err := db.Client.Disconnect(context.Background()); err != nil {
		fmt.Printf("Error disconnecting from MongoDB: %v\n", err)
	}
}

// Url returns the MongoDB connection string
func (db *DB) Url() string {
	return db.url
}
