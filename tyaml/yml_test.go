package tyaml

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/mstgnz/transformer/node"
	"gopkg.in/yaml.v3"
)

// ---- IsYml / IsYaml ----

// TestIsYml_Alias confirms the deprecated alias delegates correctly.
func TestIsYml_Alias(t *testing.T) {
	if IsYml([]byte("key: value")) != IsYaml([]byte("key: value")) {
		t.Error("IsYml and IsYaml should return the same result")
	}
	if IsYml([]byte{}) != IsYaml([]byte{}) {
		t.Error("IsYml and IsYaml should return the same result for empty input")
	}
}

func TestIsYaml(t *testing.T) {
	validYaml, err := os.ReadFile("../example/files/valid.yaml")
	if err != nil {
		t.Fatalf("Error reading valid.yaml: %v", err)
	}
	invalidYaml, err := os.ReadFile("../example/files/invalid.yaml")
	if err != nil {
		t.Fatalf("Error reading invalid.yaml: %v", err)
	}

	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{"valid yaml", validYaml, true},
		{"invalid yaml", invalidYaml, false},
		{"empty yaml", []byte{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsYaml(tt.data); got != tt.want {
				t.Errorf("IsYaml() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ---- ReadYml / ReadYaml ----

// TestReadYml_Alias confirms the deprecated alias delegates correctly.
func TestReadYml_Alias(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Write([]byte("key: value"))
	tmpfile.Close()

	a, errA := ReadYml(tmpfile.Name())
	b, errB := ReadYaml(tmpfile.Name())
	if errA != nil || errB != nil {
		t.Fatalf("unexpected errors: %v %v", errA, errB)
	}
	if string(a) != string(b) {
		t.Error("ReadYml and ReadYaml should return the same bytes")
	}
}

func TestReadYaml(t *testing.T) {
	content := []byte("name: test")
	tmpfile, err := os.CreateTemp("", "test.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// invalid YAML file
	invalidFile, err := os.CreateTemp("", "invalid*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(invalidFile.Name())
	invalidFile.Write([]byte("key: : bad"))
	invalidFile.Close()

	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{"valid file", tmpfile.Name(), false},
		{"invalid content", invalidFile.Name(), true},
		{"nonexistent file", "no.yaml", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ReadYaml(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadYaml() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ---- DecodeYml / DecodeYaml ----

// TestDecodeYml_Alias confirms the deprecated alias delegates correctly.
func TestDecodeYml_Alias(t *testing.T) {
	input := []byte("key: value")
	a, errA := DecodeYml(input)
	b, errB := DecodeYaml(input)
	if errA != nil || errB != nil {
		t.Fatalf("unexpected errors: %v %v", errA, errB)
	}
	if a.Key != b.Key || a.Value.Type != b.Value.Type {
		t.Error("DecodeYml and DecodeYaml should produce equivalent nodes")
	}
}

func TestDecodeYaml(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		wantKey  string
		wantType node.ValueType
		wantErr  bool
	}{
		{"simple object", []byte("key: value"), "root", node.TypeObject, false},
		{"float value", []byte("ratio: 3.14"), "root", node.TypeObject, false},
		{"invalid YAML", []byte("key: : value"), "", node.TypeNull, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeYaml(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeYaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Key != tt.wantKey {
					t.Errorf("DecodeYaml() key = %v, want %v", got.Key, tt.wantKey)
				}
				if got.Value.Type != tt.wantType {
					t.Errorf("DecodeYaml() type = %v, want %v", got.Value.Type, tt.wantType)
				}
			}
		})
	}
}

// ---- NodeToYaml ----

func TestNodeToYaml(t *testing.T) {
	tests := []struct {
		name    string
		node    *node.Node
		want    string // substring to check in output
		wantErr bool
	}{
		{
			name: "basic string node",
			node: &node.Node{
				Key: "root",
				Value: &node.Value{
					Type: node.TypeObject,
					Node: &node.Node{
						Key:   "key",
						Value: &node.Value{Type: node.TypeString, Worth: "value"},
					},
				},
			},
			want:    "key: value",
			wantErr: false,
		},
		{
			name: "number and boolean",
			node: &node.Node{
				Key: "root",
				Value: &node.Value{
					Type: node.TypeObject,
					Node: &node.Node{
						Key:   "count",
						Value: &node.Value{Type: node.TypeNumber, Worth: "10"},
						Next: &node.Node{
							Key:   "active",
							Value: &node.Value{Type: node.TypeBoolean, Worth: "true"},
						},
					},
				},
			},
			want:    "count:",
			wantErr: false,
		},
		{
			name: "null value",
			node: &node.Node{
				Key: "root",
				Value: &node.Value{
					Type: node.TypeObject,
					Node: &node.Node{
						Key:   "empty",
						Value: &node.Value{Type: node.TypeNull},
					},
				},
			},
			want:    "null",
			wantErr: false,
		},
		{
			name: "array node",
			node: &node.Node{
				Key: "root",
				Value: &node.Value{
					Type: node.TypeObject,
					Node: &node.Node{
						Key: "items",
						Value: &node.Value{
							Type: node.TypeArray,
							Array: []*node.Value{
								{Worth: "a", Type: node.TypeString},
								{Worth: "b", Type: node.TypeString},
							},
						},
					},
				},
			},
			want:    "- a",
			wantErr: false,
		},
		{
			name: "nested object",
			node: &node.Node{
				Key: "root",
				Value: &node.Value{
					Type: node.TypeObject,
					Node: &node.Node{
						Key: "parent",
						Value: &node.Value{
							Type: node.TypeObject,
							Node: &node.Node{
								Key:   "child",
								Value: &node.Value{Type: node.TypeString, Worth: "val"},
							},
						},
					},
				},
			},
			want:    "child: val",
			wantErr: false,
		},
		{
			name:    "nil node",
			node:    nil,
			want:    "",
			wantErr: true,
		},
		{
			name: "root with nil value node",
			node: &node.Node{
				Key:   "root",
				Value: &node.Value{Type: node.TypeObject},
			},
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NodeToYaml(tt.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeToYaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.want != "" {
				if !strings.Contains(string(got), tt.want) {
					t.Errorf("NodeToYaml() output %q does not contain %q", string(got), tt.want)
				}
			}
		})
	}
}

// ---- Round-trip ----

func TestYamlRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		yaml []byte
	}{
		{"simple", []byte("name: alice\n")},
		{"nested", []byte("person:\n  name: bob\n  age: 30\n")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, err := DecodeYml(tt.yaml)
			if err != nil {
				t.Fatalf("DecodeYml() error: %v", err)
			}
			out, err := NodeToYaml(n)
			if err != nil {
				t.Fatalf("NodeToYaml() error: %v", err)
			}
			// Re-decode to verify structural equivalence
			n2, err := DecodeYml(out)
			if err != nil {
				t.Fatalf("DecodeYml() round-trip error: %v", err)
			}
			if n2 == nil {
				t.Fatal("round-trip produced nil node")
			}
		})
	}
}

// ---- interfaceToNode edge cases via DecodeYml ----

func TestDecodeYml_TypeCoverage(t *testing.T) {
	// int value (YAML integers come as int, not float64)
	n, err := DecodeYml([]byte("port: 8080"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	child := n.Value.Node
	if child == nil || child.Value.Type != node.TypeNumber {
		t.Errorf("expected TypeNumber for integer YAML value, got %v", child)
	}

	// float value
	n2, err := DecodeYml([]byte("ratio: 3.14"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	child2 := n2.Value.Node
	if child2 == nil || child2.Value.Type != node.TypeNumber {
		t.Errorf("expected TypeNumber for float, got %v", child2)
	}

	// boolean
	n3, err := DecodeYml([]byte("enabled: true"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	child3 := n3.Value.Node
	if child3 == nil || child3.Value.Type != node.TypeBoolean {
		t.Errorf("expected TypeBoolean, got %v", child3)
	}

	// null
	n4, err := DecodeYml([]byte("empty: null"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	child4 := n4.Value.Node
	if child4 == nil || child4.Value.Type != node.TypeNull {
		t.Errorf("expected TypeNull, got %v", child4)
	}

	// array with primitive items
	n5, err := DecodeYml([]byte("tags:\n  - go\n  - yaml\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	arr := n5.Value.Node
	if arr == nil || arr.Value.Type != node.TypeArray {
		t.Errorf("expected TypeArray, got %v", arr)
	}

	// array with integer items
	n6, err := DecodeYml([]byte("ports:\n  - 80\n  - 443\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	arr6 := n6.Value.Node
	if arr6 == nil || arr6.Value.Type != node.TypeArray {
		t.Errorf("expected TypeArray for int array, got %v", arr6)
	}

	// array with boolean items
	n7, err := DecodeYml([]byte("flags:\n  - true\n  - false\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	arr7 := n7.Value.Node
	if arr7 == nil || arr7.Value.Type != node.TypeArray {
		t.Errorf("expected TypeArray for bool array, got %v", arr7)
	}

	// array with null items
	n8, err := DecodeYml([]byte("list:\n  - null\n  - value\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	arr8 := n8.Value.Node
	if arr8 == nil || arr8.Value.Type != node.TypeArray {
		t.Errorf("expected TypeArray for null array, got %v", arr8)
	}
	if len(arr8.Value.Array) < 1 || arr8.Value.Array[0].Type != node.TypeNull {
		t.Errorf("expected first item to be TypeNull")
	}

	// array with object items
	n9, err := DecodeYml([]byte("users:\n  - name: alice\n  - name: bob\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	arr9 := n9.Value.Node
	if arr9 == nil || arr9.Value.Type != node.TypeArray {
		t.Errorf("expected TypeArray for object array, got %v", arr9)
	}

	// nested array
	n10, err := DecodeYml([]byte("matrix:\n  - - 1\n    - 2\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	arr10 := n10.Value.Node
	if arr10 == nil || arr10.Value.Type != node.TypeArray {
		t.Errorf("expected TypeArray for nested array, got %v", arr10)
	}
}

// ---- NodeToYaml with array of objects ----

func TestNodeToYaml_ArrayOfObjects(t *testing.T) {
	childNode := &node.Node{
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
						{Type: node.TypeObject, Node: childNode},
					},
				},
			},
		},
	}
	out, err := NodeToYaml(n)
	if err != nil {
		t.Fatalf("NodeToYaml() error: %v", err)
	}
	if !strings.Contains(string(out), "alice") {
		t.Errorf("expected 'alice' in output, got: %s", string(out))
	}
}

// ---- scalarToValue: tag to type mapping ----

func TestScalarToValue(t *testing.T) {
	tests := []struct {
		name      string
		tag       string
		value     string
		wantType  node.ValueType
		wantWorth string
	}{
		{"null", "!!null", "null", node.TypeNull, ""},
		{"bool true", "!!bool", "true", node.TypeBoolean, "true"},
		{"bool yes", "!!bool", "yes", node.TypeBoolean, "true"},
		{"bool off", "!!bool", "off", node.TypeBoolean, "false"},
		{"int", "!!int", "42", node.TypeNumber, "42"},
		{"int hex", "!!int", "0x1F", node.TypeNumber, "31"},
		{"int underscores", "!!int", "1_000", node.TypeNumber, "1000"},
		{"big int keeps digits", "!!int", "9007199254740993", node.TypeNumber, "9007199254740993"},
		{"float", "!!float", "3.14", node.TypeNumber, "3.14"},
		{"float exponent kept", "!!float", "1e10", node.TypeNumber, "1e10"},
		{"infinity is not a number", "!!float", ".inf", node.TypeString, ".inf"},
		{"string", "!!str", "hello", node.TypeString, "hello"},
		{"timestamp falls back to string", "!!timestamp", "2026-08-28", node.TypeString, "2026-08-28"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := scalarToValue(&yaml.Node{Kind: yaml.ScalarNode, Tag: tt.tag, Value: tt.value})
			if got.Type != tt.wantType {
				t.Errorf("type = %v, want %v", got.Type, tt.wantType)
			}
			if got.Worth != tt.wantWorth {
				t.Errorf("worth = %q, want %q", got.Worth, tt.wantWorth)
			}
		})
	}
}

// ---- resolveAlias ----

func TestResolveAlias(t *testing.T) {
	target := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "target"}
	alias := &yaml.Node{Kind: yaml.AliasNode, Alias: target}
	if got := resolveAlias(alias); got != target {
		t.Errorf("resolveAlias did not reach the anchored node, got %v", got)
	}

	// A hand-built cycle must fail rather than hang.
	loop := &yaml.Node{Kind: yaml.AliasNode}
	loop.Alias = loop
	if got := resolveAlias(loop); got != nil {
		t.Errorf("resolveAlias on a cyclic alias = %v, want nil", got)
	}

	if got := resolveAlias(nil); got != nil {
		t.Errorf("resolveAlias(nil) = %v, want nil", got)
	}
}

// ---- mappingPairs: order and merge precedence ----

func TestMappingPairs_KeepsDocumentOrder(t *testing.T) {
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte("z: 1\na: 2\nm: 3\n"), &doc); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	pairs, err := mappingPairs(doc.Content[0], 0)
	if err != nil {
		t.Fatalf("mappingPairs error: %v", err)
	}
	want := []string{"z", "a", "m"}
	if len(pairs) != len(want) {
		t.Fatalf("got %d pairs, want %d", len(pairs), len(want))
	}
	for i, key := range want {
		if pairs[i].key != key {
			t.Errorf("pair %d = %q, want %q", i, pairs[i].key, key)
		}
	}
}

func TestMappingPairs_ExplicitKeyBeatsMerge(t *testing.T) {
	input := []byte(`
base: &base
  timeout: 30
  retries: 3
override:
  <<: *base
  timeout: 99
`)
	var doc yaml.Node
	if err := yaml.Unmarshal(input, &doc); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	var override *yaml.Node
	top := doc.Content[0]
	for i := 0; i+1 < len(top.Content); i += 2 {
		if top.Content[i].Value == "override" {
			override = top.Content[i+1]
		}
	}
	if override == nil {
		t.Fatal("override mapping not found")
	}

	pairs, err := mappingPairs(override, 0)
	if err != nil {
		t.Fatalf("mappingPairs error: %v", err)
	}

	got := make(map[string]string, len(pairs))
	for _, p := range pairs {
		got[p.key] = p.value.Value
	}
	if got["timeout"] != "99" {
		t.Errorf("explicit key lost to merge: timeout = %q, want 99", got["timeout"])
	}
	if got["retries"] != "3" {
		t.Errorf("merged key missing: retries = %q, want 3", got["retries"])
	}
	if len(pairs) != 2 {
		t.Errorf("got %d pairs, want 2 (timeout and retries, each once)", len(pairs))
	}
}

func TestMappingPairs_MergeSequenceFirstWins(t *testing.T) {
	input := []byte(`
a: &a
  shared: from-a
b: &b
  shared: from-b
  only-b: yes
merged:
  <<: [*a, *b]
`)
	var doc yaml.Node
	if err := yaml.Unmarshal(input, &doc); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	n, err := DecodeYaml(input)
	if err != nil {
		t.Fatalf("DecodeYaml error: %v", err)
	}

	var merged *node.Node
	for cur := n.Value.Node; cur != nil; cur = cur.Next {
		if cur.Key == "merged" {
			merged = cur
		}
	}
	if merged == nil {
		t.Fatal("merged mapping not found")
	}

	values := make(map[string]string)
	for cur := merged.Value.Node; cur != nil; cur = cur.Next {
		values[cur.Key] = cur.Value.Worth
	}
	if values["shared"] != "from-a" {
		t.Errorf("shared = %q, want from-a (the first source in the merge list)", values["shared"])
	}
	if _, ok := values["only-b"]; !ok {
		t.Error("key only present in the second merge source was dropped")
	}
}

// ---- valueToYaml: tags written for the encoder ----

func TestValueToYaml_Tags(t *testing.T) {
	tests := []struct {
		name    string
		value   *node.Value
		wantTag string
		wantVal string
	}{
		{"nil is null", nil, "!!null", "null"},
		{"null", &node.Value{Type: node.TypeNull}, "!!null", "null"},
		{"string", &node.Value{Type: node.TypeString, Worth: "x"}, "!!str", "x"},
		{"string that looks numeric", &node.Value{Type: node.TypeString, Worth: "42"}, "!!str", "42"},
		{"integer", &node.Value{Type: node.TypeNumber, Worth: "42"}, "!!int", "42"},
		{"float", &node.Value{Type: node.TypeNumber, Worth: "3.14"}, "!!float", "3.14"},
		{"unparseable number becomes text", &node.Value{Type: node.TypeNumber, Worth: "not-a-number"}, "!!str", "not-a-number"},
		{"boolean", &node.Value{Type: node.TypeBoolean, Worth: "true"}, "!!bool", "true"},
		{"boolean from yaml spelling", &node.Value{Type: node.TypeBoolean, Worth: "yes"}, "!!bool", "true"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := valueToYaml(tt.value, 0)
			if err != nil {
				t.Fatalf("valueToYaml error: %v", err)
			}
			if got.Tag != tt.wantTag {
				t.Errorf("tag = %q, want %q", got.Tag, tt.wantTag)
			}
			if got.Value != tt.wantVal {
				t.Errorf("value = %q, want %q", got.Value, tt.wantVal)
			}
		})
	}
}

func TestValueToYaml_DepthLimit(t *testing.T) {
	if _, err := valueToYaml(&node.Value{Type: node.TypeString, Worth: "x"}, maxDepth+1); err == nil {
		t.Error("expected error past the depth limit")
	}
}

// ---- itemValue ----

func TestItemValue(t *testing.T) {
	if got := itemValue(nil); got != nil {
		t.Errorf("itemValue(nil) = %v, want nil", got)
	}

	scalarItem := &node.Value{Type: node.TypeString, Worth: "x"}
	if got := itemValue(scalarItem); got != scalarItem {
		t.Error("a scalar element should be returned as it is stored")
	}

	inner := &node.Value{
		Type: node.TypeObject,
		Node: &node.Node{Key: "k", Value: &node.Value{Type: node.TypeString, Worth: "v"}},
	}
	wrapped := &node.Value{Type: node.TypeObject, Node: &node.Node{Key: "item0", Value: inner}}
	if got := itemValue(wrapped); got != inner {
		t.Error("an object element should be unwrapped to the object it holds")
	}
}

// ---- nodeToInterface non-root path ----

func TestNodeToYaml_NonRootKey(t *testing.T) {
	// NodeToYaml with a node whose Key is NOT "root" exercises convertValue directly
	n := &node.Node{
		Key:   "config",
		Value: &node.Value{Type: node.TypeString, Worth: "val"},
	}
	_, err := NodeToYaml(n)
	// Should not error, but the output may be unexpected since it's not a root node
	if err != nil {
		t.Fatalf("NodeToYaml() error: %v", err)
	}
}

// ---- YAML anchor/alias ----

// YAML anchors and aliases are automatically resolved by yaml.v3 during Unmarshal.
// The decoded values are correct; only the anchor/alias syntax is not preserved.
func TestDecodeYaml_AnchorAlias(t *testing.T) {
	input := []byte(`
defaults: &defaults
  timeout: 30
  retries: 3

production:
  <<: *defaults
  host: prod.example.com
`)
	n, err := DecodeYaml(input)
	if err != nil {
		t.Fatalf("DecodeYaml error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil node")
	}

	// yaml.v3 resolves the merge key (<<: *defaults) into the production map.
	// Find the "production" child and verify "timeout" was merged in.
	var prod *node.Node
	for cur := n.Value.Node; cur != nil; cur = cur.Next {
		if cur.Key == "production" {
			prod = cur
			break
		}
	}
	if prod == nil {
		t.Fatal("expected 'production' key in decoded node")
	}

	var timeout *node.Node
	if prod.Value != nil {
		for cur := prod.Value.Node; cur != nil; cur = cur.Next {
			if cur.Key == "timeout" {
				timeout = cur
				break
			}
		}
	}
	if timeout == nil {
		t.Error("expected 'timeout' to be merged from anchor into production via alias")
	}
}

// ---- nodeToInterface: root with nil value.Node ----

func TestNodeToYaml_RootWithNestedObject(t *testing.T) {
	n := &node.Node{
		Key: "root",
		Value: &node.Value{
			Type: node.TypeObject,
			Node: &node.Node{
				Key: "config",
				Value: &node.Value{
					Type: node.TypeObject,
					Node: &node.Node{
						Key:   "debug",
						Value: &node.Value{Type: node.TypeBoolean, Worth: "false"},
					},
				},
			},
		},
	}
	out, err := NodeToYaml(n)
	if err != nil {
		t.Fatalf("NodeToYaml() error: %v", err)
	}
	if !strings.Contains(string(out), "debug") {
		t.Errorf("expected 'debug' in output, got: %s", string(out))
	}
}

// ---- Output must be readable YAML ----

// TestNodeToYaml_OutputIsValidYaml is the regression test for the hand-rolled
// indentation the writer used to apply. Any document holding a sequence came
// out with the first item at column zero and the rest where the encoder had put
// them, which YAML cannot parse.
func TestNodeToYaml_OutputIsValidYaml(t *testing.T) {
	tests := []struct {
		name string
		yaml string
	}{
		{"flat mapping", "a: 1\nb: x\n"},
		{"one item sequence", "a:\n  - 1\n"},
		{"two item sequence", "a:\n  - 1\n  - 2\n"},
		{"sequence under a nested key", "o:\n  a:\n    - 1\n    - 2\n"},
		{"sequence of mappings", "users:\n  - name: alice\n  - name: bob\n"},
		{"sequence of sequences", "matrix:\n  - - 1\n    - 2\n  - - 3\n"},
		{"deeply nested", "a:\n  b:\n    c:\n      - x\n      - y\n"},
		{"top level sequence", "- 1\n- 2\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, err := DecodeYaml([]byte(tt.yaml))
			if err != nil {
				t.Fatalf("DecodeYaml error: %v", err)
			}
			out, err := NodeToYaml(n)
			if err != nil {
				t.Fatalf("NodeToYaml error: %v", err)
			}

			var reparsed any
			if err := yaml.Unmarshal(out, &reparsed); err != nil {
				t.Fatalf("NodeToYaml produced YAML that cannot be read back: %v\noutput:\n%s", err, out)
			}

			// The document itself must survive, not merely parse.
			var original any
			if err := yaml.Unmarshal([]byte(tt.yaml), &original); err != nil {
				t.Fatalf("test input is not valid YAML: %v", err)
			}
			if !reflect.DeepEqual(reparsed, original) {
				t.Errorf("round trip changed the document\ngot:  %#v\nwant: %#v\noutput:\n%s", reparsed, original, out)
			}
		})
	}
}

// TestNodeToYaml_ArrayScalarsKeepType covers values inside a sequence, which
// used to be written as their raw text and came back as strings.
func TestNodeToYaml_ArrayScalarsKeepType(t *testing.T) {
	n, err := DecodeYaml([]byte("nums:\n  - 1\n  - 2\nflags:\n  - true\n  - false\nwords:\n  - a\n"))
	if err != nil {
		t.Fatalf("DecodeYaml error: %v", err)
	}
	out, err := NodeToYaml(n)
	if err != nil {
		t.Fatalf("NodeToYaml error: %v", err)
	}

	var got map[string][]any
	if err := yaml.Unmarshal(out, &got); err != nil {
		t.Fatalf("output is not valid YAML: %v\n%s", err, out)
	}

	if len(got["nums"]) != 2 {
		t.Fatalf("nums = %#v", got["nums"])
	}
	for i, v := range got["nums"] {
		if _, ok := v.(int); !ok {
			t.Errorf("nums[%d] = %#v, want an integer", i, v)
		}
	}
	for i, v := range got["flags"] {
		if _, ok := v.(bool); !ok {
			t.Errorf("flags[%d] = %#v, want a boolean", i, v)
		}
	}
	if s, ok := got["words"][0].(string); !ok || s != "a" {
		t.Errorf("words[0] = %#v, want the string a", got["words"][0])
	}
}

// TestNodeToYaml_NumericStringStaysQuoted checks the other direction: text that
// looks like a number must not turn into one.
func TestNodeToYaml_NumericStringStaysQuoted(t *testing.T) {
	n := &node.Node{
		Key: "root",
		Value: &node.Value{
			Type: node.TypeObject,
			Node: &node.Node{
				Key:   "zip",
				Value: &node.Value{Type: node.TypeString, Worth: "01234"},
			},
		},
	}
	out, err := NodeToYaml(n)
	if err != nil {
		t.Fatalf("NodeToYaml error: %v", err)
	}

	var got map[string]any
	if err := yaml.Unmarshal(out, &got); err != nil {
		t.Fatalf("output is not valid YAML: %v\n%s", err, out)
	}
	if s, ok := got["zip"].(string); !ok || s != "01234" {
		t.Errorf("zip = %#v, want the string 01234", got["zip"])
	}
}

// ---- Key order ----

func TestDecodeYaml_PreservesKeyOrder(t *testing.T) {
	n, err := DecodeYaml([]byte("zebra: 1\napple: 2\nmango: 3\nbeta: 4\n"))
	if err != nil {
		t.Fatalf("DecodeYaml error: %v", err)
	}

	want := []string{"zebra", "apple", "mango", "beta"}
	var got []string
	for current := n.Value.Node; current != nil; current = current.Next {
		got = append(got, current.Key)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("key order = %v, want %v", got, want)
	}
}

func TestNodeToYaml_OutputIsDeterministic(t *testing.T) {
	input := []byte("zebra: 1\napple: 2\nmango: 3\nbeta: 4\nkilo: 5\n")

	var first string
	for i := 0; i < 25; i++ {
		n, err := DecodeYaml(input)
		if err != nil {
			t.Fatalf("DecodeYaml error: %v", err)
		}
		out, err := NodeToYaml(n)
		if err != nil {
			t.Fatalf("NodeToYaml error: %v", err)
		}
		if i == 0 {
			first = string(out)
			continue
		}
		if string(out) != first {
			t.Fatalf("output changed between runs:\nrun 0:\n%s\nrun %d:\n%s", first, i, out)
		}
	}
	if first != string(input) {
		t.Errorf("output = %q, want the input order back: %q", first, input)
	}
}

// ---- Nested sequences on the way in ----

func TestDecodeYaml_NestedSequence(t *testing.T) {
	n, err := DecodeYaml([]byte("matrix:\n  - - 1\n    - 2\n  - - 3\n"))
	if err != nil {
		t.Fatalf("DecodeYaml error: %v", err)
	}

	matrix := n.Value.Node
	if matrix == nil || matrix.Value.Type != node.TypeArray {
		t.Fatalf("expected an array node, got %v", matrix)
	}
	if len(matrix.Value.Array) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(matrix.Value.Array))
	}
	if matrix.Value.Array[0].Type != node.TypeArray {
		t.Fatalf("expected the first row to stay an array, got %v", matrix.Value.Array[0].Type)
	}
	if len(matrix.Value.Array[0].Array) != 2 {
		t.Errorf("first row = %+v, want 2 values", matrix.Value.Array[0].Array)
	}
}

// ---- Empty input ----

func TestDecodeYaml_Empty(t *testing.T) {
	n, err := DecodeYaml([]byte(""))
	if err != nil {
		t.Fatalf("DecodeYaml error: %v", err)
	}
	if n == nil || n.Value == nil || n.Value.Type != node.TypeNull {
		t.Errorf("empty document = %v, want a null root", n)
	}
}

// ---- yamlToValue edge cases ----

func TestYamlToValue_DocumentNode(t *testing.T) {
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte("key: value\n"), &doc); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	// The document node itself, not its content, to cover the unwrapping branch.
	got, err := yamlToValue(&doc, nil, 0)
	if err != nil {
		t.Fatalf("yamlToValue error: %v", err)
	}
	if got.Type != node.TypeObject || got.Node == nil || got.Node.Key != "key" {
		t.Errorf("document node was not unwrapped: %+v", got)
	}

	empty := &yaml.Node{Kind: yaml.DocumentNode}
	got2, err := yamlToValue(empty, nil, 0)
	if err != nil {
		t.Fatalf("yamlToValue error: %v", err)
	}
	if got2.Type != node.TypeNull {
		t.Errorf("empty document = %v, want null", got2.Type)
	}
}

func TestYamlToValue_UnknownKind(t *testing.T) {
	got, err := yamlToValue(&yaml.Node{}, nil, 0)
	if err != nil {
		t.Fatalf("yamlToValue error: %v", err)
	}
	if got.Type != node.TypeNull {
		t.Errorf("unknown kind = %v, want null", got.Type)
	}
}

func TestYamlToValue_DepthLimit(t *testing.T) {
	n := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "x"}
	if _, err := yamlToValue(n, nil, maxDepth+1); err == nil {
		t.Error("expected an error past the depth limit")
	}
}

func TestDecodeYaml_DepthLimit(t *testing.T) {
	depth := 110
	doc := ""
	for i := 0; i < depth; i++ {
		doc += strings.Repeat(" ", i*2) + "k:\n"
	}
	doc += strings.Repeat(" ", depth*2) + "v\n"
	if _, err := DecodeYaml([]byte(doc)); err == nil {
		t.Error("expected an error for a deeply nested document")
	}
}

func TestSequenceToValue_AliasElements(t *testing.T) {
	input := []byte("base: &base\n  k: v\nlist:\n  - *base\n  - plain\n")
	n, err := DecodeYaml(input)
	if err != nil {
		t.Fatalf("DecodeYaml error: %v", err)
	}

	var list *node.Node
	for cur := n.Value.Node; cur != nil; cur = cur.Next {
		if cur.Key == "list" {
			list = cur
		}
	}
	if list == nil || list.Value.Type != node.TypeArray {
		t.Fatalf("list not found or not an array: %v", list)
	}
	if len(list.Value.Array) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(list.Value.Array))
	}
	if list.Value.Array[0].Node == nil {
		t.Fatal("the aliased element did not resolve to the mapping it names")
	}
	if got := list.Value.Array[0].Node.Value.Node; got == nil || got.Key != "k" {
		t.Errorf("aliased mapping = %v, want the key k", got)
	}
	if list.Value.Array[1].Worth != "plain" {
		t.Errorf("second element = %q, want plain", list.Value.Array[1].Worth)
	}
}

func TestMergeSources_NonMapping(t *testing.T) {
	// A merge key pointing at something that is not a mapping contributes
	// nothing rather than breaking the decode.
	if got := mergeSources(&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "x"}); got != nil {
		t.Errorf("mergeSources(scalar) = %v, want nil", got)
	}
	if got := mergeSources(nil); got != nil {
		t.Errorf("mergeSources(nil) = %v, want nil", got)
	}

	sequence := &yaml.Node{Kind: yaml.SequenceNode, Content: []*yaml.Node{
		{Kind: yaml.MappingNode},
		{Kind: yaml.ScalarNode, Tag: "!!str", Value: "ignored"},
	}}
	if got := mergeSources(sequence); len(got) != 1 {
		t.Errorf("mergeSources(sequence) returned %d mappings, want 1", len(got))
	}
}

func TestNodeToYaml_DepthLimit(t *testing.T) {
	root := &node.Node{Key: "root", Value: &node.Value{Type: node.TypeObject}}
	current := root
	for i := 0; i <= maxDepth+2; i++ {
		child := &node.Node{Key: "a", Value: &node.Value{Type: node.TypeObject}}
		current.Value.Node = child
		child.Parent = current
		current = child
	}
	if _, err := NodeToYaml(root); err == nil {
		t.Error("expected an error past the depth limit")
	}
}
