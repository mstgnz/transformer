package main

import (
	"fmt"

	"github.com/mstgnz/transformer/node"
	"github.com/mstgnz/transformer/tjson"
	"github.com/mstgnz/transformer/txml"
	"github.com/mstgnz/transformer/tyaml"
)

func main() {
	// Create a sample JSON data
	jsonData := []byte(`{
		"name": "John Doe",
		"age": 30,
		"isStudent": false,
		"address": {
			"street": "123 Main St",
			"city": "New York"
		},
		"hobbies": ["reading", "gaming", "coding"]
	}`)

	// Parse JSON to Node
	jsonNode, err := tjson.DecodeJson(jsonData)
	if err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	fmt.Println("JSON Node:")
	jsonNode.Print()

	// Convert Node to XML
	xmlBytes, err := txml.NodeToXml(jsonNode)
	if err != nil {
		fmt.Printf("Error converting to XML: %v\n", err)
		return
	}

	fmt.Println("\nGenerated XML:")
	fmt.Println(string(xmlBytes))

	// Parse XML back to Node
	xmlNode, err := txml.DecodeXml(xmlBytes)
	if err != nil {
		fmt.Printf("Error decoding XML: %v\n", err)
		return
	}

	fmt.Println("\nXML Node:")
	xmlNode.Print()

	// Convert Node to YAML
	yamlBytes, err := tyaml.NodeToYaml(xmlNode)
	if err != nil {
		fmt.Printf("Error converting to YAML: %v\n", err)
		return
	}

	fmt.Println("\nGenerated YAML:")
	fmt.Println(string(yamlBytes))

	// Parse YAML back to Node
	yamlNode, err := tyaml.DecodeYaml(yamlBytes)
	if err != nil {
		fmt.Printf("Error decoding YAML: %v\n", err)
		return
	}

	fmt.Println("\nYAML Node:")
	yamlNode.Print()

	// Demonstrate node manipulation
	fmt.Println("\nDemonstrating node manipulation:")

	// Create a new node
	root := node.NewNode("root")

	// Add some child nodes
	person := node.NewNode("person")
	person.Value = &node.Value{
		Type: node.TypeObject,
	}

	name := node.NewNode("name")
	name.Value = &node.Value{
		Type:  node.TypeString,
		Worth: "Alice",
	}

	age := node.NewNode("age")
	age.Value = &node.Value{
		Type:  node.TypeNumber,
		Worth: "25",
	}

	// Build the tree
	root.AddToValue(&node.Value{Type: node.TypeObject})
	root.AddToStart(person)
	person.AddToStart(name)
	person.AddToEnd(age)

	fmt.Println("Created node structure:")
	root.Print()

	// Convert to different formats
	fmt.Println("\nConverted to JSON:")
	if jsonBytes, err := tjson.NodeToJson(root); err == nil {
		fmt.Println(string(jsonBytes))
	}

	fmt.Println("\nConverted to XML:")
	if xmlBytes, err := txml.NodeToXml(root); err == nil {
		fmt.Println(string(xmlBytes))
	}

	fmt.Println("\nConverted to YAML:")
	if yamlBytes, err := tyaml.NodeToYaml(root); err == nil {
		fmt.Println(string(yamlBytes))
	}
}
