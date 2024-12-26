package txml

import (
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/mstgnz/transformer/node"
)

func TestIsXml(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{
			name: "Valid XML",
			data: []byte(`<?xml version="1.0" encoding="UTF-8"?><root><key>value</key></root>`),
			want: true,
		},
		{
			name: "Invalid XML",
			data: []byte(`<root><key>value</key>`),
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
			if got := IsXml(tt.data); got != tt.want {
				t.Errorf("IsXml() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeXml(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    *node.Node
		wantErr bool
	}{
		{
			name: "Simple XML",
			data: []byte(`<?xml version="1.0" encoding="UTF-8"?><root><key>value</key></root>`),
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
			name:    "Invalid XML",
			data:    []byte(`<root><key>value</key>`),
			want:    nil,
			wantErr: true,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := DecodeXml(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeXml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Key != tt.want.Key {
					t.Errorf("DecodeXml() got = %v, want %v", got.Key, tt.want.Key)
				}
				if got.Value.Type != tt.want.Value.Type {
					t.Errorf("DecodeXml() got = %v, want %v", got.Value.Type, tt.want.Value.Type)
				}
			}
		})
	}
}

func TestReadXml(t *testing.T) {
	// Create a temporary test file
	content := []byte(`<?xml version="1.0" encoding="UTF-8"?><root><key>value</key></root>`)
	tmpfile, err := os.CreateTemp("", "test.xml")
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
			name:     "Valid XML file",
			filename: tmpfile.Name(),
			want:     content,
			wantErr:  false,
		},
		{
			name:     "Nonexistent file",
			filename: "nonexistent.xml",
			want:     nil,
			wantErr:  true,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := ReadXml(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadXml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadXml() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeToXml(t *testing.T) {
	tests := []struct {
		name    string
		node    *node.Node
		want    []byte
		wantErr bool
	}{
		{
			name: "Simple XML",
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
			want:    []byte(`<?xml version="1.0" encoding="UTF-8"?><root><key>value</key></root>`),
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
			got, err := NodeToXml(tt.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeToXml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeToXml() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
