package main

import (
    "log"
    "net/http"

    "github.com/go-chi/chi"
    "github.com/joho/godotenv"
    "github.com/tabintel/rest-api/rate-api/internal/config"
    "github.com/tabintel/rest-api/rate-api/internal/rateservice"
)

func main() {
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Error loading config: %v", err)
    }

    // Create Rate service instance
    rateService, err := rateservice.NewRate(cfg)
    if err != nil {
        log.Fatalf("Error creating rate service: %v", err)
    }

    // Create Rate HTTP handler
    rateHttpHandler := rateservice.NewHTTPHandler(rateService)

    // Set up Chi router
    router := chi.NewRouter()
    router.Mount("/v1/rate", rateHttpHandler)

    // Start HTTP server
    log.Fatal(http.ListenAndServe(":8080", router))
}
