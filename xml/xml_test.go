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
	byt := []byte(`<root><child>value</child></root>`)
	exp := &node.Node{
		Key: "root",
		Value: &node.Value{
			Node: &node.Node{
				Key: "child",
				Value: &node.Value{
					Worth: "value",
				},
			},
		},
	}

	got, _ := DecodeXml(byt)
	if !reflect.DeepEqual(got.Key, exp.Key) {
		t.Errorf("DecodeXml() got = %v, want %v", got, exp)
	}
	if !reflect.DeepEqual(got.Value.Node.Key, exp.Value.Node.Key) {
		t.Errorf("DecodeXml() got = %v, want %v", got, exp)
	}
	if got.Value.Node.Parent != nil && exp.Value.Node.Parent != nil {
		t.Errorf("DecodeXml() got = %v, want %v", got, exp)
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
	byt := []byte(`<?xml version="1.0" encoding="UTF-8" ?><root attr="root"><apiVersion>v1</apiVersion><kind>Pod</kind><thing prop="2" version="4.5">4.56</thing><metadata><name>rss-site</name><labels><app>web</app></labels></metadata><spec><containers><name>front-end</name><image>nginx</image><ports><containerPort>80</containerPort><port port="34">34</port><port>55</port><status status="on">on</status><status>off</status></ports></containers><containers><name>rss-reader</name><image img="nginx">a/rss-php-nginx:v1</image><ports><containerPort>0.23</containerPort><test>23</test><test>334</test><test type="old">old</test><test>new</test></ports></containers></spec></root>`)

	got, _ := DecodeXml(byt)
	if !reflect.DeepEqual(got.Value.Node.Next.Next.Next.Key, "metadata") {
		t.Errorf("DecodeXml() got = %v, want %v", "metadata", got.Value.Node.Next.Next.Next.Key)
	}
}
