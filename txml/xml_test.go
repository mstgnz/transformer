package txml

import (
	"os"
	"strings"
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
			name: "valid xml",
			data: []byte(`<root><child>value</child></root>`),
			want: true,
		},
		{
			name: "invalid xml",
			data: []byte(`<root><child>value</child>`),
			want: false,
		},
		{
			name: "empty xml",
			data: []byte(`<root></root>`),
			want: true,
		},
		{
			name: "xml with attributes",
			data: []byte(`<root attr="value"><child>value</child></root>`),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsXml(tt.data); got != tt.want {
				t.Errorf("IsXml() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadXml(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	filename := tmpDir + "/test.xml"
	expectedContent := `<root><child>value</child></root>`
	err := os.WriteFile(filename, []byte(expectedContent), 0644)
	if err != nil {
		t.Fatalf("Error writing test file: %v", err)
	}

	tests := []struct {
		name     string
		filename string
		want     []byte
		wantErr  bool
	}{
		{
			name:     "valid xml file",
			filename: filename,
			want:     []byte(expectedContent),
			wantErr:  false,
		},
		{
			name:     "non-existent file",
			filename: "nonexistent.xml",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "invalid xml file",
			filename: "../example/files/invalid.xml",
			want:     nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadXml(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadXml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && string(got) != string(tt.want) {
				t.Errorf("ReadXml() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestDecodeXml(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "basic xml",
			data:    []byte(`<root><child>value</child></root>`),
			wantErr: false,
		},
		{
			name:    "xml with attributes",
			data:    []byte(`<root attr="value"><child prop="test">value</child></root>`),
			wantErr: false,
		},
		{
			name: "xml with nested elements",
			data: []byte(`
				<root>
					<child>
						<grandchild>value</grandchild>
					</child>
				</root>
			`),
			wantErr: false,
		},
		{
			name:    "xml with numeric values",
			data:    []byte(`<root><number>123.45</number><integer>42</integer></root>`),
			wantErr: false,
		},
		{
			name:    "xml with boolean values",
			data:    []byte(`<root><flag>true</flag><status>false</status></root>`),
			wantErr: false,
		},
		{
			name: "complex xml",
			data: []byte(`
				<root attr="root">
					<apiVersion>v1</apiVersion>
					<kind>Pod</kind>
					<thing prop="2" version="4.5">4.56</thing>
					<metadata>
						<name>rss-site</name>
						<labels>
							<app>web</app>
						</labels>
					</metadata>
					<spec>
						<containers>
							<name>front-end</name>
							<image>nginx</image>
							<ports>
								<containerPort>80</containerPort>
								<port port="34">34</port>
								<status status="on">on</status>
							</ports>
						</containers>
					</spec>
				</root>
			`),
			wantErr: false,
		},
		{
			name:    "invalid xml",
			data:    []byte(`<root><child>value</child>`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeXml(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeXml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Verify the structure is not nil
				if got == nil {
					t.Error("DecodeXml() returned nil Node")
					return
				}

				// Convert back to XML and verify it can be parsed
				xmlBytes, err := NodeToXml(got)
				if err != nil {
					t.Errorf("NodeToXml() error = %v", err)
					return
				}
				if !IsXml(xmlBytes) {
					t.Error("NodeToXml() returned invalid XML")
				}

				// Compare normalized XML strings
				gotStr := normalizeXml(string(xmlBytes))
				wantStr := normalizeXml(string(tt.data))
				if gotStr != wantStr {
					t.Error("Converted XML does not match original")
					t.Errorf("got: %s", gotStr)
					t.Errorf("want: %s", wantStr)
				}
			}
		})
	}
}

func TestNodeToXml(t *testing.T) {
	tests := []struct {
		name    string
		node    *node.Node
		want    string
		wantErr bool
	}{
		{
			name: "basic node",
			node: &node.Node{
				Key: "root",
				Value: &node.Value{
					Type: node.TypeObject,
					Node: &node.Node{
						Key: "child",
						Value: &node.Value{
							Worth: "value",
							Type:  node.TypeString,
						},
					},
				},
			},
			want:    `<?xml version="1.0" encoding="UTF-8"?><root><child>value</child></root>`,
			wantErr: false,
		},
		{
			name: "node with attributes",
			node: &node.Node{
				Key: "root",
				Value: &node.Value{
					Type: node.TypeObject,
					Node: &node.Node{
						Key: "attr",
						Value: &node.Value{
							Type:  node.TypeString,
							Worth: "value",
						},
						Next: &node.Node{
							Key: "child",
							Value: &node.Value{
								Type: node.TypeObject,
								Node: &node.Node{
									Key: "prop",
									Value: &node.Value{
										Type:  node.TypeString,
										Worth: "test",
									},
									Next: &node.Node{
										Key: "content",
										Value: &node.Value{
											Type:  node.TypeString,
											Worth: "content",
										},
									},
								},
							},
						},
					},
				},
			},
			want:    `<?xml version="1.0" encoding="UTF-8"?><root attr="value"><child prop="test">content</child></root>`,
			wantErr: false,
		},
		{
			name: "node with numeric value",
			node: &node.Node{
				Key: "root",
				Value: &node.Value{
					Type: node.TypeObject,
					Node: &node.Node{
						Key: "number",
						Value: &node.Value{
							Worth: "123.45",
							Type:  node.TypeNumber,
						},
					},
				},
			},
			want:    `<?xml version="1.0" encoding="UTF-8"?><root><number>123.45</number></root>`,
			wantErr: false,
		},
		{
			name:    "nil node",
			node:    nil,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NodeToXml(tt.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeToXml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Normalize XML strings for comparison
				gotStr := normalizeXml(string(got))
				wantStr := normalizeXml(tt.want)

				if gotStr != wantStr {
					t.Errorf("NodeToXml() = %v, want %v", gotStr, wantStr)
				}
			}
		})
	}
}

// normalizeXml normalizes XML string by removing whitespace and newlines
func normalizeXml(xml string) string {
	// Remove XML declaration
	xml = strings.ReplaceAll(xml, `<?xml version="1.0" encoding="UTF-8"?>`, "")

	// Remove whitespace between tags
	xml = strings.ReplaceAll(xml, "\n", "")
	xml = strings.ReplaceAll(xml, "\r", "")
	xml = strings.ReplaceAll(xml, "\t", "")

	// Normalize spaces
	parts := strings.Fields(xml)
	return strings.Join(parts, "")
}
