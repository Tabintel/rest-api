// converter.go

package rate

import "fmt"

// CurrencyConverter represents a currency converter.
type CurrencyConverter struct {
    ExchangeRateService *ExchangeRateService
}

// NewCurrencyConverter creates a new instance of CurrencyConverter.
func NewCurrencyConverter(exchangeRateService *ExchangeRateService) *CurrencyConverter {
    return &CurrencyConverter{
        ExchangeRateService: exchangeRateService,
    }
}

// ConvertCurrency converts the specified amount from one currency to another.
func (c *CurrencyConverter) ConvertCurrency(amount float64, fromCurrency, toCurrency string) (float64, error) {
    // Get exchange rate for the specified currency pair
    exchangeRate, err := c.ExchangeRateService.GetExchangeRate(fmt.Sprintf("%s-%s", fromCurrency, toCurrency))
    if err != nil {
        return 0, err
    }

    // Convert currency
    convertedAmount := amount * exchangeRate
    return convertedAmount, nil
}
