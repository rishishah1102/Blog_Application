package main

import (
	"log"

	"blog-application/config"
	"blog-application/server"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize server
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	// Start server
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
