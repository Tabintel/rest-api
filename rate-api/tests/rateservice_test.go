package rateservice

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestGetRateFromCurrencyAPI_Success(t *testing.T) {
    // Mock the HTTP client and set up expectations
    httpClient := &mockHTTPClient{
        response: []byte(`{"currency_pair": "USD-EUR", "rate": 0.88}`),
        err:      nil,
    }

    // Create a new service with the mocked HTTP client
    service := NewService(httpClient)

    // Call the function under test
    response, err := service.getRateFromCurrencyAPI(context.Background(), "USD-EUR")

    // Assert the response and error
    assert.NoError(t, err)
    assert.NotNil(t, response)
    assert.Equal(t, "USD-EUR", response.CurrencyPair)
    assert.Equal(t, 0.88, response.Rate)
}

func TestGetRateFromCurrencyAPI_Failure(t *testing.T) {
    // Mock the HTTP client to simulate an error
    httpClient := &mockHTTPClient{
        response: nil,
        err:      errors.New("error from CurrencyAPI"),
    }

    // Create a new service with the mocked HTTP client
    service := NewService(httpClient)

    // Call the function under test
    response, err := service.getRateFromCurrencyAPI(context.Background(), "USD-EUR")

    // Assert the error
    assert.Error(t, err)
    assert.Nil(t, response)
    assert.Contains(t, err.Error(), "error from CurrencyAPI")
}

// Similarly, write tests for getRateFromExchangeAPI function

// MockHTTPClient is a mock HTTP client for testing
type mockHTTPClient struct {
    response []byte
    err      error
}

// Get sends an HTTP GET request and returns an HTTP response
func (m *mockHTTPClient) Get(url string) ([]byte, error) {
    return m.response, m.err
}
