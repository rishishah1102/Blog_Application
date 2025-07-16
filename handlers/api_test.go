package handlers

import (
	"blog-application/logger"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestApi(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name: "successful api creation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewAPI(logger.NewLogger(), &mongo.Database{})
			assert.NotNil(t, api)
		})
	}
}
