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
	return xml.Unmarshal(byt, new(interface{})) == nil
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

	var Parent *node.Node
	var Knot *node.Node
	isNext := true

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
			key := kind.Name.Local
			// Attr
			attr := map[string]string{}
			for _, a := range kind.Attr {
				attr[a.Name.Local] = a.Value
			}
			if isNext {
				Knot = Knot.AddToNext(Knot, Parent, key)
				Knot.Value.Attr = attr
				isNext = false
			} else {
				Knot.Key = key
				Knot.Value.Attr = attr
			}
		case xml.CharData:
			val := transformer.StripSpaces(string(kind))
			if len(val) > 0 {
				Knot.Value.Worth = val
			} else {
				Parent = Knot
				Knot.Value.Node = &node.Node{Parent: Parent}
				Knot = Knot.AddToValue(Knot, Knot.Value)
			}
		case xml.EndElement:
			Parent = Knot.Parent
			isNext = true
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
				var attrBuilder strings.Builder
				for k, v := range node.Value.Attr {
					attrBuilder.WriteString(fmt.Sprintf(" %v=\"%v\"", k, v))
				}
				xmlBuilder.WriteString(fmt.Sprintf("<%v%v>%v</%v>", node.Key, attrBuilder.String(), node.Value.Worth, node.Key))
			}
			// if Node Value.Node exists
			if node.Value.Node != nil {
				if len(node.Key) == 0 {
					node.Key = "array"
				}
				generate(node.Value.Node)
			}
			// if Node Value.Array exists
			if len(node.Value.Array) > 0 {
				for _, slc := range node.Value.Array {
					// if Array.Value.Node exists
					if slc.Node != nil {
						generate(slc.Node)
					}
				}
			}
			node = node.Next
		}
	}
	generate(knot.Reset())

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
			}
			end = true
		}
	}
}
