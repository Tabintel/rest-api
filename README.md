### Currency Exchange Rate API

This is a Go REST API that handles HTTP requests to fetch exchange rates for currency pairs from two different APIs. It uses goroutines and a wait group to concurrently fetch the exchange rates from both APIs. The fetched rates are then compared and returned as a JSON response.

#### Process for Building the Solution

- understanding the requirements
- research on currency and exchange rate APIs
- api design
- implementation of features
- testing
- documentation

#### API Services Used

- [CurrencyAPI](https://app.currencyapi.com/): Used for fetching exchange rates for currency pairs. API key securely managed using AWS KMS.
- [ExchangeratesAPI](https://exchangeratesapi.io/): Alternative service for fetching exchange rates. Not used in the final implementation due to differences in response format.

#### Concurrency

The solution uses concurrency by concurrently fetching exchange rates from multiple external services (CurrencyAPI and ExchangeratesAPI) using Goroutines and channels.

#### Logging

Logging is also implemented in the solution to provide informative messages about the exchange rate fetching process. Log messages include details such as the start of fetching, success or failure of fetching, and the name of the API service being used.

[Thunder Client  VS Code Extension](https://www.thunderclient.com/) is used to make the `POST` request to the API with this url `http://localhost:8080/exchange-rate`

#### Code Structure:
- **`main.go`:**
  - Contains the main entry point of the program.
  - Implements concurrent fetching of exchange rates and returns the first successful response.

#### Usage:
- Run `go run main.go` to execute the program.
- Ensure that the required environment variables (`CURRENCYAPI_API_KEY` and `EXCHANGERATESAPI_IO_KEY`) are set with the appropriate API keys.

#### Code Structure:
- **`main.go`:**
  - Implements the main logic for reading the JSON data, creating `Person` objects, and performing data operations.
  - Contains functions for sorting, grouping, and filtering `Person` objects.

#### Example Usage
```go
// Send a POST request to /exchange-rate with a JSON body containing the currency pair
// e.g. {"currency_pair": "USD-EUR"}

// The server will fetch the exchange rate for the given currency pair from two different APIs
// and return the rate as a JSON response
```
#### WorkFlow

1. The program loads environment variables from a .env file.
2. It sets up an HTTP handler for the /exchange-rate endpoint.
3. When a request is received, it decodes the JSON body to extract the currency pair.
4. It checks if the currency pair is empty and returns an error if it is.
5. It calls the `getExchangeRate` function to fetch the exchange rate for the currency pair from two different APIs concurrently.
6. The `getExchangeRate` function uses goroutines and a wait group to fetch the rates from both APIs.
7. It compares the fetched rates and returns the rate from the API with no errors.
8. The fetched rate is then returned as a JSON response.