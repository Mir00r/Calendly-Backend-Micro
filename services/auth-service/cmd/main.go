package main

import (
	"github.com/Calendly-Backend-Micro/auth-service/config"
	"github.com/Calendly-Backend-Micro/auth-service/db"
	"log"
	"os"

	"github.com/Calendly-Backend-Micro/auth-service/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load the configuration file
	configPath := "./config/config.yaml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}
	config.LoadConfig(configPath)

	// Connect to the database
	db.Connect()

	// Connect to Redis
	db.ConnectRedis()

	// Create a new Gin router
	router := gin.Default()

	// Register routes
	routes.SetupAuthRoutes(router)

	// Get the port from environment variable or default to 8081
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Start the server
	log.Printf("Starting server on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
