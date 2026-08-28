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

// ---- Key order ----

// TestDecodeJson_PreservesKeyOrder pins the reason the Node model is a linked
// list. Routing a document through map[string]any used to reshuffle it.
func TestDecodeJson_PreservesKeyOrder(t *testing.T) {
	n, err := DecodeJson([]byte(`{"zebra":1,"apple":2,"mango":3,"beta":4,"alpha":5}`))
	if err != nil {
		t.Fatalf("DecodeJson error: %v", err)
	}

	want := []string{"zebra", "apple", "mango", "beta", "alpha"}
	var got []string
	for current := n.Value.Node; current != nil; current = current.Next {
		got = append(got, current.Key)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("key order = %v, want %v", got, want)
	}
}

// TestNodeToJson_OutputIsDeterministic guards the same property on the way out:
// encoding/json sorts map keys, so the writer must not go through a map.
func TestNodeToJson_OutputIsDeterministic(t *testing.T) {
	input := []byte(`{"zebra":1,"apple":2,"mango":3,"beta":4,"alpha":5,"kilo":6,"delta":7}`)

	var first string
	for i := 0; i < 25; i++ {
		n, err := DecodeJson(input)
		if err != nil {
			t.Fatalf("DecodeJson error: %v", err)
		}
		out, err := NodeToJson(n)
		if err != nil {
			t.Fatalf("NodeToJson error: %v", err)
		}
		if i == 0 {
			first = string(out)
			continue
		}
		if string(out) != first {
			t.Fatalf("output changed between runs:\nrun 0: %s\nrun %d: %s", first, i, out)
		}
	}
	if first != string(input) {
		t.Errorf("output = %s, want the input order back: %s", first, input)
	}
}

// ---- Round-trip fidelity ----

// TestJsonRoundTrip_Identical checks that a canonical document survives a trip
// through the Node model unchanged, byte for byte.
func TestJsonRoundTrip_Identical(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{"key order", `{"b":2,"a":1,"c":3}`},
		{"scalars", `{"n":null,"t":true,"f":false,"s":"text"}`},
		{"numbers", `{"int":42,"neg":-17,"float":3.14,"exp":1e10}`},
		{"big integers", `{"id":9007199254740993,"huge":12345678901234567890}`},
		{"array of scalars", `{"tags":["go","json","yaml"]}`},
		{"array of objects", `{"users":[{"name":"alice"},{"name":"bob"}]}`},
		{"nested arrays", `{"matrix":[[1,2],[3,4]],"deep":[[[5]]]}`},
		{"mixed array", `{"mixed":[1,"two",true,null,{"k":"v"},[6]]}`},
		{"empty containers", `{"obj":{},"arr":[],"str":""}`},
		{"top level array", `[1,2,3]`},
		{"top level scalar", `"hello"`},
		{"unicode", `{"city":"İstanbul","emoji":"🚀"}`},
		{"deep nesting", `{"a":{"b":{"c":{"d":"e"}}}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, err := DecodeJson([]byte(tt.json))
			if err != nil {
				t.Fatalf("DecodeJson error: %v", err)
			}
			out, err := NodeToJson(n)
			if err != nil {
				t.Fatalf("NodeToJson error: %v", err)
			}
			if string(out) != tt.json {
				t.Errorf("round trip changed the document\ngot:  %s\nwant: %s", out, tt.json)
			}
		})
	}
}

// TestDecodeJson_NestedArrays covers the shape that used to be flattened into
// an empty string on the way out.
func TestDecodeJson_NestedArrays(t *testing.T) {
	n, err := DecodeJson([]byte(`{"matrix":[[1,2],[3]]}`))
	if err != nil {
		t.Fatalf("DecodeJson error: %v", err)
	}

	matrix := n.Value.Node
	if matrix == nil || matrix.Value.Type != node.TypeArray {
		t.Fatalf("expected an array node, got %v", matrix)
	}
	if len(matrix.Value.Array) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(matrix.Value.Array))
	}
	first := matrix.Value.Array[0]
	if first.Type != node.TypeArray {
		t.Fatalf("expected the first row to stay an array, got %v", first.Type)
	}
	if len(first.Array) != 2 || first.Array[0].Worth != "1" || first.Array[1].Worth != "2" {
		t.Errorf("first row lost its values: %+v", first.Array)
	}
}

// ---- Input validation ----

func TestDecodeJson_TrailingData(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
	}{
		{"single document", `{"a":1}`, false},
		{"trailing whitespace", "{\"a\":1}\n  ", false},
		{"two documents", `{"a":1}{"b":2}`, true},
		{"trailing garbage", `{"a":1} nonsense`, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeJson([]byte(tt.json))
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeJson() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNodeToJson_DepthLimit(t *testing.T) {
	root := &node.Node{Key: "root", Value: &node.Value{Type: node.TypeObject}}
	current := root
	for i := 0; i <= maxDepth+2; i++ {
		child := &node.Node{Key: "a", Value: &node.Value{Type: node.TypeObject}}
		current.Value.Node = child
		child.Parent = current
		current = child
	}
	if _, err := NodeToJson(root); err == nil {
		t.Error("expected an error past the depth limit")
	}
}

// TestNodeToJson_EscapesText checks that text which would otherwise break the
// document is escaped rather than written raw.
func TestNodeToJson_EscapesText(t *testing.T) {
	n, err := DecodeJson([]byte(`{"text":"a \"quote\", a \\backslash\\ and a <tag>"}`))
	if err != nil {
		t.Fatalf("DecodeJson error: %v", err)
	}
	out, err := NodeToJson(n)
	if err != nil {
		t.Fatalf("NodeToJson error: %v", err)
	}
	if !IsJson(out) {
		t.Fatalf("output is not valid JSON: %s", out)
	}

	// Decoding the output must give back the same text.
	n2, err := DecodeJson(out)
	if err != nil {
		t.Fatalf("DecodeJson round-trip error: %v", err)
	}
	want := `a "quote", a \backslash\ and a <tag>`
	if got := n2.Value.Node.Value.Worth; got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
}

// ---- Malformed input ----

func TestDecodeJson_Truncated(t *testing.T) {
	tests := []string{
		`{`,
		`{"a"`,
		`{"a":`,
		`{"a":1`,
		`[`,
		`[1`,
		`[1,`,
		`[1,2`,
		`{"a":[1,`,
		`{"a":{"b":`,
	}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			if _, err := DecodeJson([]byte(tt)); err == nil {
				t.Errorf("DecodeJson(%q) returned no error", tt)
			}
		})
	}
}

// ---- Values a hand-built node can hold ----

func TestNodeToJson_UnknownType(t *testing.T) {
	n := &node.Node{
		Key: "root",
		Value: &node.Value{
			Type: node.TypeObject,
			Node: &node.Node{
				Key:   "odd",
				Value: &node.Value{Type: node.ValueType(99), Worth: "value"},
			},
		},
	}
	out, err := NodeToJson(n)
	if err != nil {
		t.Fatalf("NodeToJson error: %v", err)
	}
	if want := `{"odd":"value"}`; string(out) != want {
		t.Errorf("got %s, want %s", out, want)
	}
	if !IsJson(out) {
		t.Errorf("output is not valid JSON: %s", out)
	}
}

// TestNodeToJson_NumberFormsAreNormalized covers numeric text that JSON cannot
// write as it stands. It is converted where that is possible and quoted where
// it is not, so the output parses either way.
func TestNodeToJson_NumberFormsAreNormalized(t *testing.T) {
	tests := []struct {
		name  string
		worth string
		want  string
	}{
		{"plain", "42", `{"n":42}`},
		{"hexadecimal", "0x1F", `{"n":31}`},
		{"leading zeros", "007", `{"n":7}`},
		{"infinity", "Inf", `{"n":"Inf"}`},
		{"text", "not-a-number", `{"n":"not-a-number"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &node.Node{
				Key: "root",
				Value: &node.Value{
					Type: node.TypeObject,
					Node: &node.Node{Key: "n", Value: &node.Value{Type: node.TypeNumber, Worth: tt.worth}},
				},
			}
			out, err := NodeToJson(n)
			if err != nil {
				t.Fatalf("NodeToJson error: %v", err)
			}
			if string(out) != tt.want {
				t.Errorf("got %s, want %s", out, tt.want)
			}
			if !IsJson(out) {
				t.Errorf("output is not valid JSON: %s", out)
			}
		})
	}
}

func TestNodeToJson_NilArrayElement(t *testing.T) {
	n := &node.Node{
		Key: "root",
		Value: &node.Value{
			Type:  node.TypeArray,
			Array: []*node.Value{nil, {Type: node.TypeString, Worth: "x"}},
		},
	}
	out, err := NodeToJson(n)
	if err != nil {
		t.Fatalf("NodeToJson error: %v", err)
	}
	if want := `[null,"x"]`; string(out) != want {
		t.Errorf("got %s, want %s", out, want)
	}
}
