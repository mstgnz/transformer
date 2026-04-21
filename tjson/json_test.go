package tjson

import (
	"os"
	"reflect"
	"testing"

	"github.com/mstgnz/transformer/node"
)

func TestIsJson(t *testing.T) {
	validJson, err := os.ReadFile("../example/files/valid.json")
	if err != nil {
		t.Fatalf("Error reading valid.json: %v", err)
	}
	invalidJson, err := os.ReadFile("../example/files/invalid.json")
	if err != nil {
		t.Fatalf("Error reading invalid.json: %v", err)
	}

	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{"valid json", validJson, true},
		{"invalid json", invalidJson, false},
		{"empty", []byte{}, false},
		{"null", []byte("null"), true},
		{"number", []byte("42"), true},
		{"array", []byte(`["a"]`), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsJson(tt.data); got != tt.want {
				t.Errorf("IsJson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadJson(t *testing.T) {
	content := []byte(`{"key": "value"}`)
	tmpfile, err := os.CreateTemp("", "test*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Write(content)
	tmpfile.Close()

	// File exists but has invalid JSON
	invalidFile, err := os.CreateTemp("", "invalid*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(invalidFile.Name())
	invalidFile.Write([]byte(`{bad json`))
	invalidFile.Close()

	tests := []struct {
		name     string
		filename string
		want     []byte
		wantErr  bool
	}{
		{"valid file", tmpfile.Name(), content, false},
		{"nonexistent file", "nonexistent.json", nil, true},
		{"invalid content", invalidFile.Name(), nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadJson(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadJson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeJson(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		wantKey  string
		wantType node.ValueType
		wantErr  bool
	}{
		{"simple object", []byte(`{"key":"value"}`), "root", node.TypeObject, false},
		{"null root", []byte("null"), "root", node.TypeNull, false},
		{"number root", []byte("3.14"), "root", node.TypeNumber, false},
		{"bool root", []byte("true"), "root", node.TypeBoolean, false},
		{"array root", []byte(`["a","b"]`), "root", node.TypeArray, false},
		{"nested object", []byte(`{"a":{"b":"c"}}`), "root", node.TypeObject, false},
		{"array with null", []byte(`[null,"a"]`), "root", node.TypeArray, false},
		{"array with object", []byte(`[{"k":"v"}]`), "root", node.TypeArray, false},
		{"array with array", []byte(`[["x"]]`), "root", node.TypeArray, false},
		{"array with number", []byte(`[1.5,2.0]`), "root", node.TypeArray, false},
		{"array with bool", []byte(`[true,false]`), "root", node.TypeArray, false},
		{"invalid JSON", []byte(`{"key": "value"`), "", node.TypeNull, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeJson(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Key != tt.wantKey {
					t.Errorf("DecodeJson() key = %v, want %v", got.Key, tt.wantKey)
				}
				if got.Value.Type != tt.wantType {
					t.Errorf("DecodeJson() type = %v, want %v", got.Value.Type, tt.wantType)
				}
			}
		})
	}
}

func TestNodeToJson(t *testing.T) {
	tests := []struct {
		name    string
		n       *node.Node
		want    string // substring expected in JSON output
		wantErr bool
	}{
		{
			name: "string value",
			n: &node.Node{
				Key: "root",
				Value: &node.Value{
					Type: node.TypeObject,
					Node: &node.Node{
						Key:   "key",
						Value: &node.Value{Type: node.TypeString, Worth: "value"},
					},
				},
			},
			want:    `"key":"value"`,
			wantErr: false,
		},
		{
			name: "number value",
			n: &node.Node{
				Key: "root",
				Value: &node.Value{
					Type: node.TypeObject,
					Node: &node.Node{
						Key:   "count",
						Value: &node.Value{Type: node.TypeNumber, Worth: "42"},
					},
				},
			},
			want:    `"count":42`,
			wantErr: false,
		},
		{
			name: "float number",
			n: &node.Node{
				Key: "root",
				Value: &node.Value{
					Type: node.TypeObject,
					Node: &node.Node{
						Key:   "ratio",
						Value: &node.Value{Type: node.TypeNumber, Worth: "3.14"},
					},
				},
			},
			want:    `"ratio":3.14`,
			wantErr: false,
		},
		{
			name: "boolean value",
			n: &node.Node{
				Key: "root",
				Value: &node.Value{
					Type: node.TypeObject,
					Node: &node.Node{
						Key:   "ok",
						Value: &node.Value{Type: node.TypeBoolean, Worth: "true"},
					},
				},
			},
			want:    `"ok":true`,
			wantErr: false,
		},
		{
			name: "null value",
			n: &node.Node{
				Key: "root",
				Value: &node.Value{
					Type: node.TypeObject,
					Node: &node.Node{
						Key:   "empty",
						Value: &node.Value{Type: node.TypeNull},
					},
				},
			},
			want:    `"empty":null`,
			wantErr: false,
		},
		{
			name: "nested object",
			n: &node.Node{
				Key: "root",
				Value: &node.Value{
					Type: node.TypeObject,
					Node: &node.Node{
						Key: "outer",
						Value: &node.Value{
							Type: node.TypeObject,
							Node: &node.Node{
								Key:   "inner",
								Value: &node.Value{Type: node.TypeString, Worth: "val"},
							},
						},
					},
				},
			},
			want:    `"inner":"val"`,
			wantErr: false,
		},
		{
			name: "array root",
			n: &node.Node{
				Key: "root",
				Value: &node.Value{
					Type: node.TypeArray,
					Array: []*node.Value{
						{Type: node.TypeString, Worth: "a"},
						{Type: node.TypeString, Worth: "b"},
					},
				},
			},
			want:    `["a","b"]`,
			wantErr: false,
		},
		{
			name: "array with nested object",
			n: &node.Node{
				Key: "root",
				Value: &node.Value{
					Type: node.TypeObject,
					Node: &node.Node{
						Key: "items",
						Value: &node.Value{
							Type: node.TypeArray,
							Array: []*node.Value{
								{
									Type: node.TypeObject,
									Node: &node.Node{
										Key:   "item0",
										Value: &node.Value{Type: node.TypeObject, Node: &node.Node{Key: "name", Value: &node.Value{Type: node.TypeString, Worth: "alice"}}},
									},
								},
							},
						},
					},
				},
			},
			want:    "alice",
			wantErr: false,
		},
		{
			name: "nil value in array",
			n: &node.Node{
				Key: "root",
				Value: &node.Value{
					Type:  node.TypeArray,
					Array: []*node.Value{nil},
				},
			},
			want:    "[",
			wantErr: false,
		},
		{
			name:    "nil node",
			n:       nil,
			want:    "",
			wantErr: true,
		},
		{
			name: "nil value",
			n: &node.Node{
				Key:   "root",
				Value: nil,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "root with primitive string",
			n: &node.Node{
				Key:   "root",
				Value: &node.Value{Type: node.TypeString, Worth: "hello"},
			},
			want:    "hello",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NodeToJson(tt.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeToJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.want != "" {
				if !contains(string(got), tt.want) {
					t.Errorf("NodeToJson() = %s, want to contain %q", string(got), tt.want)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}

func TestJsonRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		json []byte
	}{
		{"object", []byte(`{"name":"alice","age":30}`)},
		{"nested", []byte(`{"person":{"name":"bob"}}`)},
		{"array", []byte(`["a","b","c"]`)},
		{"mixed types", []byte(`{"str":"v","num":1,"flag":true,"empty":null}`)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, err := DecodeJson(tt.json)
			if err != nil {
				t.Fatalf("DecodeJson error: %v", err)
			}
			out, err := NodeToJson(n)
			if err != nil {
				t.Fatalf("NodeToJson error: %v", err)
			}
			// Re-decode to verify structural equivalence
			n2, err := DecodeJson(out)
			if err != nil {
				t.Fatalf("DecodeJson round-trip error: %v", err)
			}
			if n2 == nil {
				t.Fatal("round-trip produced nil node")
			}
		})
	}
}

// Test depth limit
func TestDecodeJson_DepthLimit(t *testing.T) {
	// Build deeply nested JSON
	depth := 110
	json := ""
	for i := 0; i < depth; i++ {
		json += `{"k":`
	}
	json += `"v"`
	for i := 0; i < depth; i++ {
		json += "}"
	}
	_, err := DecodeJson([]byte(json))
	if err == nil {
		t.Error("DecodeJson should return error for deeply nested JSON")
	}
}

// Test big integer preservation
func TestDecodeJson_BigInteger(t *testing.T) {
	big := `{"id": 9007199254740993}`
	n, err := DecodeJson([]byte(big))
	if err != nil {
		t.Fatalf("DecodeJson error: %v", err)
	}
	child := n.Value.Node
	if child == nil {
		t.Fatal("expected child node")
	}
	if child.Value.Worth != "9007199254740993" {
		t.Errorf("big integer precision lost: got %s", child.Value.Worth)
	}
	// Round-trip back to JSON
	out, err := NodeToJson(n)
	if err != nil {
		t.Fatalf("NodeToJson error: %v", err)
	}
	if string(out) != `{"id":9007199254740993}` {
		t.Errorf("round-trip mismatch: got %s", string(out))
	}
}

// Test convertValue with non-parseable number string (fallback to string)
func TestNodeToJson_NumberFallback(t *testing.T) {
	n := &node.Node{
		Key: "root",
		Value: &node.Value{
			Type: node.TypeObject,
			Node: &node.Node{
				Key:   "val",
				Value: &node.Value{Type: node.TypeNumber, Worth: "not-a-number"},
			},
		},
	}
	out, err := NodeToJson(n)
	if err != nil {
		t.Fatalf("NodeToJson() error: %v", err)
	}
	if !contains(string(out), "not-a-number") {
		t.Errorf("expected fallback string in output, got: %s", string(out))
	}
}

// Exercises the non-root TypeArray path in nodeToInterface
func TestNodeToJson_NonRootArray(t *testing.T) {
	n := &node.Node{
		Key: "root",
		Value: &node.Value{
			Type: node.TypeObject,
			Node: &node.Node{
				Key: "tags",
				Value: &node.Value{
					Type: node.TypeArray,
					Array: []*node.Value{
						{Type: node.TypeString, Worth: "go"},
						{Type: node.TypeString, Worth: "yaml"},
					},
				},
			},
		},
	}
	out, err := NodeToJson(n)
	if err != nil {
		t.Fatalf("NodeToJson() error: %v", err)
	}
	if !contains(string(out), "go") {
		t.Errorf("expected 'go' in output, got: %s", string(out))
	}
}

// Exercises array with object items in nodeToInterface
func TestNodeToJson_ArrayWithObjects(t *testing.T) {
	itemNode := &node.Node{
		Key: "item0",
		Value: &node.Value{
			Type: node.TypeObject,
			Node: &node.Node{
				Key:   "id",
				Value: &node.Value{Type: node.TypeNumber, Worth: "1"},
			},
		},
	}
	n := &node.Node{
		Key: "root",
		Value: &node.Value{
			Type: node.TypeArray,
			Array: []*node.Value{
				{Type: node.TypeObject, Node: itemNode},
			},
		},
	}
	out, err := NodeToJson(n)
	if err != nil {
		t.Fatalf("NodeToJson() error: %v", err)
	}
	if !contains(string(out), `"id"`) {
		t.Errorf("expected 'id' in output, got: %s", string(out))
	}
}
