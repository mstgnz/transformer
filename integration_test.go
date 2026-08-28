package transformer_test

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	transformer "github.com/mstgnz/transformer"
	"github.com/mstgnz/transformer/node"
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

// TestCrossFormat_ShapesXmlCannotHold covers the values XML has no shape of its
// own for. They travel with a type marker, so they come back as what they were
// rather than as the empty tag they would otherwise collapse into.
func TestCrossFormat_ShapesXmlCannotHold(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{"null", `{"a":null}`},
		{"null among values", `{"s":"text","n":null,"i":42}`},
		{"nulls in a list", `{"a":[null,null]}`},
		{"nested null", `{"a":{"b":null}}`},
		{"empty string", `{"s":""}`},
		{"empty strings in a list", `{"a":["","x"]}`},
		{"empty list", `{"a":[]}`},
		{"empty object", `{"o":{}}`},
		{"all three empties", `{"o":{},"a":[],"s":""}`},
		{"list holding one value", `{"a":[1]}`},
		{"list holding one object", `{"a":[{"k":1}]}`},
		{"nested lists", `{"m":[[1,2],[3,4]]}`},
		{"deeply nested lists", `{"m":[[[1]]]}`},
		{"list of one element lists", `{"m":[[1],[2]]}`},
		{"top level list", `[1,2,3]`},
		{"top level one element list", `[1]`},
		{"top level empty list", `[]`},
		{"mixed", `{"a":[1,"two",true,null],"b":[],"c":null,"d":[5]}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xmlOut, err := transformer.JsonToXml([]byte(tt.json))
			if err != nil {
				t.Fatalf("JsonToXml: %v", err)
			}
			if !txml.IsXml(xmlOut) {
				t.Fatalf("produced XML is not valid: %s", xmlOut)
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

// TestCrossFormat_XmlStaysPlainWhenItCan is the counterweight to the markers:
// a document XML can already describe must be written without one, so ordinary
// XML does not grow attributes it does not need.
func TestCrossFormat_XmlStaysPlainWhenItCan(t *testing.T) {
	tests := []string{
		`{"name":"alice","age":30,"active":true}`,
		`{"ports":[80,443,8080]}`,
		`{"users":[{"name":"alice"},{"name":"bob"}]}`,
		`{"spec":{"containers":[{"name":"web"},{"name":"db"}]}}`,
		`{"a":{"b":{"c":"deep"}}}`,
	}
	for _, in := range tests {
		t.Run(in, func(t *testing.T) {
			out, err := transformer.JsonToXml([]byte(in))
			if err != nil {
				t.Fatalf("JsonToXml: %v", err)
			}
			if strings.Contains(string(out), "xmlns:") || strings.Contains(string(out), ":type=") {
				t.Errorf("a document XML can describe on its own carried a marker: %s", out)
			}
		})
	}
}

// TestCrossFormat_MarkersSurviveYaml checks that the marked shapes also survive
// the XML to YAML direction, which is where the values leave the tree again.
func TestCrossFormat_MarkersSurviveYaml(t *testing.T) {
	input := []byte(`{"n":null,"empty":[],"one":[5],"s":""}`)

	xmlOut, err := transformer.JsonToXml(input)
	if err != nil {
		t.Fatalf("JsonToXml: %v", err)
	}
	yamlOut, err := transformer.XmlToYaml(xmlOut)
	if err != nil {
		t.Fatalf("XmlToYaml: %v", err)
	}

	var got map[string]any
	if err := yaml.Unmarshal(yamlOut, &got); err != nil {
		t.Fatalf("produced YAML cannot be read back: %v\n%s", err, yamlOut)
	}
	if got["n"] != nil {
		t.Errorf("n = %#v, want nil", got["n"])
	}
	if list, ok := got["empty"].([]any); !ok || len(list) != 0 {
		t.Errorf("empty = %#v, want an empty list", got["empty"])
	}
	if list, ok := got["one"].([]any); !ok || len(list) != 1 {
		t.Errorf("one = %#v, want a list of one", got["one"])
	}
	if s, ok := got["s"].(string); !ok || s != "" {
		t.Errorf("s = %#v, want the empty string", got["s"])
	}
}

// TestCrossFormat_KeysXmlCannotSpell covers JSON keys that are not legal XML
// element names. They used to be written straight into the tag, producing a
// document no parser accepts, with no error to say so.
func TestCrossFormat_KeysXmlCannotSpell(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{"space in key", `{"first name":"John"}`},
		{"leading digit", `{"1st":"place"}`},
		{"slash", `{"a/b":"x"}`},
		{"percent", `{"c%d":"y"}`},
		{"at sign inside", `{"key with @":"x"}`},
		{"colon that is not a prefix", `{"tel:+90":"x"}`},
		{"colon shaped key", `{"configGlossary:installationAt":"Philadelphia, PA"}`},
		{"empty key", `{"":"x"}`},
		{"turkish", `{"şehir":"İstanbul"}`},
		{"japanese", `{"日本語":"x"}`},
		{"escape shaped key", `{"_x0020_":"x"}`},
		{"nested and listed", `{"a":{"b c":{"d e f":[1,2]}}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xmlOut, err := transformer.JsonToXml([]byte(tt.json))
			if err != nil {
				t.Fatalf("JsonToXml: %v", err)
			}
			if !txml.IsXml(xmlOut) {
				t.Fatalf("produced XML is not valid: %s", xmlOut)
			}
			back, err := transformer.XmlToJson(xmlOut)
			if err != nil {
				t.Fatalf("XmlToJson: %v (intermediate XML: %s)", err, xmlOut)
			}

			// Compared as data: the JSON writer escapes < and & in text, so the
			// bytes may differ where the document itself does not.
			var want, got any
			if err := json.Unmarshal([]byte(tt.json), &want); err != nil {
				t.Fatalf("test input is not valid JSON: %v", err)
			}
			if err := json.Unmarshal(back, &got); err != nil {
				t.Fatalf("output is not valid JSON: %v\n%s", err, back)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("round trip changed the document\ngot:  %s\nwant: %s\nvia XML: %s", back, tt.json, xmlOut)
			}
		})
	}
}

// TestCrossFormat_TheWholeFixtureSurvivesXml runs the JSON fixture, whose keys
// carry colons and whose lists are deep, all the way through XML.
func TestCrossFormat_TheWholeFixtureSurvivesXml(t *testing.T) {
	fixture, err := os.ReadFile("example/files/valid.json")
	if err != nil {
		t.Fatalf("reading valid.json: %v", err)
	}

	xmlOut, err := transformer.JsonToXml(fixture)
	if err != nil {
		t.Fatalf("JsonToXml: %v", err)
	}
	if !txml.IsXml(xmlOut) {
		t.Fatalf("produced XML is not valid")
	}
	back, err := transformer.XmlToJson(xmlOut)
	if err != nil {
		t.Fatalf("XmlToJson: %v", err)
	}

	var want, got any
	if err := json.Unmarshal(fixture, &want); err != nil {
		t.Fatalf("fixture is not valid JSON: %v", err)
	}
	if err := json.Unmarshal(back, &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Error("the fixture did not survive the trip through XML")
	}
}

// ---- Search ----

// TestSearch_FindsAKeyAcrossTheWholeFixture is the search case the library
// exists for beside conversion: one key, many places, and a way to tell them
// apart.
func TestSearch_FindsAKeyAcrossTheWholeFixture(t *testing.T) {
	fixture, err := os.ReadFile("example/files/valid.json")
	if err != nil {
		t.Fatalf("reading valid.json: %v", err)
	}

	viaHelper, err := transformer.SearchJson(fixture, "name")
	if err != nil {
		t.Fatalf("SearchJson: %v", err)
	}
	if len(viaHelper) < 20 {
		t.Fatalf("found %d matches for 'name', want the whole document's worth", len(viaHelper))
	}

	// Every reported path must address the node it was reported for. The tree
	// has to be the one that was searched, so this decodes it once and searches
	// that, rather than comparing against the helper's own tree.
	root, err := tjson.DecodeJson(fixture)
	if err != nil {
		t.Fatalf("DecodeJson: %v", err)
	}
	matches := root.Search("name")
	if len(matches) != len(viaHelper) {
		t.Errorf("SearchJson found %d matches, searching the tree directly found %d", len(viaHelper), len(matches))
	}
	for _, m := range matches {
		if got := root.GetNodeByPath(m.Path); got != m.Node {
			t.Errorf("path %q does not lead back to the node the search reported", m.Path)
		}
		if m.Depth < 1 {
			t.Errorf("match %q reported depth %d", m.Path, m.Depth)
		}
	}

	// The same key sits at many different levels, and the depths say which.
	depths := root.Depths("name")
	if len(depths) < 5 {
		t.Errorf("depths for 'name' = %v, want the key to be reported at several levels", depths)
	}
	for i := 1; i < len(depths); i++ {
		if depths[i] <= depths[i-1] {
			t.Errorf("depths are not in ascending order: %v", depths)
			break
		}
	}

	// Specific hits, to pin what the paths mean.
	for path, want := range map[string]string{
		"root.name":                              "root level name",
		"root.meta.name":                         "catalog",
		"root.groups[0].members[0].profile.name": "Alice Doe",
		"root.deep[0][0].items[0].entries[0].detail.name": "deepest",
		"root.registry.servers[1].replicas.name":          "cache-1-replicas",
	} {
		node := root.GetNodeByPath(path)
		if node == nil {
			t.Errorf("%s was not found", path)
			continue
		}
		if node.Value.Worth != want {
			t.Errorf("%s = %q, want %q", path, node.Value.Worth, want)
		}
	}
}

// TestSearch_ReachesIntoNestedArrays covers keys that live under a list of
// lists, which the earlier traversal walked straight past.
func TestSearch_ReachesIntoNestedArrays(t *testing.T) {
	fixture, err := os.ReadFile("example/files/valid.json")
	if err != nil {
		t.Fatalf("reading valid.json: %v", err)
	}
	root, err := tjson.DecodeJson(fixture)
	if err != nil {
		t.Fatalf("DecodeJson: %v", err)
	}

	// "level" only ever appears under deep[][], which is an array of arrays.
	matches := root.Search("level")
	if len(matches) != 4 {
		t.Fatalf("Search(level) found %d matches, want 4:\n%v", len(matches), matches)
	}
	if !root.Exists("level") {
		t.Error("Exists(level) = false for a key that is present inside an array of arrays")
	}

	// Values inside a list of lists are addressable too.
	value, ok := root.GetValueByPath("root.matrix[1][2]")
	if !ok || value.Worth != "6" {
		t.Errorf("matrix[1][2] = %v, %v; want the number 6", value, ok)
	}
	value, ok = root.GetValueByPath("root.cube[1][0][1]")
	if !ok || value.Worth != "6" {
		t.Errorf("cube[1][0][1] = %v, %v; want the number 6", value, ok)
	}
}

// TestSearch_WorksOnEveryFormat checks that a key is found the same way whether
// the document arrived as JSON, XML or YAML.
func TestSearch_WorksOnEveryFormat(t *testing.T) {
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

	searches := []struct {
		name string
		run  func() ([]node.Match, error)
		key  string
	}{
		{"json", func() ([]node.Match, error) { return transformer.SearchJson(jsonFixture, "name") }, "name"},
		{"xml", func() ([]node.Match, error) { return transformer.SearchXml(xmlFixture, "name") }, "name"},
		{"yaml", func() ([]node.Match, error) { return transformer.SearchYaml(yamlFixture, "name") }, "name"},
	}

	for _, s := range searches {
		t.Run(s.name, func(t *testing.T) {
			matches, err := s.run()
			if err != nil {
				t.Fatalf("search: %v", err)
			}
			if len(matches) < 5 {
				t.Fatalf("found %d matches, want the key reported wherever it appears", len(matches))
			}
			seen := make(map[int]bool)
			for _, m := range matches {
				if m.Path == "" {
					t.Error("a match came back without a path")
				}
				seen[m.Depth] = true
			}
			if len(seen) < 3 {
				t.Errorf("every match landed on %d depths; the key sits at several levels", len(seen))
			}
		})
	}
}

func TestSearch_InvalidInput(t *testing.T) {
	if _, err := transformer.SearchJson([]byte(`{`), "a"); err == nil {
		t.Error("SearchJson on broken JSON returned no error")
	}
	if _, err := transformer.SearchXml([]byte(`<a>`), "a"); err == nil {
		t.Error("SearchXml on broken XML returned no error")
	}
	if _, err := transformer.SearchYaml([]byte("a: : b"), "a"); err == nil {
		t.Error("SearchYaml on broken YAML returned no error")
	}
}

// ---- Find and replace ----

// TestEdit_FindAndReplaceAcrossFormats covers the other half of what the Node
// model is for: change something wherever it appears, and write the document
// back in the format it arrived in.
func TestEdit_FindAndReplaceAcrossFormats(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		in := []byte(`{"host":"staging.example.com","db":{"host":"staging.example.com","port":5432},"peers":[{"host":"staging.example.com"}]}`)

		out, err := transformer.EditJson(in, func(n *node.Node) {
			if changed := n.ReplaceText("staging.", ""); changed != 3 {
				t.Errorf("ReplaceText changed %d values, want 3", changed)
			}
			if renamed := n.Rename("host", "address"); renamed != 3 {
				t.Errorf("Rename changed %d keys, want 3", renamed)
			}
			if removed := n.Remove("port"); removed != 1 {
				t.Errorf("Remove deleted %d members, want 1", removed)
			}
		})
		if err != nil {
			t.Fatalf("EditJson: %v", err)
		}

		want := `{"address":"example.com","db":{"address":"example.com"},"peers":[{"address":"example.com"}]}`
		if string(out) != want {
			t.Errorf("got:  %s\nwant: %s", out, want)
		}
	})

	t.Run("yaml", func(t *testing.T) {
		in := []byte("host: staging.example.com\ndb:\n  host: staging.example.com\n")

		out, err := transformer.EditYaml(in, func(n *node.Node) {
			n.ReplaceText("staging.", "")
		})
		if err != nil {
			t.Fatalf("EditYaml: %v", err)
		}

		var got map[string]any
		if err := yaml.Unmarshal(out, &got); err != nil {
			t.Fatalf("the result is not readable YAML: %v\n%s", err, out)
		}
		if got["host"] != "example.com" {
			t.Errorf("host = %v, want example.com", got["host"])
		}
	})

	t.Run("xml", func(t *testing.T) {
		in := []byte(`<root><host>staging.example.com</host><host>staging.other.com</host></root>`)

		out, err := transformer.EditXml(in, func(n *node.Node) {
			if changed := n.ReplaceText("staging.", ""); changed != 2 {
				t.Errorf("ReplaceText changed %d values, want 2", changed)
			}
			n.Rename("host", "server")
		})
		if err != nil {
			t.Fatalf("EditXml: %v", err)
		}
		want := `<root><server>example.com</server><server>other.com</server></root>`
		if got := string(out)[len(`<?xml version="1.0" encoding="UTF-8"?>`)+1:]; got != want {
			t.Errorf("got:  %s\nwant: %s", got, want)
		}
	})
}

// TestEdit_SearchThenReplaceByPath is the workflow the two halves are meant to
// be used in: find where a key sits, then change exactly the one that matters.
func TestEdit_SearchThenReplaceByPath(t *testing.T) {
	fixture, err := os.ReadFile("example/files/valid.json")
	if err != nil {
		t.Fatalf("reading valid.json: %v", err)
	}

	root, err := tjson.DecodeJson(fixture)
	if err != nil {
		t.Fatalf("DecodeJson: %v", err)
	}

	// Find every "name", then change only the deepest one.
	matches := root.Search("name")
	if len(matches) < 20 {
		t.Fatalf("found %d matches, want the whole document's worth", len(matches))
	}
	deepest := matches[0]
	for _, m := range matches {
		if m.Depth > deepest.Depth {
			deepest = m
		}
	}

	if !root.SetByPath(deepest.Path, &node.Value{Type: node.TypeString, Worth: "renamed"}) {
		t.Fatalf("SetByPath(%q) did not land", deepest.Path)
	}
	if got, ok := root.GetValueByPath(deepest.Path); !ok || got.Worth != "renamed" {
		t.Errorf("%s = %v, want renamed", deepest.Path, got)
	}

	// Nothing else moved: the count is the same and the tree is still sound.
	if after := len(root.Search("name")); after != len(matches) {
		t.Errorf("the document now holds %d names, want %d", after, len(matches))
	}
	if err := root.Validate(); err != nil {
		t.Errorf("the tree is no longer sound: %v", err)
	}

	// And it still writes back as JSON, XML and YAML.
	out, err := tjson.NodeToJson(root)
	if err != nil {
		t.Fatalf("NodeToJson: %v", err)
	}
	if !tjson.IsJson(out) {
		t.Error("the edited document is not valid JSON")
	}
	if _, err := txml.NodeToXml(root); err != nil {
		t.Errorf("NodeToXml after editing: %v", err)
	}
	if _, err := tyaml.NodeToYaml(root); err != nil {
		t.Errorf("NodeToYaml after editing: %v", err)
	}
}

// TestEdit_SurvivesAConversion checks that a document edited in one format can
// still be written in another.
func TestEdit_SurvivesAConversion(t *testing.T) {
	in := []byte(`{"a":[1,2],"b":{"c":"text"},"d":null}`)

	out, err := transformer.EditJson(in, func(n *node.Node) {
		n.SetByPath("root.a[0]", &node.Value{Type: node.TypeNumber, Worth: "10"})
		n.SetByPath("root.added", &node.Value{Type: node.TypeString, Worth: "new"})
		n.ReplaceText("text", "changed")
	})
	if err != nil {
		t.Fatalf("EditJson: %v", err)
	}
	want := `{"a":[10,2],"b":{"c":"changed"},"d":null,"added":"new"}`
	if string(out) != want {
		t.Fatalf("got:  %s\nwant: %s", out, want)
	}

	xmlOut, err := transformer.JsonToXml(out)
	if err != nil {
		t.Fatalf("JsonToXml: %v", err)
	}
	back, err := transformer.XmlToJson(xmlOut)
	if err != nil {
		t.Fatalf("XmlToJson: %v", err)
	}
	if string(back) != want {
		t.Errorf("the edited document did not survive a trip through XML\ngot:  %s\nwant: %s", back, want)
	}
}

func TestEdit_InvalidInput(t *testing.T) {
	if _, err := transformer.EditJson([]byte(`{`), nil); err == nil {
		t.Error("EditJson on broken JSON returned no error")
	}
	if _, err := transformer.EditXml([]byte(`<a>`), nil); err == nil {
		t.Error("EditXml on broken XML returned no error")
	}
	if _, err := transformer.EditYaml([]byte("a: : b"), nil); err == nil {
		t.Error("EditYaml on broken YAML returned no error")
	}

	// No callback means no change, not a failure.
	out, err := transformer.EditJson([]byte(`{"a":1}`), nil)
	if err != nil || string(out) != `{"a":1}` {
		t.Errorf("EditJson with no callback = %s, %v", out, err)
	}
}
