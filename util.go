// Package transformer provides high-level convenience functions for converting
// between JSON, XML, and YAML formats directly, and for searching a document
// without manually handling the intermediate Node representation.
package transformer

import (
	"github.com/mstgnz/transformer/node"
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

// SearchJson finds every node in a JSON document whose key is the given key.
//
// Each match carries the path that reaches it and the depth it sits at, so two
// hits on the same key in different parts of the document can be told apart:
//
//	matches, _ := transformer.SearchJson(data, "name")
//	for _, m := range matches {
//	    fmt.Println(m.Path, m.Depth, m.Node.Value.Worth)
//	}
func SearchJson(data []byte, key string) ([]node.Match, error) {
	n, err := tjson.DecodeJson(data)
	if err != nil {
		return nil, err
	}
	return n.Search(key), nil
}

// SearchXml finds every node in an XML document whose key is the given key.
// Element names, and attribute names written as "@name", are both keys.
func SearchXml(data []byte, key string) ([]node.Match, error) {
	n, err := txml.DecodeXml(data)
	if err != nil {
		return nil, err
	}
	return n.Search(key), nil
}

// SearchYaml finds every node in a YAML document whose key is the given key.
func SearchYaml(data []byte, key string) ([]node.Match, error) {
	n, err := tyaml.DecodeYaml(data)
	if err != nil {
		return nil, err
	}
	return n.Search(key), nil
}

// EditJson decodes a JSON document, hands the tree to edit, and writes the
// result back as JSON.
//
// It is the find-and-replace shape of the convenience layer: everything the
// node package can do to a tree is available inside the callback.
//
//	out, _ := transformer.EditJson(data, func(n *node.Node) {
//	    n.ReplaceText("staging.example.com", "example.com")
//	    n.Rename("labels", "tags")
//	    n.Remove("debug")
//	})
func EditJson(data []byte, edit func(*node.Node)) ([]byte, error) {
	n, err := tjson.DecodeJson(data)
	if err != nil {
		return nil, err
	}
	if edit != nil {
		edit(n)
	}
	return tjson.NodeToJson(n)
}

// EditXml decodes an XML document, hands the tree to edit, and writes the
// result back as XML.
func EditXml(data []byte, edit func(*node.Node)) ([]byte, error) {
	n, err := txml.DecodeXml(data)
	if err != nil {
		return nil, err
	}
	if edit != nil {
		edit(n)
	}
	return txml.NodeToXml(n)
}

// EditYaml decodes a YAML document, hands the tree to edit, and writes the
// result back as YAML.
func EditYaml(data []byte, edit func(*node.Node)) ([]byte, error) {
	n, err := tyaml.DecodeYaml(data)
	if err != nil {
		return nil, err
	}
	if edit != nil {
		edit(n)
	}
	return tyaml.NodeToYaml(n)
}
