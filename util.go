// Package transformer provides high-level convenience functions for converting
// between JSON, XML, and YAML formats directly, without manually handling the
// intermediate Node representation.
package transformer

import (
	"github.com/mstgnz/transformer/tjson"
	"github.com/mstgnz/transformer/txml"
	"github.com/mstgnz/transformer/tyaml"
)

// JsonToXml converts JSON bytes to XML bytes.
func JsonToXml(data []byte) ([]byte, error) {
	n, err := tjson.DecodeJson(data)
	if err != nil {
		return nil, err
	}
	return txml.NodeToXml(n)
}

// JsonToYaml converts JSON bytes to YAML bytes.
func JsonToYaml(data []byte) ([]byte, error) {
	n, err := tjson.DecodeJson(data)
	if err != nil {
		return nil, err
	}
	return tyaml.NodeToYaml(n)
}

// XmlToJson converts XML bytes to JSON bytes.
func XmlToJson(data []byte) ([]byte, error) {
	n, err := txml.DecodeXml(data)
	if err != nil {
		return nil, err
	}
	return tjson.NodeToJson(n)
}

// XmlToYaml converts XML bytes to YAML bytes.
func XmlToYaml(data []byte) ([]byte, error) {
	n, err := txml.DecodeXml(data)
	if err != nil {
		return nil, err
	}
	return tyaml.NodeToYaml(n)
}

// YamlToJson converts YAML bytes to JSON bytes.
func YamlToJson(data []byte) ([]byte, error) {
	n, err := tyaml.DecodeYaml(data)
	if err != nil {
		return nil, err
	}
	return tjson.NodeToJson(n)
}

// YamlToXml converts YAML bytes to XML bytes.
func YamlToXml(data []byte) ([]byte, error) {
	n, err := tyaml.DecodeYaml(data)
	if err != nil {
		return nil, err
	}
	return txml.NodeToXml(n)
}
