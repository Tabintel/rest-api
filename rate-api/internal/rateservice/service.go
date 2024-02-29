package rateservice

import (
	"bytes"
	"log"
	"io"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/tabintel/rest-api/rate-api/internal/config"
)

type Service struct {
	cfg         *config.Config
}

type RateResponse struct {
	CurrencyPair string  `json:"currency_pair"`
	Rate         float64 `json:"rate"`
	Err          error
}

type RateRequest struct {
	CurrencyPair string `json:"currency_pair"`
	Err          error
}

func NewRate(cfg *config.Config) (*Service, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}

	return &Service{
		cfg: cfg,
	}, nil
}


func (s *Service) RateService(ctx context.Context, currencyPair string) (*RateResponse, error) {
	rateResponseFromCurrencyAPICh := make(chan RateResponse)
	rateResponseFromExchangeAPICh := make(chan RateResponse)

	go s.getRateFromCurrencyAPI(ctx, currencyPair, rateResponseFromCurrencyAPICh)
	go s.getRateFromExchangeAPI(ctx, currencyPair, rateResponseFromExchangeAPICh)

	var currencyAPIRate, exchangeAPIRate RateResponse

	// Wait for both goroutines to complete
	for i := 0; i < 2; i++ {
		select {
		case currencyAPIRate = <-rateResponseFromCurrencyAPICh:
			if currencyAPIRate.Err != nil {
				return nil, errors.Wrap(currencyAPIRate.Err, "error from currency API")
			}
		case exchangeAPIRate = <-rateResponseFromExchangeAPICh:
			if exchangeAPIRate.Err != nil {
				return nil, errors.Wrap(exchangeAPIRate.Err, "error from exchange API")
			}
		}
	}

	// Check if both API calls failed
	if currencyAPIRate.Err != nil && exchangeAPIRate.Err != nil {
		return nil, errors.New("failed to fetch exchange rate from both services")
	}

	// Choose the rate from the API with a valid response
	if currencyAPIRate.Err == nil {
		return &currencyAPIRate, nil
	}
	return &exchangeAPIRate, nil
}

func (s *Service) getRateFromCurrencyAPI(ctx context.Context, currencyPair string, ch chan<- RateResponse) {
    url := "https://api.currencyapi.com/v3/latest"

    // Make the HTTP request to CurrencyAPI
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
    if err != nil {
        ch <- RateResponse{
            Err: errors.Wrap(err, "constructing request to currency API"),
        }
        return
    }

    // Set the API key as a header
    req.Header.Set("apikey", s.cfg.CurrencyAPIKey)

    // Send the HTTP request
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        ch <- RateResponse{
            Err: errors.Wrap(err, "sending request to currency API"),
        }
        return
    }
    defer resp.Body.Close()

    // Read the response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        ch <- RateResponse{
            Err: errors.Wrap(err, "reading response body from currency API"),
        }
        return
    }

    // Print the raw response body (for debugging)
    log.Printf("Raw response from CurrencyAPI: %s", body)

    // Decode the response from CurrencyAPI
    var rateResponse RateResponse
    err = json.NewDecoder(bytes.NewReader(body)).Decode(&rateResponse)
    if err != nil {
        ch <- RateResponse{
            Err: errors.Wrap(err, "decoding response from currency API"),
        }
        return
    }

    // Send the rate response to the channel
    ch <- rateResponse
}

func (s *Service) getRateFromExchangeAPI(ctx context.Context, currencyPair string, ch chan<- RateResponse) {
    url := "https://api.exchangeratesapi.io/v1/latest"

    // Construct the URL with the API key as a query parameter
    url += fmt.Sprintf("?access_key=%s", s.cfg.ExchangeAPIKey)

    // Make the HTTP request to ExchangeAPI
    resp, err := http.Get(url)
    if err != nil {
        ch <- RateResponse{
            Err: errors.Wrap(err, "sending request to exchange API"),
        }
        return
    }
    defer resp.Body.Close()

    // Decode the response from ExchangeAPI
    var rateResponse RateResponse
    err = json.NewDecoder(resp.Body).Decode(&rateResponse)
    if err != nil {
        ch <- RateResponse{
            Err: errors.Wrap(err, "decoding response from exchange API"),
        }
        return
    }

    // Send the rate response to the channel
    ch <- rateResponse
}