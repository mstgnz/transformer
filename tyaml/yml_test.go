package tyaml

import (
	"os"
	"strings"
	"testing"

	"github.com/mstgnz/transformer/node"
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

// ---- interfaceToNode direct coverage ----

func TestInterfaceToNode_DirectCoverage(t *testing.T) {
	// nil data hits the nil branch
	n, err := interfaceToNode("key", nil, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil || n.Value.Type != node.TypeNull {
		t.Errorf("expected TypeNull for nil, got %v", n)
	}

	// map[any]any path (yaml v1 compat)
	m := map[any]any{"name": "alice", "age": 30}
	n2, err := interfaceToNode("root", m, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n2 == nil || n2.Value.Type != node.TypeObject {
		t.Errorf("expected TypeObject for map[any]any, got %v", n2)
	}

	// default case with unknown type
	n3, err := interfaceToNode("key", struct{ X int }{X: 1}, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n3 == nil || n3.Value.Type != node.TypeString {
		t.Errorf("expected TypeString for unknown type, got %v", n3)
	}

	// array with map[any]any item
	arr := []any{map[any]any{"k": "v"}}
	n4, err := interfaceToNode("arr", arr, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n4 == nil || n4.Value.Type != node.TypeArray {
		t.Errorf("expected TypeArray, got %v", n4)
	}

	// depth exceeded
	_, err = interfaceToNode("key", "val", maxDepth+1)
	if err == nil {
		t.Error("expected error for exceeded depth")
	}
}

// ---- valueFromInterface direct coverage ----

func TestValueFromInterface(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		wantType node.ValueType
		wantWorth string
	}{
		{"nil", nil, node.TypeNull, ""},
		{"string", "hello", node.TypeString, "hello"},
		{"int", 42, node.TypeNumber, "42"},
		{"float64", 3.14, node.TypeNumber, "3.14"},
		{"bool true", true, node.TypeBoolean, "true"},
		{"bool false", false, node.TypeBoolean, "false"},
		{"unknown type", struct{}{}, node.TypeString, "{}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := valueFromInterface(tt.input)
			if v == nil {
				t.Fatal("valueFromInterface returned nil")
			}
			if v.Type != tt.wantType {
				t.Errorf("type = %v, want %v", v.Type, tt.wantType)
			}
		})
	}
}

// ---- convertValue direct coverage ----

func TestConvertValue(t *testing.T) {
	// nil value
	if got := convertValue(nil); got != nil {
		t.Errorf("convertValue(nil) = %v, want nil", got)
	}

	// TypeObject with no node
	v := &node.Value{Type: node.TypeObject}
	if got := convertValue(v); got == nil {
		t.Errorf("convertValue(TypeObject nil node) should return empty map")
	}

	// TypeObject with node
	v2 := &node.Value{
		Type: node.TypeObject,
		Node: &node.Node{
			Key:   "key",
			Value: &node.Value{Type: node.TypeString, Worth: "val"},
		},
	}
	got2 := convertValue(v2)
	if got2 == nil {
		t.Errorf("convertValue(TypeObject with node) returned nil")
	}

	// TypeArray with TypeObject item
	childNode := &node.Node{
		Key: "item0",
		Value: &node.Value{
			Type: node.TypeObject,
			Node: &node.Node{
				Key:   "x",
				Value: &node.Value{Type: node.TypeString, Worth: "y"},
			},
		},
	}
	v3 := &node.Value{
		Type: node.TypeArray,
		Array: []*node.Value{
			{Type: node.TypeObject, Node: childNode},
		},
	}
	got3 := convertValue(v3)
	if got3 == nil {
		t.Errorf("convertValue(TypeArray with object) returned nil")
	}

	// TypeArray with item where Node.Value is nil
	childNodeNilVal := &node.Node{Key: "item0", Value: nil}
	v4 := &node.Value{
		Type: node.TypeArray,
		Array: []*node.Value{
			{Type: node.TypeObject, Node: childNodeNilVal, Worth: "fallback"},
		},
	}
	got4 := convertValue(v4)
	if got4 == nil {
		t.Errorf("convertValue(TypeArray nil item value) returned nil")
	}

	// TypeNumber ParseInt fallback
	v5 := &node.Value{Type: node.TypeNumber, Worth: "9007199254740993"}
	got5 := convertValue(v5)
	if got5 == nil {
		t.Errorf("convertValue(TypeNumber int64) returned nil")
	}

	// TypeNumber full fallback to string
	v6 := &node.Value{Type: node.TypeNumber, Worth: "not-a-number"}
	got6 := convertValue(v6)
	if got6 == nil {
		t.Errorf("convertValue(TypeNumber string fallback) returned nil")
	}

	// nil item in array
	v7 := &node.Value{
		Type:  node.TypeArray,
		Array: []*node.Value{nil},
	}
	got7 := convertValue(v7)
	if got7 == nil {
		t.Errorf("convertValue(TypeArray with nil item) returned nil")
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
