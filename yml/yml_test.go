package yml

import (
	"reflect"
	"testing"

	"gitgub.com/mstgnz/transformer/node"
)

func TestIsYml(t *testing.T) {
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
			if got := IsYml(tt.args.byt); got != tt.want {
				t.Errorf("IsYml() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeToYml(t *testing.T) {
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
			got, err := NodeToYml(tt.args.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeToYml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeToYml() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYmlDecode(t *testing.T) {
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
			got, err := DecodeYml(tt.args.byt)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeYml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeYml() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYmlRead(t *testing.T) {
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
			got, err := ReadYml(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadYml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadYml() got = %v, want %v", got, tt.want)
			}
		})
	}
}
