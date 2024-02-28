// main.go

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"your_project/internal/handler"
	"your_project/internal/services/rate"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Create exchange rate service
	exchangeRateService := rate.NewExchangeRateService()

	// Create API handler
	apiHandler := handler.NewAPIHandler(exchangeRateService)

	// Setup API routes
	apiHandler.SetupRoutes(router)

	// Start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
