package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type ExchangeRateResponse struct {
	Meta struct {
		LastUpdatedAt string `json:"last_updated_at"`
	} `json:"meta"`
	Data map[string]struct {
		Code  string  `json:"code"`
		Value float64 `json:"value"`
	} `json:"data"`
}

type ExchangeRate struct {
	CurrencyPair string  `json:"currency_pair"`
	Rate         float64 `json:"rate"`
}

func main() {
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    http.HandleFunc("/exchange-rate", handleExchangeRateRequest)
    log.Fatal(http.ListenAndServe(":8080", nil))
}


func handleExchangeRateRequest(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CurrencyPair string `json:"currency_pair"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.CurrencyPair == "" {
		http.Error(w, "Currency pair is required", http.StatusBadRequest)
		return
	}

	exchangeRate, err := getExchangeRate(req.CurrencyPair)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := ExchangeRate{
		CurrencyPair: req.CurrencyPair,
		Rate:         exchangeRate,
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func getExchangeRate(currencyPair string) (float64, error) {
	var wg sync.WaitGroup
	wg.Add(2)

	var rate1, rate2 float64
	var err1, err2 error

	// Fetch exchange rate from CurrencyAPI
	go func() {
		defer wg.Done()
		rate1, err1 = fetchExchangeRate("https://api.currencyapi.com/v3/latest?apikey="+os.Getenv("CURRENCYAPI_API_KEY"), currencyPair)
	}()

	// Fetch exchange rate from ExchangeratesAPI
	go func() {
		defer wg.Done()
		rate2, err2 = fetchExchangeRate("https://api.exchangeratesapi.io/v1/latest?access_key="+os.Getenv("EXCHANGERATESAPI_IO_KEY"), currencyPair)
	}()

	wg.Wait()

	if err1 != nil && err2 != nil {
		return 0, fmt.Errorf("failed to fetch exchange rate from both services: %v, %v", err1, err2)
	}

	if err1 == nil {
		return rate1, nil
	}

	return rate2, nil
}

func fetchExchangeRate(url, currencyPair string) (float64, error) {
    // Determine the service name based on the URL
    var serviceName string
    switch {
    case strings.Contains(url, "currencyapi"):
        serviceName = "CurrencyAPI"
    case strings.Contains(url, "exchangeratesapi"):
        serviceName = "ExchangeratesAPI"
    default:
        serviceName = "Unknown"
    }

    // Log the start of the fetch process
    log.Printf("🔍 Fetching exchange rate for currency pair %s from %s", currencyPair, serviceName)

    resp, err := http.Get(url)
    if err != nil {
        return 0, fmt.Errorf("❌ Failed to get exchange rate from %s: %v", serviceName, err)
    }
    defer resp.Body.Close()

    var exchangeRateResp ExchangeRateResponse
    if err := json.NewDecoder(resp.Body).Decode(&exchangeRateResp); err != nil {
        return 0, fmt.Errorf("❌ Failed to decode response from %s: %v", serviceName, err)
    }

    // Extract currency codes from the currency pair
    currencies := strings.Split(currencyPair, "-")

    // Check if both currencies exist in the response
    for _, cur := range currencies {
        if _, ok := exchangeRateResp.Data[cur]; !ok {
            return 0, fmt.Errorf("❌ Currency code %s not found in response from %s", cur, serviceName)
        }
    }

    // Calculate the exchange rate
    rate := exchangeRateResp.Data[currencies[1]].Value / exchangeRateResp.Data[currencies[0]].Value

    // Log the successful fetch
    log.Printf("✅ Successfully fetched exchange rate for currency pair %s from %s", currencyPair, serviceName)

    return rate, nil
}

