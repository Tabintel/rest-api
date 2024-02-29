package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchExchangeRate(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a mock response
		ExchangeRateResponse := 
		exchangeRateResp := ExchangeRateResponse{
			Meta: struct {
				LastUpdatedAt string `json:"last_updated_at"`
			}{
				LastUpdatedAt: "2024-02-28T18:52:04.281Z",
			},
			Data: map[string]struct {
				Code  string  `json:"code"`
				Value float64 `json:"value"`
			}{
				"USD": {
					Code:  "USD",
					Value: 1.0,
				},
				"EUR": {
					Code:  "EUR",
					Value: 0.92,
				},
			},
		}

		// Encode the mock response as JSON
		jsonResp, err := json.Marshal(exchangeRateResp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	}))
	defer mockServer.Close()

	// Override the base URL for testing
	baseURL = mockServer.URL

	// Call the function being tested
	rate, err := fetchExchangeRate(mockServer.URL, "USD-EUR")

	// Verify the result
	if err != nil {
		t.Errorf("fetchExchangeRate() returned an error: %v", err)
	}
	if rate != 0.92 {
		t.Errorf("fetchExchangeRate() returned incorrect rate: got %f, want 0.92", rate)
	}
}
