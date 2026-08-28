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

// ---- Repeated elements ----

// TestDecodeXml_RepeatedElementsBecomeArray pins the fix for the loss that used
// to happen one step later: two siblings sharing a name reached the JSON writer
// as two nodes with one key, and the map it built kept only the last of them.
func TestDecodeXml_RepeatedElementsBecomeArray(t *testing.T) {
	n, err := DecodeXml([]byte(`<root><item>1</item><item>2</item><item>3</item></root>`))
	if err != nil {
		t.Fatalf("DecodeXml error: %v", err)
	}

	item := n.Value.Node
	if item == nil || item.Key != "item" {
		t.Fatalf("expected a single 'item' node, got %v", item)
	}
	if item.Next != nil {
		t.Errorf("expected the repeated elements to be merged, found sibling %q", item.Next.Key)
	}
	if item.Value.Type != node.TypeArray {
		t.Fatalf("item type = %v, want array", item.Value.Type)
	}
	if len(item.Value.Array) != 3 {
		t.Fatalf("expected 3 values, got %d", len(item.Value.Array))
	}
	for i, want := range []string{"1", "2", "3"} {
		if got := item.Value.Array[i].Worth; got != want {
			t.Errorf("value %d = %q, want %q", i, got, want)
		}
	}
}

func TestDecodeXml_RepeatedElementsKeepAttributes(t *testing.T) {
	n, err := DecodeXml([]byte(`<root><port p="34">34</port><port>55</port></root>`))
	if err != nil {
		t.Fatalf("DecodeXml error: %v", err)
	}

	port := n.Value.Node
	if port == nil || port.Value.Type != node.TypeArray || len(port.Value.Array) != 2 {
		t.Fatalf("expected an array of 2 ports, got %v", port)
	}

	first := port.Value.Array[0]
	if first.Node == nil || first.Node.Value == nil {
		t.Fatal("the element carrying an attribute lost its node")
	}
	var attr, text *node.Node
	for cur := first.Node.Value.Node; cur != nil; cur = cur.Next {
		switch cur.Key {
		case "@p":
			attr = cur
		case textKey:
			text = cur
		}
	}
	if attr == nil || attr.Value.Worth != "34" {
		t.Errorf("attribute lost: %v", attr)
	}
	if text == nil || text.Value.Worth != "34" {
		t.Errorf("text lost: %v", text)
	}

	if second := port.Value.Array[1]; second.Worth != "55" {
		t.Errorf("plain element = %q, want 55", second.Worth)
	}
}

// TestXmlRoundTrip_RepeatedElements checks the writer puts the list back the
// way the document held it.
func TestXmlRoundTrip_RepeatedElements(t *testing.T) {
	tests := []string{
		`<root><item>1</item><item>2</item></root>`,
		`<root><port p="34">34</port><port>55</port></root>`,
		`<root><list><a>1</a><a>2</a></list><other>x</other></root>`,
		`<root><group><n>1</n></group><group><n>2</n></group></root>`,
	}
	for _, in := range tests {
		t.Run(in, func(t *testing.T) {
			n, err := DecodeXml([]byte(in))
			if err != nil {
				t.Fatalf("DecodeXml error: %v", err)
			}
			out, err := NodeToXml(n)
			if err != nil {
				t.Fatalf("NodeToXml error: %v", err)
			}
			if got := normalizeXml(string(out)); got != in {
				t.Errorf("round trip changed the document\ngot:  %s\nwant: %s", got, in)
			}
		})
	}
}

// ---- Attributes beside text ----

// TestDecodeXml_AttributesAndText covers the element shape whose attributes
// used to be unreachable once its text was read.
func TestDecodeXml_AttributesAndText(t *testing.T) {
	n, err := DecodeXml([]byte(`<root><name lang="en" id="7">John</name></root>`))
	if err != nil {
		t.Fatalf("DecodeXml error: %v", err)
	}

	name := n.Value.Node
	if name == nil || name.Key != "name" {
		t.Fatalf("expected a 'name' node, got %v", name)
	}
	if name.Value.Type != node.TypeObject {
		t.Fatalf("type = %v, want object so attributes and text can sit side by side", name.Value.Type)
	}

	got := make(map[string]string)
	for cur := name.Value.Node; cur != nil; cur = cur.Next {
		got[cur.Key] = cur.Value.Worth
	}
	for key, want := range map[string]string{"@lang": "en", "@id": "7", textKey: "John"} {
		if got[key] != want {
			t.Errorf("%s = %q, want %q (all of: %v)", key, got[key], want, got)
		}
	}
}

func TestXmlRoundTrip_AttributesAndText(t *testing.T) {
	in := `<root id="5"><name lang="en">John</name><plain>text</plain><bare a="1"/></root>`
	n, err := DecodeXml([]byte(in))
	if err != nil {
		t.Fatalf("DecodeXml error: %v", err)
	}
	out, err := NodeToXml(n)
	if err != nil {
		t.Fatalf("NodeToXml error: %v", err)
	}
	if got := normalizeXml(string(out)); got != in {
		t.Errorf("round trip changed the document\ngot:  %s\nwant: %s", got, in)
	}
}

// TestDecodeXml_TextOnlyElementStaysScalar keeps the common case simple: an
// element with nothing but text must not grow a "#text" child.
func TestDecodeXml_TextOnlyElementStaysScalar(t *testing.T) {
	n, err := DecodeXml([]byte(`<root><name>John</name><count>42</count><ok>true</ok></root>`))
	if err != nil {
		t.Fatalf("DecodeXml error: %v", err)
	}

	want := map[string]node.ValueType{
		"name":  node.TypeString,
		"count": node.TypeNumber,
		"ok":    node.TypeBoolean,
	}
	for cur := n.Value.Node; cur != nil; cur = cur.Next {
		if cur.Value.Node != nil {
			t.Errorf("%s grew a child node it does not need", cur.Key)
		}
		if got := cur.Value.Type; got != want[cur.Key] {
			t.Errorf("%s type = %v, want %v", cur.Key, got, want[cur.Key])
		}
	}
}

// TestDecodeXml_TextTypeIsNotOverEager checks the narrow reading of what counts
// as a number or a boolean. strconv.ParseBool would call "t" true, and
// strconv.ParseFloat would call "Inf" a number.
func TestDecodeXml_TextTypeIsNotOverEager(t *testing.T) {
	n, err := DecodeXml([]byte(`<root><a>t</a><b>Inf</b><c>NaN</c><d>0x1F</d><e>1</e></root>`))
	if err != nil {
		t.Fatalf("DecodeXml error: %v", err)
	}

	want := map[string]node.ValueType{
		"a": node.TypeString,
		"b": node.TypeString,
		"c": node.TypeString,
		"d": node.TypeString,
		"e": node.TypeNumber,
	}
	for cur := n.Value.Node; cur != nil; cur = cur.Next {
		if got := cur.Value.Type; got != want[cur.Key] {
			t.Errorf("%s (%q) type = %v, want %v", cur.Key, cur.Value.Worth, got, want[cur.Key])
		}
	}
}

// ---- Writing arrays ----

// TestNodeToXml_ArrayWritesOneElementPerValue covers a list that arrived from
// JSON or YAML rather than from XML.
func TestNodeToXml_ArrayWritesOneElementPerValue(t *testing.T) {
	n := &node.Node{
		Key: "root",
		Value: &node.Value{
			Type: node.TypeObject,
			Node: &node.Node{
				Key: "tag",
				Value: &node.Value{
					Type: node.TypeArray,
					Array: []*node.Value{
						{Type: node.TypeString, Worth: "go"},
						{Type: node.TypeString, Worth: "xml"},
					},
				},
			},
		},
	}
	out, err := NodeToXml(n)
	if err != nil {
		t.Fatalf("NodeToXml error: %v", err)
	}
	want := `<root><tag>go</tag><tag>xml</tag></root>`
	if got := normalizeXml(string(out)); got != want {
		t.Errorf("got:  %s\nwant: %s", got, want)
	}
}

func TestNodeToXml_EmptyArrayKeepsTheTag(t *testing.T) {
	n := &node.Node{
		Key: "root",
		Value: &node.Value{
			Type: node.TypeObject,
			Node: &node.Node{
				Key:   "tags",
				Value: &node.Value{Type: node.TypeArray, Array: []*node.Value{}},
			},
		},
	}
	out, err := NodeToXml(n)
	if err != nil {
		t.Fatalf("NodeToXml error: %v", err)
	}
	if got := normalizeXml(string(out)); got != `<root><tags/></root>` {
		t.Errorf("got: %s, want <root><tags/></root>", got)
	}
}

func TestNodeToXml_RootArray(t *testing.T) {
	n := &node.Node{
		Key: "root",
		Value: &node.Value{
			Type: node.TypeArray,
			Array: []*node.Value{
				{Type: node.TypeNumber, Worth: "1"},
				{Type: node.TypeNumber, Worth: "2"},
			},
		},
	}
	out, err := NodeToXml(n)
	if err != nil {
		t.Fatalf("NodeToXml error: %v", err)
	}
	want := `<root><item>1</item><item>2</item></root>`
	if got := normalizeXml(string(out)); got != want {
		t.Errorf("got:  %s\nwant: %s", got, want)
	}
}

// ---- Fixture ----

// TestXmlRoundTrip_Fixture runs the whole example document, which carries
// repeated elements, attributes and attributes beside text all at once.
func TestXmlRoundTrip_Fixture(t *testing.T) {
	data, err := os.ReadFile("../example/files/valid.xml")
	if err != nil {
		t.Fatalf("reading valid.xml: %v", err)
	}

	n, err := DecodeXml(data)
	if err != nil {
		t.Fatalf("DecodeXml error: %v", err)
	}
	out, err := NodeToXml(n)
	if err != nil {
		t.Fatalf("NodeToXml error: %v", err)
	}

	if got, want := normalizeXml(string(out)), normalizeXml(string(data)); got != want {
		t.Errorf("round trip changed the document\ngot:  %s\nwant: %s", got, want)
	}

	// The values that used to be dropped must all still be reachable.
	for _, key := range []string{"containers", "port", "status", "test"} {
		found := n.GetNode(key)
		if len(found) == 0 {
			t.Errorf("%q is missing from the decoded tree", key)
			continue
		}
		if found[0].Value.Type != node.TypeArray {
			t.Errorf("%q type = %v, want array", key, found[0].Value.Type)
		}
	}
}

// ---- Mixed content ----

// TestDecodeXml_MixedContent covers text interleaved with child elements. The
// text runs are collected into the element's "#text" node, so none of them is
// lost, though they are written back before the child elements rather than
// between them.
func TestDecodeXml_MixedContent(t *testing.T) {
	n, err := DecodeXml([]byte(`<root>hello<b>x</b>world</root>`))
	if err != nil {
		t.Fatalf("DecodeXml error: %v", err)
	}

	var text *node.Node
	for cur := n.Value.Node; cur != nil; cur = cur.Next {
		if cur.Key == textKey {
			text = cur
		}
	}
	if text == nil {
		t.Fatal("the text runs are missing entirely")
	}
	if text.Value.Type != node.TypeArray || len(text.Value.Array) != 2 {
		t.Fatalf("expected both text runs, got %+v", text.Value)
	}
	if text.Value.Array[0].Worth != "hello" || text.Value.Array[1].Worth != "world" {
		t.Errorf("text runs = %q and %q, want hello and world",
			text.Value.Array[0].Worth, text.Value.Array[1].Worth)
	}

	out, err := NodeToXml(n)
	if err != nil {
		t.Fatalf("NodeToXml error: %v", err)
	}
	want := `<root>helloworld<b>x</b></root>`
	if got := normalizeXml(string(out)); got != want {
		t.Errorf("got:  %s\nwant: %s", got, want)
	}
}

// ---- Nested lists ----

// TestNodeToXml_NestedArray covers a list of lists, which XML has no name for.
// The inner values are written as <item> elements so nothing disappears.
func TestNodeToXml_NestedArray(t *testing.T) {
	n := &node.Node{
		Key: "root",
		Value: &node.Value{
			Type: node.TypeObject,
			Node: &node.Node{
				Key: "m",
				Value: &node.Value{
					Type: node.TypeArray,
					Array: []*node.Value{
						{Type: node.TypeArray, Array: []*node.Value{
							{Type: node.TypeNumber, Worth: "1"},
							{Type: node.TypeNumber, Worth: "2"},
						}},
					},
				},
			},
		},
	}
	out, err := NodeToXml(n)
	if err != nil {
		t.Fatalf("NodeToXml error: %v", err)
	}
	want := `<root><m><item>1</item><item>2</item></m></root>`
	if got := normalizeXml(string(out)); got != want {
		t.Errorf("got:  %s\nwant: %s", got, want)
	}
}

// ---- Helper edge cases ----

func TestUnwrapItem(t *testing.T) {
	if got := unwrapItem(nil); got != nil {
		t.Errorf("unwrapItem(nil) = %v, want nil", got)
	}

	scalarItem := &node.Value{Type: node.TypeString, Worth: "x"}
	if got := unwrapItem(scalarItem); got != scalarItem {
		t.Error("a scalar element should be returned as it is stored")
	}

	inner := &node.Value{Type: node.TypeObject}
	wrapped := &node.Value{Type: node.TypeObject, Node: &node.Node{Key: "item0", Value: inner}}
	if got := unwrapItem(wrapped); got != inner {
		t.Error("an object element should be unwrapped to the object it holds")
	}
}

func TestIsEmptyElement(t *testing.T) {
	tests := []struct {
		name  string
		value *node.Value
		want  bool
	}{
		{"nil", nil, true},
		{"null", &node.Value{Type: node.TypeNull}, true},
		{"empty string", &node.Value{Type: node.TypeString}, true},
		{"text", &node.Value{Type: node.TypeString, Worth: "x"}, false},
		{"empty list", &node.Value{Type: node.TypeArray, Array: []*node.Value{}}, true},
		{"list with a value", &node.Value{Type: node.TypeArray, Array: []*node.Value{{}}}, false},
		{"object with no children", &node.Value{Type: node.TypeObject}, true},
		{
			"object with only an attribute",
			&node.Value{Type: node.TypeObject, Node: &node.Node{Key: "@id", Value: &node.Value{Type: node.TypeString, Worth: "1"}}},
			true,
		},
		{
			"object with a child element",
			&node.Value{Type: node.TypeObject, Node: &node.Node{Key: "a", Value: &node.Value{Type: node.TypeString, Worth: "1"}}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isEmptyElement(tt.value); got != tt.want {
				t.Errorf("isEmptyElement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroupRepeated_NilAndScalar(t *testing.T) {
	// None of these should panic or change anything.
	groupRepeated(nil)
	groupRepeated(&node.Node{Key: "a"})
	groupRepeated(&node.Node{Key: "a", Value: &node.Value{Type: node.TypeString, Worth: "x"}})

	// An array node reached directly recurses into the elements it holds.
	inner := &node.Node{
		Key: "item0",
		Value: &node.Value{
			Type: node.TypeObject,
			Node: &node.Node{
				Key:   "dup",
				Value: &node.Value{Type: node.TypeString, Worth: "1"},
				Next:  &node.Node{Key: "dup", Value: &node.Value{Type: node.TypeString, Worth: "2"}},
			},
		},
	}
	root := &node.Node{
		Key: "root",
		Value: &node.Value{
			Type:  node.TypeArray,
			Array: []*node.Value{{Type: node.TypeObject, Node: inner}},
		},
	}
	groupRepeated(root)

	dup := inner.Value.Node
	if dup == nil || dup.Value.Type != node.TypeArray || len(dup.Value.Array) != 2 {
		t.Errorf("repeated keys inside an array element were not grouped: %+v", dup)
	}
}

func TestNodeToXml_RootArrayDepthLimit(t *testing.T) {
	// A root array whose element is nested past the limit must report it
	// rather than writing a truncated document.
	deep := &node.Value{Type: node.TypeObject}
	current := deep
	for i := 0; i <= maxDepth+2; i++ {
		child := &node.Node{Key: "a", Value: &node.Value{Type: node.TypeObject}}
		current.Node = child
		current = child.Value
	}
	root := &node.Node{
		Key:   "root",
		Value: &node.Value{Type: node.TypeArray, Array: []*node.Value{deep}},
	}
	if _, err := NodeToXml(root); err == nil {
		t.Error("expected an error past the depth limit")
	}
}
