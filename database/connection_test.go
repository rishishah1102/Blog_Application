package database

import (
	"blog-application/logger"
	"context"
	"errors"
	"testing"
	"time"

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

func TestNewMongoClient(t *testing.T) {
	_ = logger.NewLogger()

	// Save original newClient function and restore after test
	originalNewClient := NewClient
	defer func() { NewClient = originalNewClient }()

	tests := []struct {
		name        string
		uri         string
		mockClient  MongoClientInterface
		mockErr     error
		expectError bool
	}{
		{
			name: "successful connection",
			uri:  "mongodb://testhost:27017",
			mockClient: &mockClient{
				connectErr: nil,
				pingErr:    nil,
			},
			mockErr:     nil,
			expectError: false,
		},
		{
			name:        "connect failure - invalid uri",
			uri:         "invalid::",
			mockClient:  nil,
			mockErr:     errors.New("connection failed"),
			expectError: true,
		},
		{
			name: "ping failure",
			uri:  "mongodb://testhost:2701",
			mockClient: &mockClient{
				connectErr: nil,
				pingErr:    errors.New("ping failed"),
			},
			mockErr:     nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Override the newClient function for this test case
			NewClient = func(ctx context.Context, uri string) (MongoClientInterface, error) {
				return tt.mockClient, tt.mockErr
			}

			client, err := NewMongoClient(tt.uri, 10*time.Second)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		})
	}
}
