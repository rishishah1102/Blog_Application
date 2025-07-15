package database

import (
	"blog-application/logger"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoClient connects the application with mongoDB database and creates new mongo client
func NewMongoClient(uri string, timeout time.Duration) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, logger.WrapError(err, "failed to connect to MongoDB")
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, logger.WrapError(err, "failed to ping MongoDB")
	}

	return client, nil
}
