package tyaml

import (
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/mstgnz/transformer/node"
)

func TestIsYml(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{
			name: "Valid YAML object",
			data: []byte("key: value"),
			want: true,
		},
		{
			name: "Valid YAML array",
			data: []byte("- 1\n- 2\n- 3"),
			want: true,
		},
		{
			name: "Invalid YAML",
			data: []byte("key: : value"),
			want: false,
		},
		{
			name: "Empty input",
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

func TestNodeToYml(t *testing.T) {
	tests := []struct {
		name    string
		node    *node.Node
		want    []byte
		wantErr bool
	}{
		{
			name: "Simple object",
			node: &node.Node{
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
			want:    []byte("key: value\n"),
			wantErr: false,
		},
		{
			name:    "Nil node",
			node:    nil,
			want:    nil,
			wantErr: true,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := NodeToYml(tt.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeToYml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeToYml() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
