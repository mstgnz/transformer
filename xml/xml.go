package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/mstgnz/transformer"
	"github.com/mstgnz/transformer/node"
	"github.com/pkg/errors"
)

// IsXml
// Checks if the given file is in xml format.
func IsXml(byt []byte) bool {
	return xml.Unmarshal(byt, new(any)) == nil
}

// ReadXml
// Reads the given file, returns as byt
func ReadXml(filename string) ([]byte, error) {
	byt, err := os.ReadFile(filename)
	if err != nil {
		return byt, errors.Wrap(err, "cannot read the file")
	}
	if ok := IsXml(byt); !ok {
		return byt, errors.Wrap(errors.New("this file is not xml"), "this file is not xml")
	}
	return byt, nil
}

// DecodeXml
// Converts a byte array to a key value struct.
func DecodeXml(byt []byte) (*node.Node, error) {

	dec := xml.NewDecoder(strings.NewReader(string(byt)))

	var Knot *node.Node
	isNext := false
	start := false

	for {
		t, err := dec.Token()
		if err == io.EOF || err != nil {
			return Knot, errors.Wrap(err, "no more")
		}
		if Knot == nil && reflect.TypeOf(t).Name() != "StartElement" {
			continue
		}

		switch kind := t.(type) {
		case xml.StartElement:
			if Knot == nil {
				Knot = &node.Node{Value: &node.Value{}}
			}
			key := kind.Name.Local
			// Attr
			attr := map[string]string{}
			for _, a := range kind.Attr {
				attr[a.Name.Local] = a.Value
			}
			if isNext {
				Knot.Next = &node.Node{Key: key, Parent: Knot.Parent, Prev: Knot, Value: &node.Value{Attr: attr}}
				Knot = Knot.Next
				isNext = false
			} else {
				if start {
					Knot.Value.Node = &node.Node{Key: key, Parent: Knot, Value: &node.Value{}}
					Knot = Knot.Value.Node
				} else {
					Knot.Key = key
					Knot.Value.Attr = attr
				}
			}
			start = true
		case xml.CharData:
			val := transformer.StripSpaces(string(kind))
			if len(val) > 0 {
				Knot.Value.Worth = val
			}
		case xml.EndElement:
			if Knot.Parent != nil && isNext {
				Knot = Knot.Parent
			}
			isNext = true
			start = false
		}
	}
}

// NodeToXml byte
// TODO nested structure cannot be provided, will be refactored.
func NodeToXml(knot *node.Node) (string, error) {
	var xmlBuilder strings.Builder
	xmlBuilder.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\" ?>")

	var generate func(node *node.Node)
	generate = func(node *node.Node) {
		for node != nil {
			if len(node.Key) > 0 {
				// Write start tag with attributes
				var attrBuilder strings.Builder
				for k, v := range node.Value.Attr {
					attrBuilder.WriteString(fmt.Sprintf(" %v=\"%v\"", k, v))
				}
				xmlBuilder.WriteString(fmt.Sprintf("<%v%v>", node.Key, attrBuilder.String()))

				// Write value if any
				if len(node.Value.Worth) > 0 {
					xmlBuilder.WriteString(node.Value.Worth)
				}

				// Process nested Node
				if node.Value.Node != nil {
					generate(node.Value.Node)
				}

				// Process nested Array
				if len(node.Value.Array) > 0 {
					for _, arrayNode := range node.Value.Array {
						if arrayNode.Node != nil {
							generate(arrayNode.Node)
						}
					}
				}
				// Write end tag
				xmlBuilder.WriteString(fmt.Sprintf("</%v>", node.Key))
			}
			node = node.Next
		}
	}
	generate(knot)

	return xmlBuilder.String(), nil
}

// ParseXml relies only on key and value values. It ignores attribute values and duplicate keys in xml.
func ParseXml(byt []byte) map[string]any {
	key := ""
	end := false
	result := make(map[string]any)
	step := result
	parent := make([]map[string]any, 0)

	dec := xml.NewDecoder(strings.NewReader(string(byt)))
	for {
		t, err := dec.Token()
		if err == io.EOF || err != nil {
			return result
		}
		if len(result) == 0 && reflect.TypeOf(t).Name() != "StartElement" {
			continue
		}
		switch kind := t.(type) {
		case xml.StartElement:
			key = kind.Name.Local
			step[key] = ""
			end = false
		case xml.CharData:
			if !end {
				val := transformer.StripSpaces(string(kind))
				if len(val) > 0 {
					step[key] = val
				} else {
					step[key] = make(map[string]any)
					parent = append(parent, step)
					step = step[key].(map[string]any)
				}
			}
		case xml.EndElement:
			count := len(parent)
			if end && count > 0 {
				step = parent[count-1]
				parent = parent[:count-1]
			}
			end = true
		}
	}
}
