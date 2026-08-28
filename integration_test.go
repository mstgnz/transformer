package transformer_test

import (
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

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

// ---- Round trips across formats ----

// TestCrossFormat_JsonYamlJson checks that a document survives the trip through
// YAML unchanged. Before the writers kept key order, this came back reshuffled;
// before the YAML writer stopped reindenting by hand, it did not come back at
// all, because the YAML it produced could not be parsed.
func TestCrossFormat_JsonYamlJson(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{"flat", `{"name":"alice","age":30,"active":true}`},
		{"key order", `{"zebra":1,"apple":2,"mango":3}`},
		{"nested", `{"server":{"host":"localhost","port":5432}}`},
		{"array of scalars", `{"ports":[80,443,8080]}`},
		{"array of objects", `{"users":[{"name":"alice"},{"name":"bob"}]}`},
		{"nested arrays", `{"matrix":[[1,2],[3,4]]}`},
		{"array under a nested key", `{"spec":{"containers":[{"name":"web"},{"name":"db"}]}}`},
		{"mixed scalars", `{"n":null,"t":true,"f":false,"i":42,"d":3.14,"s":"text"}`},
		{"big integer", `{"id":9007199254740993}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			yamlOut, err := transformer.JsonToYaml([]byte(tt.json))
			if err != nil {
				t.Fatalf("JsonToYaml: %v", err)
			}
			back, err := transformer.YamlToJson(yamlOut)
			if err != nil {
				t.Fatalf("YamlToJson: %v (intermediate YAML:\n%s)", err, yamlOut)
			}
			if string(back) != tt.json {
				t.Errorf("round trip changed the document\ngot:  %s\nwant: %s\nvia YAML:\n%s", back, tt.json, yamlOut)
			}
		})
	}
}

// TestCrossFormat_JsonXmlJson checks the same trip through XML. A list needs at
// least two values to survive it: one repeated element is indistinguishable
// from a single element, which is a property of XML rather than of this code.
func TestCrossFormat_JsonXmlJson(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{"flat", `{"name":"alice","age":30,"active":true}`},
		{"key order", `{"zebra":1,"apple":2,"mango":3}`},
		{"nested", `{"server":{"host":"localhost","port":5432}}`},
		{"array of scalars", `{"ports":[80,443,8080]}`},
		{"array of objects", `{"users":[{"name":"alice"},{"name":"bob"}]}`},
		{"array under a nested key", `{"spec":{"containers":[{"name":"web"},{"name":"db"}]}}`},
		{"big integer", `{"id":9007199254740993}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xmlOut, err := transformer.JsonToXml([]byte(tt.json))
			if err != nil {
				t.Fatalf("JsonToXml: %v", err)
			}
			back, err := transformer.XmlToJson(xmlOut)
			if err != nil {
				t.Fatalf("XmlToJson: %v (intermediate XML: %s)", err, xmlOut)
			}
			if string(back) != tt.json {
				t.Errorf("round trip changed the document\ngot:  %s\nwant: %s\nvia XML: %s", back, tt.json, xmlOut)
			}
		})
	}
}

// TestCrossFormat_XmlRepeatedElementsSurvive is the end-to-end version of the
// loss that used to happen inside the JSON writer's map.
func TestCrossFormat_XmlRepeatedElementsSurvive(t *testing.T) {
	input := []byte(`<root><item>1</item><item>2</item><item>3</item></root>`)

	jsonOut, err := transformer.XmlToJson(input)
	if err != nil {
		t.Fatalf("XmlToJson: %v", err)
	}
	if want := `{"item":[1,2,3]}`; string(jsonOut) != want {
		t.Errorf("XmlToJson = %s, want %s", jsonOut, want)
	}

	yamlOut, err := transformer.XmlToYaml(input)
	if err != nil {
		t.Fatalf("XmlToYaml: %v", err)
	}
	for _, want := range []string{"- 1", "- 2", "- 3"} {
		if !strings.Contains(string(yamlOut), want) {
			t.Errorf("missing %q in YAML output:\n%s", want, yamlOut)
		}
	}
}

// TestCrossFormat_XmlAttributesSurvive covers the element that carries an
// attribute and text at once, whose attribute used to be dropped.
func TestCrossFormat_XmlAttributesSurvive(t *testing.T) {
	input := []byte(`<root id="5"><name lang="en">John</name></root>`)

	jsonOut, err := transformer.XmlToJson(input)
	if err != nil {
		t.Fatalf("XmlToJson: %v", err)
	}
	want := `{"@id":"5","name":{"@lang":"en","#text":"John"}}`
	if string(jsonOut) != want {
		t.Errorf("XmlToJson = %s, want %s", jsonOut, want)
	}

	back, err := transformer.JsonToXml(jsonOut)
	if err != nil {
		t.Fatalf("JsonToXml: %v", err)
	}
	if !strings.Contains(string(back), `<name lang="en">John</name>`) {
		t.Errorf("attributes lost on the way back to XML: %s", back)
	}
}

// TestCrossFormat_YamlXmlYaml runs a Kubernetes-shaped document, the kind this
// library exists for, through XML and back.
func TestCrossFormat_YamlXmlYaml(t *testing.T) {
	input := []byte("apiVersion: apps/v1\nkind: Deployment\nspec:\n  replicas: 2\n  ports:\n    - 80\n    - 443\n")

	xmlOut, err := transformer.YamlToXml(input)
	if err != nil {
		t.Fatalf("YamlToXml: %v", err)
	}
	jsonOut, err := transformer.XmlToJson(xmlOut)
	if err != nil {
		t.Fatalf("XmlToJson: %v", err)
	}
	want := `{"apiVersion":"apps/v1","kind":"Deployment","spec":{"replicas":2,"ports":[80,443]}}`
	if string(jsonOut) != want {
		t.Errorf("document changed on the way through XML\ngot:  %s\nwant: %s\nvia XML: %s", jsonOut, want, xmlOut)
	}
}

// TestCrossFormat_OutputIsDeterministic runs every conversion repeatedly. Go
// map iteration is randomized per run, so a single pass would not have caught
// the reshuffling this guards against.
func TestCrossFormat_OutputIsDeterministic(t *testing.T) {
	jsonInput := []byte(`{"zebra":1,"apple":{"inner":true},"mango":[1,2],"beta":"x","kilo":null}`)
	yamlInput := []byte("zebra: 1\napple:\n  inner: true\nmango:\n  - 1\n  - 2\n")
	xmlInput := []byte(`<root><zebra>1</zebra><apple><inner>true</inner></apple><mango>1</mango><mango>2</mango></root>`)

	conversions := []struct {
		name string
		run  func() ([]byte, error)
	}{
		{"JsonToXml", func() ([]byte, error) { return transformer.JsonToXml(jsonInput) }},
		{"JsonToYaml", func() ([]byte, error) { return transformer.JsonToYaml(jsonInput) }},
		{"YamlToJson", func() ([]byte, error) { return transformer.YamlToJson(yamlInput) }},
		{"YamlToXml", func() ([]byte, error) { return transformer.YamlToXml(yamlInput) }},
		{"XmlToJson", func() ([]byte, error) { return transformer.XmlToJson(xmlInput) }},
		{"XmlToYaml", func() ([]byte, error) { return transformer.XmlToYaml(xmlInput) }},
	}

	for _, c := range conversions {
		t.Run(c.name, func(t *testing.T) {
			var first string
			for i := 0; i < 25; i++ {
				out, err := c.run()
				if err != nil {
					t.Fatalf("%s: %v", c.name, err)
				}
				if i == 0 {
					first = string(out)
					continue
				}
				if string(out) != first {
					t.Fatalf("%s output changed between runs:\nrun 0: %s\nrun %d: %s", c.name, first, i, out)
				}
			}
		})
	}
}

// TestCrossFormat_YamlOutputIsAlwaysReadable converts every fixture to YAML and
// reads the result back. A writer that emits YAML nothing can parse is the
// failure this whole file exists to keep out.
func TestCrossFormat_YamlOutputIsAlwaysReadable(t *testing.T) {
	jsonFixture, err := os.ReadFile("example/files/valid.json")
	if err != nil {
		t.Fatalf("reading valid.json: %v", err)
	}
	xmlFixture, err := os.ReadFile("example/files/valid.xml")
	if err != nil {
		t.Fatalf("reading valid.xml: %v", err)
	}
	yamlFixture, err := os.ReadFile("example/files/valid.yaml")
	if err != nil {
		t.Fatalf("reading valid.yaml: %v", err)
	}

	sources := []struct {
		name string
		run  func() ([]byte, error)
	}{
		{"from JSON", func() ([]byte, error) { return transformer.JsonToYaml(jsonFixture) }},
		{"from XML", func() ([]byte, error) { return transformer.XmlToYaml(xmlFixture) }},
		{"from YAML", func() ([]byte, error) {
			n, err := tyaml.DecodeYaml(yamlFixture)
			if err != nil {
				return nil, err
			}
			return tyaml.NodeToYaml(n)
		}},
	}

	for _, s := range sources {
		t.Run(s.name, func(t *testing.T) {
			out, err := s.run()
			if err != nil {
				t.Fatalf("%s: %v", s.name, err)
			}
			var reparsed any
			if err := yaml.Unmarshal(out, &reparsed); err != nil {
				t.Fatalf("produced YAML cannot be read back: %v\n%s", err, out)
			}
		})
	}
}

// TestCrossFormat_JsonOutputIsAlwaysValid does the same for the JSON writer,
// which now escapes and writes the document itself instead of handing it to
// encoding/json.
func TestCrossFormat_JsonOutputIsAlwaysValid(t *testing.T) {
	xmlFixture, err := os.ReadFile("example/files/valid.xml")
	if err != nil {
		t.Fatalf("reading valid.xml: %v", err)
	}
	yamlFixture, err := os.ReadFile("example/files/valid.yaml")
	if err != nil {
		t.Fatalf("reading valid.yaml: %v", err)
	}
	jsonFixture, err := os.ReadFile("example/files/valid.json")
	if err != nil {
		t.Fatalf("reading valid.json: %v", err)
	}

	sources := []struct {
		name string
		run  func() ([]byte, error)
	}{
		{"from XML", func() ([]byte, error) { return transformer.XmlToJson(xmlFixture) }},
		{"from YAML", func() ([]byte, error) { return transformer.YamlToJson(yamlFixture) }},
		{"from JSON", func() ([]byte, error) {
			n, err := tjson.DecodeJson(jsonFixture)
			if err != nil {
				return nil, err
			}
			return tjson.NodeToJson(n)
		}},
	}

	for _, s := range sources {
		t.Run(s.name, func(t *testing.T) {
			out, err := s.run()
			if err != nil {
				t.Fatalf("%s: %v", s.name, err)
			}
			if !tjson.IsJson(out) {
				t.Fatalf("produced JSON is not valid:\n%s", out)
			}
		})
	}
}
