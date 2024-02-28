### Currency Exchange Rate API

#### Overview

This project implements a REST API for retrieving exchange rates for currency pairs from multiple external services concurrently. It ensures secure API key management using AWS Key Management Service (KMS) and stores the encrypted keys in AWS Systems Manager Parameter Store.

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

The solution implements concurrency by concurrently fetching exchange rates from multiple external services (CurrencyAPI and ExchangeratesAPI) using Goroutines and channels.

#### Fulfillment of Requirements

- **REST API**: Implemented a REST API for receiving JSON requests with currency pairs and retrieving exchange rates.
- **Concurrency**: Concurrently fetches exchange rates from multiple external services.
- **Secure API Key Management**: Utilizes AWS KMS for secure storage and retrieval of API keys. API keys are encrypted before storage.
- **Documentation**: README file and inline comments provide detailed documentation of the solution and its components.

#### Logging

Logging is implemented in the solution to provide informative messages about the exchange rate fetching process. Log messages include details such as the start of fetching, success or failure of fetching, and the name of the API service being used.

### Task 2