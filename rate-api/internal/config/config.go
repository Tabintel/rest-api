package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
    Environment         string `default:"local" envconfig:"ENVIRONMENT"`
    CurrencyAPIBaseURI  string `default:"https://api.currencyapi.com" envconfig:"CURRENCY_API_BASE_URI"`
    ExchangeAPIBaseURI  string `default:"https://api.exchangeratesapi.io" envconfig:"EXCHANGE_API_BASE_URI"`
    CurrencyAPIKey      string `envconfig:"CURRENCY_API_KEY"`
    ExchangeAPIKey      string `envconfig:"EXCHANGE_API_KEY"`
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

