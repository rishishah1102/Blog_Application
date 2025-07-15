package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config is the struct for blog-application config
type Config struct {
	Server  ServerConfig
	MongoDB MongoDBConfig
}

// ServerConfig is the struct for Server
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// MongoDBConfig is the struct for MongoDB
type MongoDBConfig struct {
	MongoURI     string
	DatabaseName string
	Timeout      time.Duration
}

// LoadConfig loads the environment variables
func LoadConfig() *Config {
	godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		MongoDB: MongoDBConfig{
			MongoURI:     getEnv("MONGO_URI", "mongodb://localhost:27017"),
			DatabaseName: getEnv("DATABASE_NAME", "blog-application"),
			Timeout:      10 * time.Second,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
