# Transformer

Transformer is a Go library that enables conversion between different data formats (JSON, XML, YAML). Using a common data structure, you can perform lossless conversions between these formats.

[Turkish Documentation (Türkçe Dokümantasyon)](README_TR.md)

## Features

- Convert between JSON, XML, and YAML formats
- Consistent conversion with common data structure
- Easy to use
- Type safety
- Customizable conversion rules
- High test coverage
- Thread-safe operations
- Minimal dependencies

## Requirements

- Go 1.23 or higher
- Dependencies:
  - `gopkg.in/yaml.v3` for YAML operations

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
yamlData, err := tyaml.NodeToYaml(node) // returns []byte
if err != nil {
    log.Fatal(err)
}
```

### Cross-Format Conversions (one-call API)

```go
import transformer "github.com/mstgnz/transformer"

jsonData := []byte(`{"name": "John", "age": 30}`)

// JSON -> XML
xmlData, err := transformer.JsonToXml(jsonData)

// JSON -> YAML
yamlData, err := transformer.JsonToYaml(jsonData)

// XML -> JSON
xmlInput := []byte(`<root><name>John</name></root>`)
jsonOut, err := transformer.XmlToJson(xmlInput)

// XML -> YAML
yamlOut, err := transformer.XmlToYaml(xmlInput)

// YAML -> JSON
yamlInput := []byte("name: John\nage: 30\n")
jsonOut, err = transformer.YamlToJson(yamlInput)

// YAML -> XML
xmlOut, err := transformer.YamlToXml(yamlInput)
```

### Cross-Format Conversions (via Node)

For more control — inspecting or modifying values between steps — use the Node API directly:

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
yamlData := []byte("name: John\nage: 30")
node, _ := tyaml.DecodeYaml(yamlData)
jsonData, _ := tjson.NodeToJson(node)
```

## Format-Specific Behavior

### XML Namespaces

Namespace prefixes are preserved end-to-end. An element `soap:Envelope` is stored in the Node with key `"soap:Envelope"`, and namespace declarations are stored as `@xmlns:soap` attributes. Round-tripping a namespaced XML document produces identical output.

```xml
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body/>
</soap:Envelope>
```

decodes to a Node whose key is `"soap:Envelope"` with an `@xmlns:soap` attribute, then encodes back to the same XML.

### JSON Large Integers

JSON numbers are decoded using `json.Decoder.UseNumber()`, which preserves large integers exactly. Numbers larger than 2^53 that cannot be represented as `float64` without precision loss are kept as strings internally and written back as-is to JSON output.

```go
// 9007199254740993 is preserved exactly - no float64 rounding
node, _ := tjson.DecodeJson([]byte(`{"id": 9007199254740993}`))
out, _ := tjson.NodeToJson(node)
// out == {"id":9007199254740993}
```

### YAML Anchors and Aliases

YAML anchors (`&anchor`) and aliases (`*alias`) are resolved automatically by the YAML parser before the data reaches the Node structure. The resolved values are preserved correctly, but the anchor/alias syntax itself is not retained on re-encoding.

```yaml
defaults: &defaults
  timeout: 30

production:
  <<: *defaults   # merged: production.timeout == 30 after decode
  host: prod.example.com
```

### Nesting Depth Limit

All decoders enforce a maximum nesting depth of 100 levels and return an error if exceeded. This prevents stack overflows on pathological inputs.

## Package Structure

- `node`: Contains core data structure and operations
  - Node structure for representing hierarchical data
  - Value types and type conversion operations
  - Tree traversal and manipulation functions
- `tjson`: Handles JSON conversion operations
  - JSON encoding/decoding
  - JSON validation
  - JSON file operations
- `txml`: Handles XML conversion operations
  - XML encoding/decoding
  - XML validation
  - XML file operations
  - XML attribute handling
- `tyaml`: Handles YAML conversion operations
  - YAML encoding/decoding
  - YAML validation
  - YAML file operations
- `example`: Contains example usages
  - Basic conversion examples
  - Complex data structure examples
  - Error handling examples

## Data Types

The Node structure supports the following data types:

- `TypeNull`: Null value
- `TypeObject`: Object (key-value pairs)
  - Supports nested objects
  - Maintains key order
  - Handles circular references
- `TypeArray`: Array
  - Supports mixed types
  - Preserves order
- `TypeString`: String
- `TypeNumber`: Number (integers and floating-point)
- `TypeBoolean`: Boolean

## Error Handling

The library provides detailed error information for various scenarios:

- File operations errors
- Format validation errors
- Conversion errors
- Type mismatch errors
- Structure validation errors

Example error handling:

```go
if err := validateAndConvert(); err != nil {
    switch e := err.(type) {
    case *FormatError:
        log.Printf("Invalid format: %v", e)
    case *ConversionError:
        log.Printf("Conversion failed: %v", e)
    default:
        log.Printf("Unexpected error: %v", e)
    }
}
```

## Testing

The library has comprehensive test coverage. You can use the following make commands to run tests:

### General Test Commands
```bash
# Run all tests
make test

# Run all tests with verbose output
make test-verbose

# Run tests with coverage
make test-cover

# Generate HTML coverage report
make test-coverage-report
```

### Package Specific Tests
```bash
# Run JSON tests
make test-json

# Run XML tests
make test-xml

# Run YAML tests
make test-yaml

# Run Node tests
make test-node

# Run Benchmark tests
make test-bench
```

### Package Specific Coverage Reports
```bash
# Run JSON tests with coverage
make test-json-cover

# Run XML tests with coverage
make test-xml-cover

# Run YAML tests with coverage
make test-yaml-cover

# Run Node tests with coverage
make test-node-cover
```

Current test coverage: >90%

## Performance

The library is optimized for:
- Memory efficiency
- CPU usage
- Large file handling
- Concurrent operations

### Benchmark Results

```bash
goos: darwin
goarch: arm64
cpu: Apple M1
BenchmarkJSONMarshal-8           4416622               261.0 ns/op           192 B/op          2 allocs/op
BenchmarkXMLMarshal-8             975189              1230 ns/op            4704 B/op         10 allocs/op
BenchmarkYAMLMarshal-8            213493              5284 ns/op           16728 B/op         47 allocs/op
BenchmarkJSONUnmarshal-8         1000000              1742 ns/op             272 B/op          9 allocs/op
BenchmarkXMLUnmarshal-8           370683              3104 ns/op            2328 B/op         56 allocs/op
BenchmarkYAMLUnmarshal-8          142972              8640 ns/op           10128 B/op        108 allocs/op
BenchmarkLargeJSONMarshal-8        66734             17580 ns/op           10953 B/op          2 allocs/op
BenchmarkLargeXMLMarshal-8         12298             97192 ns/op           33456 B/op         15 allocs/op
BenchmarkLargeYAMLMarshal-8         2500            466568 ns/op         1581555 B/op       3149 allocs/op
```

#### Analysis
- **JSON** shows the best performance in both marshaling and unmarshaling operations
  - Marshal: ~261 ns/op with only 2 allocations
  - Unmarshal: ~1.7 µs/op with 9 allocations
- **XML** performs slower than JSON
  - Marshal: ~1.2 µs/op with 10 allocations
  - Unmarshal: ~3.1 µs/op with 56 allocations
- **YAML** shows the highest resource usage
  - Marshal: ~5.2 µs/op with 47 allocations
  - Unmarshal: ~8.6 µs/op with 108 allocations
- For large data operations:
  - JSON maintains efficiency with minimal allocations
  - XML shows moderate performance degradation
  - YAML shows significant increase in both time and memory usage

## Security

- Input validation to prevent XML entity attacks
- Memory limit checks for large files
- Safe type conversions
- No external command execution

## Contributing

This project is open-source, and contributions are welcome. Feel free to contribute or provide feedback of any kind.


## License

This project is licensed under the Apache License, Version 2.0. See the [LICENSE](LICENSE) file for details.