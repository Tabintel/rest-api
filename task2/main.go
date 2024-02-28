package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
)

// Salary represents the salary value and currency.
type Salary struct {
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

// Person represents an individual with an ID, name, and salary.
type Person struct {
	ID     string `json:"id"`
	Name   string `json:"personName"`
	Salary Salary `json:"salary"`
}

// Persons represents a collection of Person records.
type Persons struct {
	Data []Person `json:"data"`
}

// ExchangeRateAPI defines the interface for fetching exchange rates.
type ExchangeRateAPI interface {
	GetExchangeRate(currency string) (float64, error)
}

// CurrencyAPI implements the ExchangeRateAPI interface.
type CurrencyAPI struct{}

// GetExchangeRate fetches the exchange rate for a given currency.
func (c *CurrencyAPI) GetExchangeRate(currency string) (float64, error) {
	// Simulate fetching exchange rate from CurrencyAPI
	// For demonstration purposes, assume a fixed exchange rate for different currencies
	exchangeRates := map[string]float64{
		"USD": 1.0, // USD to USD rate is 1.0
		"EUR": 0.85,
		"GBP": 0.75,
		// Add more exchange rates as needed
	}
	rate, ok := exchangeRates[currency]
	if !ok {
		return 0, fmt.Errorf("exchange rate for currency %s not found", currency)
	}
	return rate, nil
}

// SortBySalaryAscending sorts the Persons data by salary in ascending order.
func (p *Persons) SortBySalaryAscending() {
	sort.Slice(p.Data, func(i, j int) bool {
		return p.Data[i].Salary.Value < p.Data[j].Salary.Value
	})
}

// SortBySalaryDescending sorts the Persons data by salary in descending order.
func (p *Persons) SortBySalaryDescending() {
	sort.Slice(p.Data, func(i, j int) bool {
		return p.Data[i].Salary.Value > p.Data[j].Salary.Value
	})
}

// GroupByCurrency groups the Persons data by salary currency.
func (p *Persons) GroupByCurrency() map[string][]Person {
	groups := make(map[string][]Person)
	for _, person := range p.Data {
		groups[person.Salary.Currency] = append(groups[person.Salary.Currency], person)
	}
	return groups
}

// FilterByCurrency filters the Persons data by salary currency and converts salary to USD.
func (p *Persons) FilterByCurrency(api ExchangeRateAPI, currency string) ([]Person, error) {
	var filtered []Person
	for _, person := range p.Data {
		if person.Salary.Currency == currency {
			rate, err := api.GetExchangeRate(currency)
			if err != nil {
				return nil, err
			}
			person.Salary.Value *= rate // Convert salary to USD
			filtered = append(filtered, person)
		}
	}
	return filtered, nil
}

func main() {
	// Read data from JSON file
	data, err := ioutil.ReadFile("data.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	// Unmarshal JSON data into Persons struct
	var persons Persons
	if err := json.Unmarshal(data, &persons); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// Print original data
	fmt.Println("Original Data:")
	for _, p := range persons.Data {
		fmt.Printf("ID: %s, Name: %s, Salary: %.2f %s\n", p.ID, p.Name, p.Salary.Value, p.Salary.Currency)
	}

	// Sort data by salary in ascending order
	persons.SortBySalaryAscending()
	fmt.Println("\nSorted Data by Salary (Ascending):")
	for _, p := range persons.Data {
		fmt.Printf("ID: %s, Name: %s, Salary: %.2f %s\n", p.ID, p.Name, p.Salary.Value, p.Salary.Currency)
	}

	// Sort data by salary in descending order
	persons.SortBySalaryDescending()
	fmt.Println("\nSorted Data by Salary (Descending):")
	for _, p := range persons.Data {
		fmt.Printf("ID: %s, Name: %s, Salary: %.2f %s\n", p.ID, p.Name, p.Salary.Value, p.Salary.Currency)
	}

	// Group data by salary currency
	groups := persons.GroupByCurrency()
	fmt.Println("\nGrouped Data by Currency:")
	for currency, group := range groups {
		fmt.Printf("%s:\n", currency)
		for _, p := range group {
			fmt.Printf("ID: %s, Name: %s, Salary: %.2f %s\n", p.ID, p.Name, p.Salary.Value, p.Salary.Currency)
		}
	}

	// Filter data by currency and convert salary to USD
	currencyAPI := &CurrencyAPI{}
	filtered, err := persons.FilterByCurrency(currencyAPI, "EUR")
	if err != nil {
		fmt.Println("Error filtering data by currency:", err)
		return
	}
	fmt.Println("\nFiltered Data by Currency (EUR) and Converted to USD:")
	for _, p := range filtered {
		fmt.Printf("ID: %s, Name: %s, Salary: %.2f USD\n", p.ID, p.Name, p.Salary.Value)
	}
}
