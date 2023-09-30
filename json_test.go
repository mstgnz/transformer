package transformer

import (
	"reflect"
	"testing"
)

func TestIsJSON(t *testing.T) {
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
			if got := IsJSON(tt.args.byt); got != tt.want {
				t.Errorf("IsJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonDecode(t *testing.T) {
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
			got, err := JsonDecode(tt.args.byt)
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonDecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonDecode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonRead(t *testing.T) {
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
			got, err := JsonRead(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonRead() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonRead() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeToJson(t *testing.T) {
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
			got, err := NodeToJson(tt.args.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeToJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeToJson() got = %v, want %v", got, tt.want)
			}
		})
	}
}
