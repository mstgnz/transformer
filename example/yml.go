package main

import (
	"fmt"
	"log"

	"github.com/mstgnz/transformer/tyaml"
)

// YmlExample demonstrates the usage of the tyaml package
func YmlExample() {
	// Read YAML file
	data, err := tyaml.ReadYaml("example/files/valid.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Check if it's valid YAML
	if !tyaml.IsYaml(data) {
		log.Fatal("Invalid YAML format")
	}

	// Decode YAML to Node structure
	node, err := tyaml.DecodeYaml(data)
	if err != nil {
		log.Fatal(err)
	}

	// Convert Node back to YAML
	yamlData, err := tyaml.NodeToYaml(node)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Converted YAML:")
	fmt.Println(string(yamlData))

	// Read invalid YAML file
	data, err = tyaml.ReadYaml("example/files/invalid.yaml")
	if err != nil {
		fmt.Printf("Error reading invalid YAML: %v\n", err)
	}

	// Check if it's valid YAML
	if !tyaml.IsYaml(data) {
		fmt.Println("Invalid YAML format detected")
	}

	// Try to decode invalid YAML
	_, err = tyaml.DecodeYaml(data)
	if err != nil {
		fmt.Printf("Error decoding invalid YAML: %v\n", err)
	}
}
