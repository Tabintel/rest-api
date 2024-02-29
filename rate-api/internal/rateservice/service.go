package rateservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"errors"
	"github.com/tabintel/rest-api/rate-api/internal/config"
)

type Service struct {
	cfg *config.Config
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

func NewRate(cfg *config.Config) *Service {
	return &Service{
		cfg: cfg,
	}
}


func (s *Service) RateService(ctx context.Context, currencyPair string) (*RateResponse, error) {
    rateResponseFromCurrencyAPICh := make(chan RateResponse)
    rateResponseFromExchangeAPICh := make(chan RateResponse)

    go s.getRateFromCurrencyAPI(ctx, currencyPair, rateResponseFromCurrencyAPICh)
    go s.getRateFromExchangeAPI(ctx, currencyPair, rateResponseFromExchangeAPICh)

    var currencyAPIRate, exchangeAPIRate RateResponse
    var currencyAPIErr, exchangeAPIErr error

    select {
    case currencyAPIRate = <-rateResponseFromCurrencyAPICh:
        currencyAPIErr = currencyAPIRate.Err
    case exchangeAPIRate = <-rateResponseFromExchangeAPICh:
        exchangeAPIErr = exchangeAPIRate.Err
    }

    if currencyAPIErr == nil {
        return &currencyAPIRate, nil
    }

    if exchangeAPIErr == nil {
        return &exchangeAPIRate, nil
    }

    return nil, fmt.Errorf("failed to fetch exchange rate from both services: %v, %v", currencyAPIErr, exchangeAPIErr)
}

func (s *Service) getRateFromCurrencyAPI(ctx context.Context, currencyPair string, ch chan<- RateResponse) {
    url := fmt.Sprintf("%s/v3/rate", s.cfg.CurrencyAPIBaseURI)
    request := RateRequest{
        CurrencyPair: currencyPair,
    }

    // Construct the request payload
    requestBody, err := json.Marshal(request)
    if err != nil {
        ch <- RateResponse{
            Err: errors.Wrap(err, "encoding request body for currency API"),
        }
        return
    }

    // Make the HTTP request to CurrencyAPI
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
    if err != nil {
        ch <- RateResponse{
            Err: errors.Wrap(err, "constructing request to currency API"),
        }
        return
    }

    // Update headers for the currency API
    req.Header.Set("Content-Type", "application/json")
    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.cfg.CurrencyAPIKey))

    // Send the HTTP request
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        ch <- RateResponse{
            Err: errors.Wrap(err, "sending request to currency API"),
        }
        return
    }
    defer resp.Body.Close()

    // Decode the response from CurrencyAPI
    var rateResponse RateResponse
    err = json.NewDecoder(resp.Body).Decode(&rateResponse)
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
    url := fmt.Sprintf("%s/v3/latest?apikey=%s", s.cfg.ExchangeAPIBaseURI, s.cfg.ExchangeAPIKey)
    request := RateRequest{
        CurrencyPair: currencyPair,
    }

    // Construct the request payload
    requestBody, err := json.Marshal(request)
    if err != nil {
        ch <- RateResponse{
            Err: errors.Wrap(err, "encoding request body for exchange API"),
        }
        return
    }

    // Make the HTTP request to ExchangeAPI
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
    if err != nil {
        ch <- RateResponse{
            Err: errors.Wrap(err, "constructing request to exchange API"),
        }
        return
    }

    // Update headers for the exchange API
    req.Header.Set("Content-Type", "application/json")
    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.cfg.ExchangeAPIKey))

    // Send the HTTP request
    resp, err := http.DefaultClient.Do(req)
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



func (s *Service) getRateFromCurrencyAPI(ctx context.Context, currencyPair string, rateResponseFromCurrencyAPICh chan RateResponse, ) {
	url := fmt.Sprintf("%s/v3/rate", s.cfg.CurrencyAPIBaseURI)
	request := RateRequest{
		CurrencyPair: currencyPair,
	}

	payloadBuf := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuf).Encode(request)
	if err != nil {
		rateResponseFromCurrencyAPICh <- RateResponse{
			Err: errors.Wrap(err, "encoding payload - currency api"),
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, payloadBuf)
	if err != nil {
		rateResponseFromCurrencyAPICh <- RateResponse{
			Err: errors.Wrap(err, "constructing request to get rate - currency api"),
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.cfg.CurrencyAPIBaseURI)) //api key

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		rateResponseFromCurrencyAPICh <- RateResponse{
			Err: errors.Wrap(err, "request to get rate - currency api"),
		}
	}

	defer func() { err = resp.Body.Close() }()

	response := RateResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		rateResponseFromCurrencyAPICh <- RateResponse{
			Err: errors.Wrap(err, "decoding response to RateResponse - currency api"),
		}
	}

	rateResponseFromCurrencyAPICh <- response
	return
}

func (s *Service) getRateFromExchangeAPI(ctx context.Context, currencyPair string) (*RateResponse, error) {
	url := fmt.Sprintf("%s/v3/latest?apikey=%s", s.cfg.CurrencyAPIBaseURI, s.cfg.CurrencyAPIKEY)
	request := RateRequest{
		CurrencyPair: currencyPair,
	}

	payloadBuf := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuf).Encode(request)
	if err != nil {
		return nil, errors.Wrap(err, "encoding payload - currency api")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, payloadBuf)
	if err != nil {
		return nil, errors.Wrap(err, "constructing request to get rate - currency api")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.cfg.CurrencyAPIBaseURI)) //api key

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "request to get rate - currency api")
	}

	defer func() { err = resp.Body.Close() }()

	var response RateResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, errors.Wrap(err, "decoding response to RateResponse - currency api")
	}

	return &response, nil
}
