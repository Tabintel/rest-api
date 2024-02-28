// handlers.go

package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "your_project/internal/services/rate"
)

// APIHandler represents the API handler for currency exchange rate endpoints.
type APIHandler struct {
    exchangeRateService *rate.ExchangeRateService
}

// NewAPIHandler creates a new instance of APIHandler with the provided exchange rate service.
func NewAPIHandler(exchangeRateService *rate.ExchangeRateService) *APIHandler {
    return &APIHandler{
        exchangeRateService: exchangeRateService,
    }
}

// SetupRoutes sets up API routes with versioning.
func (h *APIHandler) SetupRoutes(router *gin.Engine) {
    v1 := router.Group("/v1")
    {
        v1.POST("/rate", h.GetExchangeRateHandler)
    }
}

// GetExchangeRateHandler handles requests to retrieve the exchange rate for a given currency pair.
func (h *APIHandler) GetExchangeRateHandler(c *gin.Context) {
    // Get currency pair from request body
    var requestBody struct {
        CurrencyPair string `json:"currency-pair" binding:"required"`
    }
    if err := c.BindJSON(&requestBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // Get exchange rate for the specified currency pair
    rate, err := h.exchangeRateService.GetExchangeRate(requestBody.CurrencyPair)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Return the exchange rate
    c.JSON(http.StatusOK, gin.H{requestBody.CurrencyPair: rate})
}
