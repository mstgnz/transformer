# Transformer

Transformer is a Go library that enables conversion between different data formats (JSON, XML, YAML). Using a common data structure, you can perform lossless conversions between these formats.

[Turkish Documentation (Türkçe Dokümantasyon)](README_TR.md)

## Features

- Convert between JSON, XML, and YAML formats
- Consistent conversion with common data structure
- Easy to use
- Type safety
- Customizable conversion rules

## Installation

```bash
go get github.com/mstgnz/transformer
```

## Usage

### JSON Conversions

```go
import "github.com/mstgnz/transformer/tjson"

// Read JSON file
data, err := tjson.ReadJson("data.json")
if err != nil {
    log.Fatal(err)
}

// Validate JSON format
if !tjson.IsJson(data) {
    log.Fatal("Invalid JSON format")
}

// Convert JSON to Node structure
node, err := tjson.DecodeJson(data)
if err != nil {
    log.Fatal(err)
}

// Convert Node structure to JSON
jsonData, err := tjson.NodeToJson(node)
if err != nil {
    log.Fatal(err)
}
```

### XML Conversions

```go
import "github.com/mstgnz/transformer/txml"

// Read XML file
data, err := txml.ReadXml("data.xml")
if err != nil {
    log.Fatal(err)
}

// Validate XML format
if !txml.IsXml(data) {
    log.Fatal("Invalid XML format")
}

// Convert XML to Node structure
node, err := txml.DecodeXml(data)
if err != nil {
    log.Fatal(err)
}

// Convert Node structure to XML
xmlData, err := txml.NodeToXml(node)
if err != nil {
    log.Fatal(err)
}
```

### YAML Conversions

```go
import "github.com/mstgnz/transformer/tyaml"

// Read YAML file
data, err := tyaml.ReadYaml("data.yaml")
if err != nil {
    log.Fatal(err)
}

// Validate YAML format
if !tyaml.IsYaml(data) {
    log.Fatal("Invalid YAML format")
}

// Convert YAML to Node structure
node, err := tyaml.DecodeYaml(data)
if err != nil {
    log.Fatal(err)
}

// Convert Node structure to YAML
yamlData, err := tyaml.NodeToYaml(node)
if err != nil {
    log.Fatal(err)
}
```

### Cross-Format Conversions

```go
// JSON -> XML conversion
jsonData := []byte(`{"name": "John", "age": 30}`)
node, _ := tjson.DecodeJson(jsonData)
xmlData, _ := txml.NodeToXml(node)

// XML -> YAML conversion
xmlData := []byte(`<root><name>John</name><age>30</age></root>`)
node, _ := txml.DecodeXml(xmlData)
yamlData, _ := tyaml.NodeToYaml(node)

// YAML -> JSON conversion
yamlData := []byte(`name: John\nage: 30`)
node, _ := tyaml.DecodeYaml(yamlData)
jsonData, _ := tjson.NodeToJson(node)
```

## Package Structure

- `node`: Contains core data structure and operations
- `tjson`: Handles JSON conversion operations
- `txml`: Handles XML conversion operations
- `tyaml`: Handles YAML conversion operations
- `example`: Contains example usages

## Data Types

The Node structure supports the following data types:

- `TypeNull`: Null value
- `TypeObject`: Object (key-value pairs)
- `TypeArray`: Array
- `TypeString`: String
- `TypeNumber`: Number
- `TypeBoolean`: Boolean

## Contributing

1. Fork this repository
2. Create a new branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Create a Pull Request

## License

This project is licensed under the Apache License, Version 2.0. See the [LICENSE](LICENSE) file for details.

## Contact

Mesut GENEZ - [@mstgnz](https://github.com/mstgnz)

Project Link: [https://github.com/mstgnz/transformer](https://github.com/mstgnz/transformer)