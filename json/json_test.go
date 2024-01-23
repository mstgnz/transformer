package json

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/mstgnz/transformer/node"
)

func TestIsJSON(t *testing.T) {
	type args struct {
		byt []byte
	}
	tests := []struct {
		args args
		want bool
	}{
		// TODO: Add test cases.
		{},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := IsJson(tt.args.byt); got != tt.want {
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
		args    args
		want    *node.Node
		wantErr bool
	}{
		// TODO: Add test cases.
		{},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := DecodeJson(tt.args.byt)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeJson() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonRead(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := ReadJson(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadJson() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeToJson(t *testing.T) {
	type args struct {
		node *node.Node
	}
	tests := []struct {
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
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
