// rate.go

package rate

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
)

// ExchangeRateService represents the service for fetching exchange rates from external APIs.
type ExchangeRateService struct {
    APIKeys map[string]string
}

// NewExchangeRateService creates a new instance of ExchangeRateService with API keys loaded from environment variables.
func NewExchangeRateService() *ExchangeRateService {
    // Load API keys from environment variables
    apiKeys := make(map[string]string)
    apiKeys["exchangeratesapi.io"] = os.Getenv("EXCHANGERATESAPI_IO_KEY")
    apiKeys["openexchangerates.org"] = os.Getenv("OPENEXCHANGERATES_ORG_KEY")

    return &ExchangeRateService{
        APIKeys: apiKeys,
    }
}

// GetExchangeRate retrieves the exchange rate for the specified currency pair from two external services concurrently.
func (s *ExchangeRateService) GetExchangeRate(currencyPair string) (float64, error) {
    // Channels for receiving exchange rates from each service
    resultChannel := make(chan float64)
    errorChannel := make(chan error, 2)

    // Concurrently fetch exchange rates from two external services
    go s.fetchExchangeRate("exchangeratesapi.io", currencyPair, resultChannel, errorChannel)
    go s.fetchExchangeRate("openexchangerates.org", currencyPair, resultChannel, errorChannel)

    // Return the first exchange rate received
    select {
    case rate := <-resultChannel:
        return rate, nil
    case err := <-errorChannel:
        return 0, err
    }
}

// fetchExchangeRate fetches the exchange rate for the specified currency pair from the specified external service.
func (s *ExchangeRateService) fetchExchangeRate(service, currencyPair string, resultChannel chan<- float64, errorChannel chan<- error) {
    apiKey, ok := s.APIKeys[service]
    if !ok {
        errorChannel <- fmt.Errorf("API key not found for service: %s", service)
        return
    }

    url := fmt.Sprintf("https://%s/latest?base=%s&symbols=%s", service, currencyPair[:3], currencyPair[4:])
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        errorChannel <- err
        return
    }
    req.Header.Set("Authorization", apiKey)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        errorChannel <- err
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        errorChannel <- fmt.Errorf("failed to fetch exchange rate from %s: %s", service, resp.Status)
        return
    }

    var exchangeRateData struct {
        Rates map[string]float64 `json:"rates"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&exchangeRateData); err != nil {
        errorChannel <- err
        return
    }

    rate, ok := exchangeRateData.Rates[currencyPair[4:]]
    if !ok {
        errorChannel <- fmt.Errorf("exchange rate for currency pair %s not found in the response from %s", currencyPair, service)
        return
    }

    resultChannel <- rate
}
