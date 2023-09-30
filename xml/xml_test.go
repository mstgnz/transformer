package xml

import (
	"reflect"
	"testing"

	"gitgub.com/mstgnz/transformer/node"
)

func TestIsXml(t *testing.T) {
	type args struct {
		byt []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsXml(tt.args.byt); got != tt.want {
				t.Errorf("IsXml() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeToXml(t *testing.T) {
	type args struct {
		node *node.Node
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestXmlDecode(t *testing.T) {
	type args struct {
		byt []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *node.Node
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestXmlRead(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
