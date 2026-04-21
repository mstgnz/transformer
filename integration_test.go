package transformer_test

import (
	"strings"
	"testing"

	transformer "github.com/mstgnz/transformer"
	"github.com/mstgnz/transformer/tjson"
	"github.com/mstgnz/transformer/txml"
	"github.com/mstgnz/transformer/tyaml"
)

// ---- Convenience functions ----

func TestJsonToXml(t *testing.T) {
	out, err := transformer.JsonToXml([]byte(`{"city":"Istanbul","pop":15000000}`))
	if err != nil {
		t.Fatalf("JsonToXml: %v", err)
	}
	for _, want := range []string{"Istanbul", "15000000"} {
		if !strings.Contains(string(out), want) {
			t.Errorf("missing %q in XML output: %s", want, out)
		}
	}
}

func TestJsonToYaml(t *testing.T) {
	out, err := transformer.JsonToYaml([]byte(`{"key":"val","num":42}`))
	if err != nil {
		t.Fatalf("JsonToYaml: %v", err)
	}
	for _, want := range []string{"val", "42"} {
		if !strings.Contains(string(out), want) {
			t.Errorf("missing %q in YAML output: %s", want, out)
		}
	}
}

func TestXmlToJson(t *testing.T) {
	out, err := transformer.XmlToJson([]byte(`<root><name>alice</name><age>30</age></root>`))
	if err != nil {
		t.Fatalf("XmlToJson: %v", err)
	}
	for _, want := range []string{"alice", "30"} {
		if !strings.Contains(string(out), want) {
			t.Errorf("missing %q in JSON output: %s", want, out)
		}
	}
}

func TestXmlToYaml(t *testing.T) {
	out, err := transformer.XmlToYaml([]byte(`<root><host>db.local</host><port>5432</port></root>`))
	if err != nil {
		t.Fatalf("XmlToYaml: %v", err)
	}
	for _, want := range []string{"db.local", "5432"} {
		if !strings.Contains(string(out), want) {
			t.Errorf("missing %q in YAML output: %s", want, out)
		}
	}
}

func TestYamlToJson(t *testing.T) {
	out, err := transformer.YamlToJson([]byte("name: bob\nage: 25\n"))
	if err != nil {
		t.Fatalf("YamlToJson: %v", err)
	}
	for _, want := range []string{"bob", "25"} {
		if !strings.Contains(string(out), want) {
			t.Errorf("missing %q in JSON output: %s", want, out)
		}
	}
}

func TestYamlToXml(t *testing.T) {
	out, err := transformer.YamlToXml([]byte("env: production\ndebug: false\n"))
	if err != nil {
		t.Fatalf("YamlToXml: %v", err)
	}
	for _, want := range []string{"production", "false"} {
		if !strings.Contains(string(out), want) {
			t.Errorf("missing %q in XML output: %s", want, out)
		}
	}
}

func TestConvenienceFunctions_InvalidInput(t *testing.T) {
	bad := []byte(`{bad json`)
	if _, err := transformer.JsonToXml(bad); err == nil {
		t.Error("JsonToXml should fail on invalid JSON")
	}
	if _, err := transformer.JsonToYaml(bad); err == nil {
		t.Error("JsonToYaml should fail on invalid JSON")
	}

	badXML := []byte(`<unclosed>`)
	if _, err := transformer.XmlToJson(badXML); err == nil {
		t.Error("XmlToJson should fail on invalid XML")
	}
	if _, err := transformer.XmlToYaml(badXML); err == nil {
		t.Error("XmlToYaml should fail on invalid XML")
	}

	badYAML := []byte("key: : bad")
	if _, err := transformer.YamlToJson(badYAML); err == nil {
		t.Error("YamlToJson should fail on invalid YAML")
	}
	if _, err := transformer.YamlToXml(badYAML); err == nil {
		t.Error("YamlToXml should fail on invalid YAML")
	}
}

// ---- XML nested <root> element fix ----

func TestDecodeXml_NestedRootElement(t *testing.T) {
	// A <root> element nested inside another element should not hijack the document root.
	// DecodeXml always returns a wrapper node with Key "root"; the first child is <data>,
	// and <data>'s first child is the nested <root> element.
	input := []byte(`<data><root><name>alice</name></root></data>`)
	n, err := txml.DecodeXml(input)
	if err != nil {
		t.Fatalf("DecodeXml: %v", err)
	}
	// Outer wrapper has Key "root" (internal IR convention).
	if n.Key != "root" {
		t.Errorf("expected wrapper key 'root', got %q", n.Key)
	}
	// First child is the actual document root element <data>.
	if n.Value == nil || n.Value.Node == nil || n.Value.Node.Key != "data" {
		t.Error("expected 'data' as the first child of the wrapper node")
	}
	// <data>'s first child must be the nested <root>, not lost.
	dataNode := n.Value.Node
	if dataNode.Value == nil || dataNode.Value.Node == nil || dataNode.Value.Node.Key != "root" {
		t.Error("expected nested 'root' element as a child of 'data'")
	}
}

// ---- Cross-format pipeline tests ----

func TestCrossFormat_JSONToXMLToYAMLToJSON(t *testing.T) {
	original := []byte(`{"name":"alice","age":30,"active":true}`)

	n1, err := tjson.DecodeJson(original)
	if err != nil {
		t.Fatalf("DecodeJson: %v", err)
	}
	xmlBytes, err := txml.NodeToXml(n1)
	if err != nil {
		t.Fatalf("NodeToXml: %v", err)
	}
	n2, err := txml.DecodeXml(xmlBytes)
	if err != nil {
		t.Fatalf("DecodeXml: %v", err)
	}
	yamlBytes, err := tyaml.NodeToYaml(n2)
	if err != nil {
		t.Fatalf("NodeToYaml: %v", err)
	}
	n3, err := tyaml.DecodeYaml(yamlBytes)
	if err != nil {
		t.Fatalf("DecodeYaml: %v", err)
	}
	finalJSON, err := tjson.NodeToJson(n3)
	if err != nil {
		t.Fatalf("NodeToJson: %v", err)
	}
	for _, field := range []string{"alice", "30", "true"} {
		if !strings.Contains(string(finalJSON), field) {
			t.Errorf("field %q missing from final JSON: %s", field, string(finalJSON))
		}
	}
}

func TestCrossFormat_YAMLToJSONToXML(t *testing.T) {
	input := []byte("host: db.local\nport: 5432\nssl: true\n")

	n, err := tyaml.DecodeYaml(input)
	if err != nil {
		t.Fatalf("DecodeYaml: %v", err)
	}
	jsonBytes, err := tjson.NodeToJson(n)
	if err != nil {
		t.Fatalf("NodeToJson: %v", err)
	}
	n2, err := tjson.DecodeJson(jsonBytes)
	if err != nil {
		t.Fatalf("DecodeJson: %v", err)
	}
	xmlBytes, err := txml.NodeToXml(n2)
	if err != nil {
		t.Fatalf("NodeToXml: %v", err)
	}
	for _, field := range []string{"db.local", "5432", "true"} {
		if !strings.Contains(string(xmlBytes), field) {
			t.Errorf("field %q missing from XML: %s", field, string(xmlBytes))
		}
	}
}

func TestCrossFormat_XMLToYAML(t *testing.T) {
	input := []byte(`<?xml version="1.0" encoding="UTF-8"?>` +
		`<root><city>Istanbul</city><population>15000000</population></root>`)

	n, err := txml.DecodeXml(input)
	if err != nil {
		t.Fatalf("DecodeXml: %v", err)
	}
	yamlBytes, err := tyaml.NodeToYaml(n)
	if err != nil {
		t.Fatalf("NodeToYaml: %v", err)
	}
	for _, field := range []string{"Istanbul", "15000000"} {
		if !strings.Contains(string(yamlBytes), field) {
			t.Errorf("field %q missing from YAML: %s", field, string(yamlBytes))
		}
	}
}

func TestCrossFormat_LargeInteger(t *testing.T) {
	original := []byte(`{"id":9007199254740993}`)

	n1, err := tjson.DecodeJson(original)
	if err != nil {
		t.Fatalf("DecodeJson: %v", err)
	}
	xmlBytes, err := txml.NodeToXml(n1)
	if err != nil {
		t.Fatalf("NodeToXml: %v", err)
	}
	n2, err := txml.DecodeXml(xmlBytes)
	if err != nil {
		t.Fatalf("DecodeXml: %v", err)
	}
	yamlBytes, err := tyaml.NodeToYaml(n2)
	if err != nil {
		t.Fatalf("NodeToYaml: %v", err)
	}
	if !strings.Contains(string(yamlBytes), "9007199254740993") {
		t.Errorf("large integer lost in YAML output: %s", string(yamlBytes))
	}
}
