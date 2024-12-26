package tyaml

import (
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/mstgnz/transformer/node"
)

func TestIsYml(t *testing.T) {
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
		{
			name: "valid yaml",
			data: validYaml,
			want: true,
		},
		{
			name: "invalid yaml",
			data: invalidYaml,
			want: false,
		},
		{
			name: "empty yaml",
			data: []byte{},
			want: false,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := IsYml(tt.data); got != tt.want {
				t.Errorf("IsYml() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeYml(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    *node.Node
		wantErr bool
	}{
		{
			name: "Simple object",
			data: []byte("key: value"),
			want: &node.Node{
				Key: "root",
				Value: &node.Value{
					Type: node.TypeObject,
					Node: &node.Node{
						Key: "key",
						Value: &node.Value{
							Type:  node.TypeString,
							Worth: "value",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "Invalid YAML",
			data:    []byte("key: : value"),
			want:    nil,
			wantErr: true,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := DecodeYml(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeYml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Key != tt.want.Key {
					t.Errorf("DecodeYml() got = %v, want %v", got.Key, tt.want.Key)
				}
				if got.Value.Type != tt.want.Value.Type {
					t.Errorf("DecodeYml() got = %v, want %v", got.Value.Type, tt.want.Value.Type)
				}
			}
		})
	}
}

func TestReadYml(t *testing.T) {
	// Create a temporary test file
	content := []byte("key: value")
	tmpfile, err := os.CreateTemp("", "test.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		filename string
		want     []byte
		wantErr  bool
	}{
		{
			name:     "Valid YAML file",
			filename: tmpfile.Name(),
			want:     content,
			wantErr:  false,
		},
		{
			name:     "Nonexistent file",
			filename: "nonexistent.yml",
			want:     nil,
			wantErr:  true,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := ReadYml(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadYml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadYml() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeToYaml(t *testing.T) {
	tests := []struct {
		name    string
		node    *node.Node
		want    string
		wantErr bool
	}{
		{
			name: "basic node",
			node: &node.Node{
				Key: "root",
				Value: &node.Value{
					Type: node.TypeObject,
					Node: &node.Node{
						Key: "key",
						Value: &node.Value{
							Worth: "value",
							Type:  node.TypeString,
						},
					},
				},
			},
			want: `key: value
`,
			wantErr: false,
		},
		{
			name: "node with different types",
			node: &node.Node{
				Key: "root",
				Value: &node.Value{
					Type: node.TypeObject,
					Node: &node.Node{
						Key: "number",
						Value: &node.Value{
							Worth: "123.45",
							Type:  node.TypeNumber,
						},
						Next: &node.Node{
							Key: "boolean",
							Value: &node.Value{
								Worth: "true",
								Type:  node.TypeBoolean,
							},
							Next: &node.Node{
								Key: "null",
								Value: &node.Value{
									Type: node.TypeNull,
								},
							},
						},
					},
				},
			},
			want: `number: 123.45
boolean: true
"null": null
`,
			wantErr: false,
		},
		{
			name: "node with array",
			node: &node.Node{
				Key: "root",
				Value: &node.Value{
					Type: node.TypeObject,
					Node: &node.Node{
						Key: "array",
						Value: &node.Value{
							Type: node.TypeArray,
							Array: []*node.Value{
								{Worth: "item1", Type: node.TypeString},
								{Worth: "item2", Type: node.TypeString},
							},
						},
					},
				},
			},
			want: `array:
- item1
- item2
`,
			wantErr: false,
		},
		{
			name:    "nil node",
			node:    nil,
			want:    "",
			wantErr: true,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := NodeToYaml(tt.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeToYaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeToYaml() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
