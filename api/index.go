package main

import (
	"log"
	"net/http"

	"blog-application/config"
	"blog-application/server"

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

// Initializes the configuration and server
func init() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize server
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	Router = srv.Router
}

// Start the backend server on vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	Router.ServeHTTP(w, r)
}
