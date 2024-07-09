package xml

import (
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/mstgnz/transformer/node"
)

func TestIsXml(t *testing.T) {
	type args struct {
		byt []byte
	}
	tests := []struct {
		args args
		want bool
	}{
		{args: args{byt: []byte(`<root><child>value</child></root>`)}, want: true},
		{args: args{byt: []byte(`<root><child>value</child>`)}, want: false},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := IsXml(tt.args.byt); got != tt.want {
				t.Errorf("IsXml() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadXml(t *testing.T) {
	filename := "test.xml"
	expectedContent := `<root><child>value</child></root>`
	err := os.WriteFile(filename, []byte(expectedContent), 0644)
	if err != nil {
		t.Fatalf("Error writing test file: %v", err)
	}
	defer os.Remove(filename)

	type args struct {
		filename string
	}
	tests := []struct {
		args    args
		want    []byte
		wantErr bool
	}{
		{args: args{filename: filename}, want: []byte(expectedContent), wantErr: false},
		{args: args{filename: "nonexistent.xml"}, want: nil, wantErr: true},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := ReadXml(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadXml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadXml() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeXml(t *testing.T) {
	type args struct {
		byt []byte
	}
	tests := []struct {
		args    args
		want    *node.Node
		wantErr bool
	}{
		{
			args: args{byt: []byte(`<root><child>value</child></root>`)},
			want: &node.Node{
				Key: "root",
				Value: &node.Value{
					Node: &node.Node{
						Key: "child",
						Value: &node.Value{
							Worth: "value",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			args:    args{byt: []byte(`<root><child>value</child>`)},
			want:    nil,
			wantErr: true,
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := DecodeXml(tt.args.byt)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeXml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeXml() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeToXml(t *testing.T) {
	type args struct {
		node *node.Node
	}
	tests := []struct {
		args    args
		want    string
		wantErr bool
	}{
		{
			args: args{node: &node.Node{
				Key: "root",
				Value: &node.Value{
					Node: &node.Node{
						Key: "child",
						Value: &node.Value{
							Worth: "value",
						},
					},
				},
			}},
			want:    `<?xml version="1.0" encoding="UTF-8" ?><root><child>value</child></root>`,
			wantErr: false,
		},
		// Add more test cases as needed
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := NodeToXml(tt.args.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeToXml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeToXml() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseXml(t *testing.T) {
	type args struct {
		byt []byte
	}
	tests := []struct {
		args args
		want map[string]any
	}{
		{
			args: args{byt: []byte(`<root><child>value</child></root>`)},
			want: map[string]any{
				"root": map[string]any{
					"child": "value",
				},
			},
		},
		{
			args: args{byt: []byte(`<root><child1>value1</child1><child2>value2</child2></root>`)},
			want: map[string]any{
				"root": map[string]any{
					"child1": "value1",
					"child2": "value2",
				},
			},
		},
		{
			args: args{byt: []byte(`<root><child><subchild>value</subchild></child></root>`)},
			want: map[string]any{
				"root": map[string]any{
					"child": map[string]any{
						"subchild": "value",
					},
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := ParseXml(tt.args.byt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseXml() = %v, want %v", got, tt.want)
			}
		})
	}
}
