# Transformer

Transformer is a Go library that enables conversion between different data formats (JSON, XML, YAML). A document is decoded into one common tree, and every writer reads that tree, so key order, value types and repeated elements survive the trip. What each format cannot express is listed under [Format-Specific Behavior](#format-specific-behavior).

[Turkish Documentation (Türkçe Dokümantasyon)](README_TR.md)

## Features

- Convert between JSON, XML, and YAML formats
- Key order, value types and repeated elements survive every conversion
- Search a document for a key and get back where each hit sits, with its path and depth
- Find and replace: change every match at once, or exactly the one a path names
- Address any node or value by path, array indices included
- Numbers are never routed through float64, so long identifiers keep every digit
- High test coverage
- One dependency: `gopkg.in/yaml.v3`

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

## Searching a Document

Conversion is one half of what the Node model is for. The other is looking things up: a key can appear in twenty places in a large document, and which one is meant is a question the route answers.

```go
matches, _ := transformer.SearchJson(data, "name")
for _, m := range matches {
    fmt.Println(m)
}
// root.meta.name (depth 2) = catalog
// root.groups[0].members[0].name (depth 5) = alice
// root.groups[0].members[0].profile.name (depth 6) = Alice Doe
// root.deep[0][0].items[0].entries[0].detail.name (depth 9) = deepest
```

`SearchJson`, `SearchXml` and `SearchYaml` decode and search in one call. On a tree you already hold, the same is `n.Search(key)`.

Each `Match` carries three things:

- `Node` — the node that matched
- `Path` — how to address it, opening with the searched node's own key: `root.groups[0].members[1].profile.name`
- `Depth` — how many steps that route takes. A direct child is at depth 1, and an array index is a step of its own.

### Which levels does a key sit on

```go
root.Depths("name")   // [1 2 3 4 5 6 7 8 9 10]
```

`Depths` reports the levels a key appears on, in ascending order and without repeats, so a key that shows up at one level can be told from one that runs through the whole document.

### Reading a value back

A path from a match reads back with `GetNodeByPath`, and `GetValueByPath` also reaches plain values inside a list, which have no node of their own:

```go
root.GetNodeByPath("root.groups[0].members[1].profile.name")   // *Node
root.GetValueByPath("root.matrix[1][2]")                       // *Value, true
root.GetValueByPath("root.cube[1][0][1]")                      // *Value, true
```

A key that itself contains a dot or a bracket cannot be addressed this way; use the `Node` the match carries.

### Searching by anything else

```go
// every list in the document
lists := root.SearchFunc(func(n *node.Node) bool { return n.Type() == node.TypeArray })

// does this key exist anywhere, including inside lists
root.Exists("milestones")   // true
```

Keys inside arrays, and inside arrays of arrays, are reached like any other. A plain value inside a list has no node of its own, so it cannot be a search result; search its list and read the values from it.


## Changing a Document

Finding something is usually the first half of a job. `EditJson`, `EditXml` and `EditYaml` decode, hand you the tree, and write the result back in the same format:

```go
out, _ := transformer.EditJson(data, func(n *node.Node) {
    n.ReplaceText("staging.example.com", "example.com")
    n.Rename("labels", "tags")
    n.Remove("debug")
})
```

### Changing everything that matches

| | |
| --- | --- |
| `Replace(key, value)` | give every node with that key a copy of the value |
| `ReplaceFunc(predicate, replace)` | hand every node the predicate accepts to `replace` |
| `Rename(oldKey, newKey)` | rename every key that reads `oldKey` |
| `Remove(key)` | delete every object member with that key |
| `ReplaceText(old, new)` | rewrite text inside every string value, lists included |

Each returns how many nodes it changed, so a replacement that matched nothing is not silent.

Every target gets its own copy of the value, so changing one result afterwards does not change the rest. `ReplaceText` keeps to text: a key, a number or a boolean that happens to contain the same characters is left alone, and `Rename` is what changes keys.

### Changing exactly one

A path from a search addresses a single node, which is how you change the one occurrence you meant rather than all of them:

```go
matches := root.Search("name")
root.SetByPath(matches[3].Path, &node.Value{Type: node.TypeString, Worth: "renamed"})
```

| | |
| --- | --- |
| `SetByPath(path, value)` | replace the value a path addresses, or add the member when the last segment names one that is not there |
| `RemoveByPath(path)` | remove an object member, or an element from a list |

Both report whether the path led anywhere. `SetByPath` also reaches list elements, including inside lists of lists (`root.matrix[1][2]`), and wraps an object handed to a list element the way the decoders wrap one.

A tree stays sound through all of this: `Validate` still passes after any of them, and the document can be written back in any of the three formats.


## Format-Specific Behavior

### Key Order

Keys keep the order the source document wrote them in, in every direction. Neither decoding nor encoding goes through a Go map, because map iteration is randomized per run and would reshuffle the document differently on every conversion.

```go
out, _ := transformer.JsonToXml([]byte(`{"zebra":1,"apple":2}`))
// <root><zebra>1</zebra><apple>2</apple></root>, on every run
```

### XML Repeated Elements

Sibling elements sharing a name are collected into one array node, and an array node is written back as one element per value. This is what lets a list survive XML, which has no syntax for one.

```go
out, _ := transformer.XmlToJson([]byte(`<root><item>1</item><item>2</item></root>`))
// {"item":[1,2]}
```

### XML Attributes and Text

Attributes are stored as `@name` child nodes. An element carrying both attributes and text keeps its text in a `#text` child, so neither displaces the other. An element carrying only text stays a plain string, number or boolean.

```go
out, _ := transformer.XmlToJson([]byte(`<root id="5"><name lang="en">John</name></root>`))
// {"@id":"5","name":{"@lang":"en","#text":"John"}}
```

The same shape converts back to the original XML.

### XML Type Markers

A few values JSON and YAML can write have no shape of their own in XML. An empty tag could be a null, an empty string or an empty object, and a single `<port>` element could be one value or a list holding one. Where the shape would otherwise be lost, the writer marks the element in this library's own namespace and the reader reads it back:

```xml
<n t:type="null"/>                          <!-- null -->
<s t:type="string"/>                        <!-- the empty string -->
<a t:type="array"/>                         <!-- the empty list -->
<a t:type="array"><item>80</item></a>       <!-- a list holding one value -->
```

The marker is written only where it is needed, and the namespace is declared on the root element only when the document uses one. A list of two or more plain values still travels as plain repeated elements, so ordinary documents come out exactly as they did before markers existed.

```go
out, _ := transformer.JsonToXml([]byte(`{"ports":[80,443]}`))
// <root><ports>80</ports><ports>443</ports></root>   no marker needed

out, _ = transformer.JsonToXml([]byte(`{"ports":[80]}`))
// <root xmlns:t="..."><ports t:type="array"><item>80</item></ports></root>
```

### XML Element Names

XML names are not free text: an element cannot be called `first name` or `1st`. A key that is not a legal name is escaped on the way out and read back on the way in, so the key survives and the document stays valid. Names that are already legal, including non-ASCII ones, are written untouched.

```go
out, _ := transformer.JsonToXml([]byte(`{"first name":"John","şehir":"İstanbul"}`))
// <root><first_x0020_name>John</first_x0020_name><şehir>İstanbul</şehir></root>
```

A colon is treated as a namespace prefix only when the document declares that prefix. `soap:Envelope` beside an `xmlns:soap` declaration stays a prefixed name; a key that merely contains a colon has it escaped, because writing it as a prefix would name a namespace that was never declared.

### XML Namespaces

Namespace prefixes are preserved end-to-end. An element `soap:Envelope` is stored in the Node with key `"soap:Envelope"`, and namespace declarations are stored as `@xmlns:soap` attributes. Round-tripping a namespaced XML document produces identical output.

```xml
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body/>
</soap:Envelope>
```

decodes to a Node whose key is `"soap:Envelope"` with an `@xmlns:soap` attribute, then encodes back to the same XML.

### JSON Large Integers

JSON numbers are decoded using `json.Decoder.UseNumber()` and held as the text the document wrote, never as `float64`. A 19-digit identifier therefore comes back with all 19 digits, and `1e10` comes back in the notation it arrived in.

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
  <<: *defaults # merged: production.timeout == 30 after decode
  host: prod.example.com
```

### Nesting Depth Limit

All decoders enforce a maximum nesting depth of 100 levels and return an error if exceeded. This prevents stack overflows on pathological inputs.

### What a Round Trip Does Not Preserve

- **Mixed content position.** Text interleaved with child elements is kept, but the text runs are written back before the child elements rather than between them.
- **YAML anchor syntax.** Anchors and aliases are resolved on the way in, and merge keys (`<<`) are expanded in place. The values are correct; the syntax is not written back.
- **XML comments and processing instructions.** Dropped during decode.

One thing to know about XML that this library did not write: a bare `<a/>` is read as an empty object, because without a marker there is nothing to say whether it was meant as a null, an empty string or an empty object.

Everything else survives. JSON, XML and YAML each round-trip through the Node model unchanged, and a document converted to another format and back comes out as the document it started as.

## Package Structure

- `node`: Contains core data structure and operations
  - Node structure for representing hierarchical data
  - Value types and type conversion operations
  - Tree traversal and manipulation functions
  - Search: `Search`, `SearchFunc`, `Depths`, `Exists`, `GetNodeByPath`, `GetValueByPath`
  - Edit: `Replace`, `ReplaceFunc`, `ReplaceText`, `Rename`, `Remove`, `SetByPath`, `RemoveByPath`
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
