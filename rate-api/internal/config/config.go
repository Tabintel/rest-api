package config

import (
    "errors"

    "github.com/kelseyhightower/envconfig"
)

type Config struct {
    Environment        string `default:"local" envconfig:"ENVIRONMENT"`
    CurrencyAPIBaseURI string `default:"https://api.currencyapi.com/v3/latest" envconfig:"CURRENCYAPI_API_KEY"`
    ExchangeAPIBaseURI string `default:"https://api.exchangeratesapi.io/v1/latest" envconfig:"EXCHANGERATESAPI_IO_KEY"`
    CurrencyAPIKey     string `envconfig:"CURRENCYAPI_API_KEY"`
    ExchangeAPIKey     string `envconfig:"EXCHANGERATESAPI_IO_KEY"`
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    err := envconfig.Process("", cfg)
    if err != nil {
        return nil, err
    }

    if cfg.CurrencyAPIKey == "" || cfg.ExchangeAPIKey == "" {
        return nil, errors.New("CurrencyAPI key and ExchangeAPI key must be provided")
    }

    return cfg, nil
}
