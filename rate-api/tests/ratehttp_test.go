package rateservice_test

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"

    "github.com/tabintel/rest-api/rate-api/internal/rateservice"
)

func TestRateHTTPHandler(t *testing.T) {
    // Create a new rate service
    service := &mockRateService{}

    // Create a new HTTP request
    req := httptest.NewRequest("POST", "/v1/rate", bytes.NewBufferString(`{"currency_pair": "USD-EUR"}`))

    // Create a new HTTP recorder
    rr := httptest.NewRecorder()

    // Create a new HTTP handler
    handler := rateservice.NewHTTPHandler(service)

    // Serve HTTP
    handler.ServeHTTP(rr, req)

    // Check the status code
    assert.Equal(t, http.StatusOK, rr.Code)

    // Check the response body
    expectedResponseBody := `{"currency_pair":"USD-EUR","rate":0.88}`
    assert.Equal(t, expectedResponseBody, rr.Body.String())
}

// MockRateService is a mock implementation of the rate service
type mockRateService struct{}

// GetRate returns a mocked exchange rate
func (m *mockRateService) GetRate(currencyPair string) (float64, error) {
    // Mocked implementation, return 0.88 for USD-EUR
    return 0.88, nil
}
