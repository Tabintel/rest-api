package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
)

// Salary struct represents the salary of a person
type Salary struct {
	Value    interface{} `json:"value"`
	Currency string      `json:"currency"`
}

// Person struct represents a person's information
type Person struct {
	ID         string `json:"id"`
	PersonName string `json:"personName"`
	Salary     Salary `json:"salary"`
}

// Persons struct represents an array of Person objects
type Persons struct {
	Data []Person `json:"data"`
}

func main() {
	// Read JSON data from file
	data, err := ioutil.ReadFile("data.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Unmarshal JSON data into Persons struct
	var persons Persons
	err = json.Unmarshal(data, &persons)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Print original data
	fmt.Println("Original Data:")
	printPersons(persons.Data)

	// Sort data by salary (ascending)
	sort.SliceStable(persons.Data, func(i, j int) bool {
		return convertSalary(persons.Data[i].Salary) < convertSalary(persons.Data[j].Salary)
	})
	fmt.Println("\nSorted Data by Salary (Ascending):")
	printPersons(persons.Data)

	// Sort data by salary (descending)
	sort.SliceStable(persons.Data, func(i, j int) bool {
		return convertSalary(persons.Data[i].Salary) > convertSalary(persons.Data[j].Salary)
	})
	fmt.Println("\nSorted Data by Salary (Descending):")
	printPersons(persons.Data)

	// Group persons by salary currency
	fmt.Println("\nGrouped Data by Currency:")
	groupedData := groupByCurrency(persons.Data)
	printGroupedData(groupedData)

	// Filter persons by salary criteria in USD
	fmt.Println("\nFiltered Data by Currency (USD):")
	filteredData := filterByCurrency(groupedData, "USD")
	printPersons(filteredData)
}

// printPersons prints the details of persons
func printPersons(persons []Person) {
	for _, p := range persons {
		fmt.Printf("ID: %s, Name: %s, Salary: %v %s\n", p.ID, p.PersonName, p.Salary.Value, p.Salary.Currency)
	}
}

// convertSalary converts salary value to float64
func convertSalary(s Salary) float64 {
	var value float64
	switch v := s.Value.(type) {
	case float64:
		value = v
	case string:
		fmt.Sscanf(v, "%f", &value)
	}
	return value
}

// groupByCurrency groups persons by salary currency into hash maps
func groupByCurrency(persons []Person) map[string][]Person {
	groupedData := make(map[string][]Person)
	for _, p := range persons {
		groupedData[p.Salary.Currency] = append(groupedData[p.Salary.Currency], p)
	}
	return groupedData
}

// filterByCurrency filters persons by salary criteria in the specified currency
func filterByCurrency(groupedData map[string][]Person, currency string) []Person {
	return groupedData[currency]
}

// printGroupedData prints the grouped data
func printGroupedData(groupedData map[string][]Person) {
	for currency, persons := range groupedData {
		fmt.Println(currency + ":")
		printPersons(persons)
	}
}
