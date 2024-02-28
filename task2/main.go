package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
)

// Person represents a single person record
type Person struct {
	ID     string  `json:"id"`
	Name   string  `json:"personName"`
	Salary Salary  `json:"salary"`
}

// Salary represents the salary of a person
type Salary struct {
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

// Persons represents a collection of Person records
type Persons struct {
	Data []Person `json:"data"`
}

func main() {
	// Read data from JSON file
	data, err := ioutil.ReadFile("data.json")
	if err != nil {
		log.Fatalf("failed to read JSON file: %v", err)
	}

	// Unmarshal JSON data into Persons struct
	var persons Persons
	if err := json.Unmarshal(data, &persons); err != nil {
		log.Fatalf("failed to unmarshal JSON data: %v", err)
	}

	// Print original data
	fmt.Println("Original Data:")
	for _, person := range persons.Data {
		fmt.Printf("ID: %s, Name: %s, Salary: %.2f %s\n", person.ID, person.Name, person.Salary.Value, person.Salary.Currency)
	}

	// Sort data by salary in ascending order
	sort.Slice(persons.Data, func(i, j int) bool {
		return persons.Data[i].Salary.Value < persons.Data[j].Salary.Value
	})

	// Print sorted data (ascending order)
	fmt.Println("\nSorted Data by Salary (Ascending):")
	for _, person := range persons.Data {
		fmt.Printf("ID: %s, Name: %s, Salary: %.2f %s\n", person.ID, person.Name, person.Salary.Value, person.Salary.Currency)
	}

	// Sort data by salary in descending order
	sort.Slice(persons.Data, func(i, j int) bool {
		return persons.Data[i].Salary.Value > persons.Data[j].Salary.Value
	})

	// Print sorted data (descending order)
	fmt.Println("\nSorted Data by Salary (Descending):")
	for _, person := range persons.Data {
		fmt.Printf("ID: %s, Name: %s, Salary: %.2f %s\n", person.ID, person.Name, person.Salary.Value, person.Salary.Currency)
	}

	// Group data by currency
	currencyMap := make(map[string][]Person)
	for _, person := range persons.Data {
		currencyMap[person.Salary.Currency] = append(currencyMap[person.Salary.Currency], person)
	}

	// Print grouped data
	fmt.Println("\nGrouped Data by Currency:")
	for currency, people := range currencyMap {
		fmt.Printf("%s:\n", currency)
		for _, person := range people {
			fmt.Printf("ID: %s, Name: %s, Salary: %.2f %s\n", person.ID, person.Name, person.Salary.Value, person.Salary.Currency)
		}
	}

	// Filter data by currency (USD) and convert to USD
	usdCurrency := "USD"
	fmt.Printf("\nFiltered Data by Currency (%s) and Converted to USD:\n", usdCurrency)
	for _, person := range persons.Data {
		if person.Salary.Currency == usdCurrency {
			fmt.Printf("ID: %s, Name: %s, Salary: %.2f USD\n", person.ID, person.Name, person.Salary.Value)
		}
	}
}
