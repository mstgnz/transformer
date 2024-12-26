package txml

import (
	"os"
	"strings"
	"testing"
)

func TestIsXml(t *testing.T) {
	validXml, err := os.ReadFile("../example/files/valid.xml")
	if err != nil {
		t.Fatalf("Error reading valid.xml: %v", err)
	}

	invalidXml, err := os.ReadFile("../example/files/invalid.xml")
	if err != nil {
		t.Fatalf("Error reading invalid.xml: %v", err)
	}

	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{
			name: "valid xml",
			data: validXml,
			want: true,
		},
		{
			name: "invalid xml",
			data: invalidXml,
			want: false,
		},
		{
			name: "empty xml",
			data: []byte{},
			want: false,
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
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "valid xml file",
			filename: "../example/files/valid.xml",
			wantErr:  false,
		},
		{
			name:     "invalid xml file",
			filename: "../example/files/invalid.xml",
			wantErr:  true,
		},
		{
			name:     "non-existent file",
			filename: "nonexistent.xml",
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
			if !tt.wantErr && !IsXml(got) {
				t.Error("ReadXml() returned invalid XML")
			}
		})
	}
}

func TestDecodeXml(t *testing.T) {
	validXml, err := os.ReadFile("../example/files/valid.xml")
	if err != nil {
		t.Fatalf("Error reading valid.xml: %v", err)
	}

	invalidXml, err := os.ReadFile("../example/files/invalid.xml")
	if err != nil {
		t.Fatalf("Error reading invalid.xml: %v", err)
	}

	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "valid xml",
			data:    validXml,
			wantErr: false,
		},
		{
			name:    "invalid xml",
			data:    invalidXml,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := DecodeXml(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeXml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Node'u tekrar XML'e dönüştür ve karşılaştır
				xmlBytes, err := NodeToXml(node)
				if err != nil {
					t.Errorf("NodeToXml() error = %v", err)
					return
				}

				// XML'leri normalize et ve karşılaştır
				gotXml := normalizeXml(string(xmlBytes))
				wantXml := normalizeXml(string(tt.data))

				if gotXml != wantXml {
					t.Errorf("XML conversion mismatch\ngot:  %s\nwant: %s", gotXml, wantXml)
				}
			}
		})
	}
}

func TestNodeToXml(t *testing.T) {
	tests := []struct {
		name string
		xml  string
		want string
	}{
		{
			name: "simple element",
			xml:  "<root><test>value</test></root>",
			want: "<root><test>value</test></root>",
		},
		{
			name: "element with attribute",
			xml:  "<root><test attr=\"value\">content</test></root>",
			want: "<root><test attr=\"value\">content</test></root>",
		},
		{
			name: "nested elements",
			xml:  "<root><parent><child>value</child></parent></root>",
			want: "<root><parent><child>value</child></parent></root>",
		},
		{
			name: "multiple attributes",
			xml:  "<root><test a=\"1\" b=\"2\">value</test></root>",
			want: "<root><test a=\"1\" b=\"2\">value</test></root>",
		},
		{
			name: "empty element",
			xml:  "<root><test/></root>",
			want: "<root><test/></root>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// XML'i Node'a dönüştür
			node, err := DecodeXml([]byte(tt.xml))
			if err != nil {
				t.Fatalf("DecodeXml() error = %v", err)
			}

			// Node'u tekrar XML'e dönüştür
			gotXml, err := NodeToXml(node)
			if err != nil {
				t.Fatalf("NodeToXml() error = %v", err)
			}

			// XML'leri normalize et ve karşılaştır
			got := normalizeXml(string(gotXml))
			want := normalizeXml(tt.want)

			if got != want {
				t.Errorf("NodeToXml() = %v, want %v", got, want)
			}
		})
	}
}

// normalizeXml normalizes XML string by removing whitespace and newlines
func normalizeXml(xml string) string {
	// XML başlığını kaldır
	xml = strings.ReplaceAll(xml, `<?xml version="1.0" encoding="UTF-8"?>`, "")
	xml = strings.ReplaceAll(xml, `<?xml version="1.0" encoding="UTF-8" ?>`, "")

	// Boşlukları ve yeni satırları kaldır
	xml = strings.ReplaceAll(xml, "\n", "")
	xml = strings.ReplaceAll(xml, "\r", "")
	xml = strings.ReplaceAll(xml, "\t", "")
	xml = strings.ReplaceAll(xml, "  ", " ")
	xml = strings.TrimSpace(xml)

	return xml
}
