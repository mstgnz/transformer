package txml

import (
	"os"
	"strings"
	"testing"

	"github.com/mstgnz/transformer/node"
)

func TestIsXml(t *testing.T) {
	validXml, err := os.ReadFile("../example/files/valid.xml")
	if err != nil {
		t.Fatalf("Error reading valid.xml: %v", err)
	}
	invalidXml, err := os.ReadFile("../example/files/invalid.xml")
	if err != nil {
		t.Fatalf("Error reading invalid.xml: %v", err)
	}

	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{"valid xml", validXml, true},
		{"invalid xml", invalidXml, false},
		{"empty xml", []byte{}, false},
		{"simple element", []byte("<root/>"), true},
		{"with attributes", []byte(`<root id="1"/>`), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsXml(tt.data); got != tt.want {
				t.Errorf("IsXml() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadXml(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{"valid xml file", "../example/files/valid.xml", false},
		{"invalid xml file", "../example/files/invalid.xml", true},
		{"non-existent file", "nonexistent.xml", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadXml(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadXml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !IsXml(got) {
				t.Error("ReadXml() returned invalid XML")
			}
		})
	}
}

func TestDecodeXml(t *testing.T) {
	validXml, err := os.ReadFile("../example/files/valid.xml")
	if err != nil {
		t.Fatalf("Error reading valid.xml: %v", err)
	}
	invalidXml, err := os.ReadFile("../example/files/invalid.xml")
	if err != nil {
		t.Fatalf("Error reading invalid.xml: %v", err)
	}

	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{"valid xml", validXml, false},
		{"invalid xml", invalidXml, true},
		{"element with number text", []byte("<root><count>42</count></root>"), false},
		{"element with bool text", []byte("<root><active>true</active></root>"), false},
		{"element with attributes", []byte(`<root><item id="1" name="test">value</item></root>`), false},
		{"nested elements", []byte("<root><a><b><c>deep</c></b></a></root>"), false},
		{"empty element", []byte("<root><empty/></root>"), false},
		{"root element", []byte("<root><child>val</child></root>"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, err := DecodeXml(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeXml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				xmlBytes, err := NodeToXml(n)
				if err != nil {
					t.Errorf("NodeToXml() error = %v", err)
					return
				}
				gotXml := normalizeXml(string(xmlBytes))
				wantXml := normalizeXml(string(tt.data))
				if gotXml != wantXml {
					t.Errorf("XML conversion mismatch\ngot:  %s\nwant: %s", gotXml, wantXml)
				}
			}
		})
	}
}

func TestNodeToXml(t *testing.T) {
	tests := []struct {
		name string
		xml  string
		want string
	}{
		{
			name: "simple element",
			xml:  "<root><test>value</test></root>",
			want: "<root><test>value</test></root>",
		},
		{
			name: "element with attribute",
			xml:  `<root><test attr="value">content</test></root>`,
			want: `<root><test attr="value">content</test></root>`,
		},
		{
			name: "nested elements",
			xml:  "<root><parent><child>value</child></parent></root>",
			want: "<root><parent><child>value</child></parent></root>",
		},
		{
			name: "multiple attributes",
			xml:  `<root><test a="1" b="2">value</test></root>`,
			want: `<root><test a="1" b="2">value</test></root>`,
		},
		{
			name: "empty element",
			xml:  "<root><test/></root>",
			want: "<root><test/></root>",
		},
		{
			name: "number text",
			xml:  "<root><count>99</count></root>",
			want: "<root><count>99</count></root>",
		},
		{
			name: "bool text",
			xml:  "<root><ok>true</ok></root>",
			want: "<root><ok>true</ok></root>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, err := DecodeXml([]byte(tt.xml))
			if err != nil {
				t.Fatalf("DecodeXml() error = %v", err)
			}
			gotXml, err := NodeToXml(n)
			if err != nil {
				t.Fatalf("NodeToXml() error = %v", err)
			}
			got := normalizeXml(string(gotXml))
			want := normalizeXml(tt.want)
			if got != want {
				t.Errorf("NodeToXml() = %v, want %v", got, want)
			}
		})
	}
}

func TestNodeToXml_NilNode(t *testing.T) {
	_, err := NodeToXml(nil)
	if err == nil {
		t.Error("NodeToXml(nil) should return error")
	}
}

func TestNodeToXml_SpecialTypes(t *testing.T) {
	// TypeNull node
	n := &node.Node{
		Key:   "root",
		Value: &node.Value{Type: node.TypeNull},
	}
	out, err := NodeToXml(n)
	if err != nil {
		t.Fatalf("NodeToXml(TypeNull) error: %v", err)
	}
	if !strings.Contains(string(out), "root") {
		t.Errorf("expected 'root' tag in output, got: %s", string(out))
	}

	// TypeBoolean node
	n2 := &node.Node{
		Key:   "root",
		Value: &node.Value{Type: node.TypeObject},
	}
	boolChild := &node.Node{
		Key:   "active",
		Value: &node.Value{Type: node.TypeBoolean, Worth: "true"},
	}
	n2.AddToEnd(boolChild)
	out2, err := NodeToXml(n2)
	if err != nil {
		t.Fatalf("NodeToXml(TypeBoolean) error: %v", err)
	}
	if !strings.Contains(string(out2), "true") {
		t.Errorf("expected 'true' in output, got: %s", string(out2))
	}
}

func TestNodeToXml_ArrayWithObjects(t *testing.T) {
	// Construct an array node with object items
	itemNode := &node.Node{
		Key: "item0",
		Value: &node.Value{
			Type: node.TypeObject,
			Node: &node.Node{
				Key:   "name",
				Value: &node.Value{Type: node.TypeString, Worth: "alice"},
			},
		},
	}
	n := &node.Node{
		Key: "root",
		Value: &node.Value{
			Type: node.TypeObject,
			Node: &node.Node{
				Key: "users",
				Value: &node.Value{
					Type: node.TypeArray,
					Array: []*node.Value{
						{Type: node.TypeObject, Node: itemNode},
						{Type: node.TypeString, Worth: "item"},
					},
				},
			},
		},
	}
	out, err := NodeToXml(n)
	if err != nil {
		t.Fatalf("NodeToXml() error: %v", err)
	}
	if !strings.Contains(string(out), "alice") {
		t.Errorf("expected 'alice' in output, got: %s", string(out))
	}
}

func TestEscapeXml(t *testing.T) {
	// Build a node with special chars that need escaping
	n := &node.Node{
		Key: "root",
		Value: &node.Value{
			Type: node.TypeObject,
			Node: &node.Node{
				Key:   "text",
				Value: &node.Value{Type: node.TypeString, Worth: `a & b < c > d " e ' f`},
			},
		},
	}
	out, err := NodeToXml(n)
	if err != nil {
		t.Fatalf("NodeToXml() error: %v", err)
	}
	outStr := string(out)
	checks := []string{"&amp;", "&lt;", "&gt;", "&quot;", "&apos;"}
	for _, ch := range checks {
		if !strings.Contains(outStr, ch) {
			t.Errorf("expected %q in escaped output, got: %s", ch, outStr)
		}
	}
}

func TestNamespace_PrefixedElement(t *testing.T) {
	input := `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><msg>hello</msg></soap:Body></soap:Envelope>`
	n, err := DecodeXml([]byte(input))
	if err != nil {
		t.Fatalf("DecodeXml error: %v", err)
	}

	// DecodeXml wraps in an internal "root" node; the real root is the first child.
	envelope := n.Value.Node
	if envelope == nil || envelope.Key != "soap:Envelope" {
		t.Fatalf("expected first child 'soap:Envelope', got %v", envelope)
	}

	// The xmlns:soap declaration should be on the soap:Envelope node.
	found := false
	current := envelope.Value.Node
	for current != nil {
		if current.Key == "@xmlns:soap" {
			found = true
			if current.Value.Worth != "http://schemas.xmlsoap.org/soap/envelope/" {
				t.Errorf("unexpected xmlns:soap value: %q", current.Value.Worth)
			}
		}
		current = current.Next
	}
	if !found {
		t.Error("expected @xmlns:soap attribute on soap:Envelope")
	}

	// Round-trip
	out, err := NodeToXml(n)
	if err != nil {
		t.Fatalf("NodeToXml error: %v", err)
	}
	outStr := normalizeXml(string(out))
	if !strings.Contains(outStr, "soap:Envelope") {
		t.Errorf("expected 'soap:Envelope' in output, got: %s", outStr)
	}
	if !strings.Contains(outStr, `xmlns:soap=`) {
		t.Errorf("expected xmlns:soap declaration in output, got: %s", outStr)
	}
}

func TestNamespace_DefaultNamespace(t *testing.T) {
	// Root element named "root" replaces the internal root node
	input := `<root xmlns="http://www.w3.org/1999/xhtml"><body><p>hello</p></body></root>`
	n, err := DecodeXml([]byte(input))
	if err != nil {
		t.Fatalf("DecodeXml error: %v", err)
	}

	// @xmlns should be on the root node itself
	found := false
	current := n.Value.Node
	for current != nil {
		if current.Key == "@xmlns" {
			found = true
		}
		current = current.Next
	}
	if !found {
		t.Error("expected @xmlns attribute for default namespace")
	}

	// Round-trip preserves xmlns declaration
	out, err := NodeToXml(n)
	if err != nil {
		t.Fatalf("NodeToXml error: %v", err)
	}
	if !strings.Contains(string(out), `xmlns="http://www.w3.org/1999/xhtml"`) {
		t.Errorf("expected xmlns declaration in output, got: %s", string(out))
	}
}

func TestNamespace_MultipleNamespaces(t *testing.T) {
	// Root element named "root" replaces the internal root; xmlns attrs go on it
	input := `<root xmlns:a="http://a.com" xmlns:b="http://b.com"><a:item>1</a:item><b:item>2</b:item></root>`
	n, err := DecodeXml([]byte(input))
	if err != nil {
		t.Fatalf("DecodeXml error: %v", err)
	}

	// Skip xmlns attributes on root, find first element child
	current := n.Value.Node
	for current != nil && strings.HasPrefix(current.Key, "@") {
		current = current.Next
	}
	if current == nil || current.Key != "a:item" {
		t.Errorf("expected first non-attr child 'a:item', got %v", current)
	}

	// Round-trip
	out, err := NodeToXml(n)
	if err != nil {
		t.Fatalf("NodeToXml error: %v", err)
	}
	outStr := string(out)
	if !strings.Contains(outStr, "a:item") || !strings.Contains(outStr, "b:item") {
		t.Errorf("expected prefixed elements in output, got: %s", outStr)
	}
}

func TestNamespace_NamespaceInheritedFromParent(t *testing.T) {
	// Namespace declared on parent, child elements use the prefix
	input := `<root xmlns:ns="http://example.com"><ns:child><ns:leaf>val</ns:leaf></ns:child></root>`
	n, err := DecodeXml([]byte(input))
	if err != nil {
		t.Fatalf("DecodeXml error: %v", err)
	}

	out, err := NodeToXml(n)
	if err != nil {
		t.Fatalf("NodeToXml error: %v", err)
	}
	if !strings.Contains(string(out), "ns:child") {
		t.Errorf("expected 'ns:child' in output, got: %s", string(out))
	}
}

func TestNamespace_PrefixedAttribute(t *testing.T) {
	input := `<root xmlns:xml="http://www.w3.org/XML/1998/namespace"><item xml:lang="en">hello</item></root>`
	n, err := DecodeXml([]byte(input))
	if err != nil {
		t.Fatalf("DecodeXml error: %v", err)
	}
	out, err := NodeToXml(n)
	if err != nil {
		t.Fatalf("NodeToXml error: %v", err)
	}
	if !strings.Contains(string(out), "hello") {
		t.Errorf("expected content in output, got: %s", string(out))
	}
}

func TestXmlRoundTrip(t *testing.T) {
	xmlData := []byte(`<root><person><name>Alice</name><age>30</age></person></root>`)
	n, err := DecodeXml(xmlData)
	if err != nil {
		t.Fatalf("DecodeXml error: %v", err)
	}
	out, err := NodeToXml(n)
	if err != nil {
		t.Fatalf("NodeToXml error: %v", err)
	}
	n2, err := DecodeXml(out)
	if err != nil {
		t.Fatalf("DecodeXml round-trip error: %v", err)
	}
	if n2 == nil {
		t.Fatal("round-trip produced nil node")
	}
}

func TestNodeToXml_OnlyAttributes(t *testing.T) {
	// Node with only attribute children should produce self-closing tag
	xml := `<root><item id="1"/></root>`
	n, err := DecodeXml([]byte(xml))
	if err != nil {
		t.Fatalf("DecodeXml error: %v", err)
	}
	out, err := NodeToXml(n)
	if err != nil {
		t.Fatalf("NodeToXml error: %v", err)
	}
	if !strings.Contains(string(out), `id="1"`) {
		t.Errorf("expected attribute in output, got: %s", string(out))
	}
}

func TestNodeToXml_NilValueNode(t *testing.T) {
	// Node with nil Value — writeNodeToXml isEmpty path
	n := &node.Node{
		Key:   "root",
		Value: nil,
	}
	out, err := NodeToXml(n)
	if err != nil {
		t.Fatalf("NodeToXml error: %v", err)
	}
	if !strings.Contains(string(out), "<root/>") {
		t.Errorf("expected self-closing tag for nil value, got: %s", string(out))
	}
}

func TestHasNonAttributeChildren_NilValue(t *testing.T) {
	// hasNonAttributeChildren with nil Value
	n := &node.Node{Key: "test", Value: nil}
	if hasNonAttributeChildren(n) {
		t.Error("hasNonAttributeChildren should return false for nil Value")
	}
	// with non-nil Value but nil Node
	n2 := &node.Node{Key: "test", Value: &node.Value{Type: node.TypeObject}}
	if hasNonAttributeChildren(n2) {
		t.Error("hasNonAttributeChildren should return false for nil Value.Node")
	}
}

func TestDecodeXml_DepthLimit(t *testing.T) {
	depth := 110
	xml := "<root>"
	for i := 0; i < depth; i++ {
		xml += "<a>"
	}
	xml += "val"
	for i := 0; i < depth; i++ {
		xml += "</a>"
	}
	xml += "</root>"
	_, err := DecodeXml([]byte(xml))
	if err == nil {
		t.Error("expected error for deeply nested XML")
	}
}

func TestNodeToXml_DepthLimit(t *testing.T) {
	// Build a deeply nested node structure
	root := &node.Node{Key: "root", Value: &node.Value{Type: node.TypeObject}}
	current := root
	for i := 0; i <= maxDepth+2; i++ {
		child := &node.Node{Key: "a", Value: &node.Value{Type: node.TypeObject}}
		current.Value.Node = child
		child.Parent = current
		current = child
	}
	_, err := NodeToXml(root)
	if err == nil {
		t.Error("expected error for deeply nested node")
	}
}

// normalizeXml normalizes XML string by removing whitespace and newlines
func normalizeXml(xml string) string {
	xml = strings.ReplaceAll(xml, `<?xml version="1.0" encoding="UTF-8"?>`, "")
	xml = strings.ReplaceAll(xml, `<?xml version="1.0" encoding="UTF-8" ?>`, "")
	xml = strings.ReplaceAll(xml, "\n", "")
	xml = strings.ReplaceAll(xml, "\r", "")
	xml = strings.ReplaceAll(xml, "\t", "")
	xml = strings.ReplaceAll(xml, "  ", " ")
	xml = strings.TrimSpace(xml)
	return xml
}
