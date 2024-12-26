package main

import (
	"fmt"
	"log"

	"github.com/mstgnz/transformer/txml"
)

// XmlExample demonstrates the usage of the txml package
func XmlExample() {
	// Read XML file
	data, err := txml.ReadXml("example/files/valid.xml")
	if err != nil {
		log.Fatal(err)
	}

	// Check if it's valid XML
	if !txml.IsXml(data) {
		log.Fatal("Invalid XML format")
	}

	// Decode XML to Node structure
	node, err := txml.DecodeXml(data)
	if err != nil {
		log.Fatal(err)
	}

	// Convert Node back to XML
	xmlData, err := txml.NodeToXml(node)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Converted XML:")
	fmt.Println(string(xmlData))

	// Read invalid XML file
	data, err = txml.ReadXml("example/files/invalid.xml")
	if err != nil {
		fmt.Printf("Error reading invalid XML: %v\n", err)
	}

	// Check if it's valid XML
	if !txml.IsXml(data) {
		fmt.Println("Invalid XML format detected")
	}

	// Try to decode invalid XML
	node, err = txml.DecodeXml(data)
	if err != nil {
		fmt.Printf("Error decoding invalid XML: %v\n", err)
	}
}
