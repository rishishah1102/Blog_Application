package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected Config
	}{
		{
			name: "default values",
			expected: Config{
				Server: ServerConfig{
					Port:         "8080",
					ReadTimeout:  10 * time.Second,
					WriteTimeout: 10 * time.Second,
				},
				MongoDB: MongoDBConfig{
					MongoURI:     "mongodb://localhost:27017",
					DatabaseName: "blog-application",
					Timeout:      10 * time.Second,
				},
			},
		},
		{
			name: "custom env values",
			envVars: map[string]string{
				"SERVER_PORT":   "9090",
				"MONGO_URI":     "mongodb://remote:27017",
				"DATABASE_NAME": "testdb",
			},
			expected: Config{
				Server: ServerConfig{
					Port:         "9090",
					ReadTimeout:  10 * time.Second,
					WriteTimeout: 10 * time.Second,
				},
				MongoDB: MongoDBConfig{
					MongoURI:     "mongodb://remote:27017",
					DatabaseName: "testdb",
					Timeout:      10 * time.Second,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			cfg := LoadConfig()
			assert.Equal(t, tt.expected, *cfg)
		})
	}
}
