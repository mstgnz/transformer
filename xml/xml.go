package xml

import (
	"encoding/xml"
	"io"
	"os"
	"reflect"
	"strings"

	"gitgub.com/mstgnz/transformer"
	"gitgub.com/mstgnz/transformer/node"
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
	var (
		Knot   *node.Node
		Parent *node.Node
		isNext bool
	)
	dec := xml.NewDecoder(strings.NewReader(string(byt)))

	for {
		t, err := dec.Token()
		if err == io.EOF || err != nil {
			return Knot, errors.Wrap(err, "no more")
		}
		if Knot == nil && reflect.TypeOf(t).Name() != "StartElement" {
			continue
		}
		Knot = &node.Node{}
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
				Knot = Knot.AddToValue(Knot, node.Value{Node: &node.Node{Parent: Parent}})
			}
		case xml.EndElement:
			Parent = Knot.Parent
			isNext = true
		}
	}
}

// NodeToXml
// TODO implement
func NodeToXml(node *node.Node) ([]byte, error) {
	return nil, nil
}
