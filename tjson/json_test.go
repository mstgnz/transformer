package tjson

import (
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/mstgnz/transformer/node"
)

func TestIsJson(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{
			name: "Valid JSON object",
			data: []byte(`{"key": "value"}`),
			want: true,
		},
		{
			name: "Valid JSON array",
			data: []byte(`[1, 2, 3]`),
			want: true,
		},
		{
			name: "Invalid JSON",
			data: []byte(`{"key": "value"`),
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
			if got := IsJson(tt.data); got != tt.want {
				t.Errorf("IsJson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeJson(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    *node.Node
		wantErr bool
	}{
		{
			name: "Simple object",
			data: []byte(`{"key": "value"}`),
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
			name:    "Invalid JSON",
			data:    []byte(`{"key": "value"`),
			want:    nil,
			wantErr: true,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := DecodeJson(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Key != tt.want.Key {
					t.Errorf("DecodeJson() got = %v, want %v", got.Key, tt.want.Key)
				}
				if got.Value.Type != tt.want.Value.Type {
					t.Errorf("DecodeJson() got = %v, want %v", got.Value.Type, tt.want.Value.Type)
				}
			}
		})
	}
}

func TestReadJson(t *testing.T) {
	// Create a temporary test file
	content := []byte(`{"key": "value"}`)
	tmpfile, err := os.CreateTemp("", "test.json")
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
			name:     "Valid JSON file",
			filename: tmpfile.Name(),
			want:     content,
			wantErr:  false,
		},
		{
			name:     "Nonexistent file",
			filename: "nonexistent.json",
			want:     nil,
			wantErr:  true,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
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

func TestNodeToJson(t *testing.T) {
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
			want:    []byte(`{"key":"value"}`),
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
			got, err := NodeToJson(tt.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeToJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeToJson() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
