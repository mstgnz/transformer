package transformer

import (
	"reflect"
	"testing"
)

func TestIsYaml(t *testing.T) {
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
			if got := IsYaml(tt.args.byt); got != tt.want {
				t.Errorf("IsYaml() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeToYml(t *testing.T) {
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

func TestYamlDecode(t *testing.T) {
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
			got, err := YamlDecode(tt.args.byt)
			if (err != nil) != tt.wantErr {
				t.Errorf("YamlDecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("YamlDecode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYamlRead(t *testing.T) {
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
			got, err := YamlRead(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("YamlRead() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("YamlRead() got = %v, want %v", got, tt.want)
			}
		})
	}
}
