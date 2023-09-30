package transformer

import (
	"reflect"
	"testing"
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
		node *Node
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
		want    *Node
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := XmlDecode(tt.args.byt)
			if (err != nil) != tt.wantErr {
				t.Errorf("XmlDecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("XmlDecode() got = %v, want %v", got, tt.want)
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
			got, err := XmlRead(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("XmlRead() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("XmlRead() got = %v, want %v", got, tt.want)
			}
		})
	}
}
