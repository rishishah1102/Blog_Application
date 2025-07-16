package database

import (
	"blog-application/logger"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoClientInterface is the wrapper interface for mongo client
type MongoClientInterface interface {
	Connect(ctx context.Context) error
	Ping(ctx context.Context, rp *readpref.ReadPref) error
	Database(name string, opts ...*options.DatabaseOptions) *mongo.Database
	Disconnect(ctx context.Context) error
}

// MongoWrapper is the wrapper for mongo client
type MongoWrapper struct {
	*mongo.Client
}

// NewClient is the wrapper function of mongo.Connect
var NewClient = func(ctx context.Context, uri string) (MongoClientInterface, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &MongoWrapper{client}, nil
}

// NewMongoClient connects the application with mongoDB database and creates new mongo client
func NewMongoClient(uri string, timeout time.Duration) (MongoClientInterface, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := NewClient(ctx, uri)
	if err != nil {
		return nil, logger.WrapError(err, "failed to connect to MongoDB")
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, logger.WrapError(err, "failed to ping MongoDB")
	}

	return client, nil
}
