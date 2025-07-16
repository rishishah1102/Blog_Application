package server

import (
	"blog-application/config"
	"blog-application/database"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mockClient struct {
	connectErr error
	pingErr    error
}

func (m *mockClient) Connect(ctx context.Context) error {
	return m.connectErr
}

func (m *mockClient) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	return m.pingErr
}

func (m *mockClient) Database(name string, opts ...*options.DatabaseOptions) *mongo.Database {
	return nil
}

func (m *mockClient) Disconnect(ctx context.Context) error {
	return nil
}

func TestNewServer(t *testing.T) {
	// Save original newClient function and restore after test
	originalNewClient := database.NewClient
	defer func() { database.NewClient = originalNewClient }()

	tests := []struct {
		name        string
		cfg         *config.Config
		mockClient  database.MongoClientInterface
		mockErr     error
		expectError bool
	}{
		{
			name: "successful server creation",
			cfg: &config.Config{
				MongoDB: config.MongoDBConfig{
					MongoURI:     "mongodb://testhost:27017",
					DatabaseName: "testdb",
					Timeout:      10,
				},
			},
			mockClient: &mockClient{
				connectErr: nil,
				pingErr:    nil,
			},
			mockErr:     nil,
			expectError: false,
		},
		{
			name: "MongoDB connection failure",
			cfg: &config.Config{
				MongoDB: config.MongoDBConfig{
					MongoURI:     "invalid-uri",
					DatabaseName: "testdb",
					Timeout:      10,
				},
			},
			mockClient:  &mockClient{},
			mockErr:     errors.New("connection failed"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			database.NewClient = func(ctx context.Context, uri string) (database.MongoClientInterface, error) {
				return tt.mockClient, tt.mockErr
			}

			server, err := NewServer(tt.cfg)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, server)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, server)
				assert.NotNil(t, server.logger)
				assert.NotNil(t, server.Router)
				assert.Equal(t, tt.cfg, server.cfg)
			}
		})
	}
}
