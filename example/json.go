package main

import (
	"fmt"
	"log"

	"github.com/mstgnz/transformer/tjson"
)

// JsonExample demonstrates the usage of the tjson package
func JsonExample() {
	// Read JSON file
	data, err := tjson.ReadJson("example/files/valid.json")
	if err != nil {
		log.Fatal(err)
	}

	// Check if it's valid JSON
	if !tjson.IsJson(data) {
		log.Fatal("Invalid JSON format")
	}

	// Decode JSON to Node structure
	node, err := tjson.DecodeJson(data)
	if err != nil {
		log.Fatal(err)
	}

	// Convert Node back to JSON
	jsonData, err := tjson.NodeToJson(node)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Converted JSON:")
	fmt.Println(string(jsonData))

	// Read invalid JSON file
	data, err = tjson.ReadJson("example/files/invalid.json")
	if err != nil {
		fmt.Printf("Error reading invalid JSON: %v\n", err)
	}

	// Check if it's valid JSON
	if !tjson.IsJson(data) {
		fmt.Println("Invalid JSON format detected")
	}

	// Try to decode invalid JSON
	node, err = tjson.DecodeJson(data)
	if err != nil {
		fmt.Printf("Error decoding invalid JSON: %v\n", err)
	}
}
