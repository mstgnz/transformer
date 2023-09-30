package xml

import (
	"encoding/xml"
	"io"
	"os"
	"strings"

	"gitgub.com/mstgnz/transformer"
	"gitgub.com/mstgnz/transformer/node"
	"github.com/pkg/errors"
)

// IsXml Checks if the given file is in xml format.
func IsXml(byt []byte) bool {
	return xml.Unmarshal(byt, new(interface{})) == nil
}

// ReadXml Reads the given file, returns as byt
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

// DecodeXml Converts a byte array to a key value struct.
func DecodeXml(byt []byte) (*node.Node, error) {
	var (
		knot   *node.Node
		parent *node.Node
	)
	dec := xml.NewDecoder(strings.NewReader(string(byt)))
	var (
		key        string
		arrCount   int
		objStart   bool
		firstChild bool
		attr       map[string]string
	)

	for {
		t, err := dec.Token()
		if err == io.EOF || err != nil {
			return knot, errors.Wrap(err, "no more")
		}
		switch kind := t.(type) {
		case xml.ProcInst:
			continue
		case xml.StartElement:
			key = kind.Name.Local
			if key == "root" {
				continue
			}
			// Attr
			attr = map[string]string{}
			if len(kind.Attr) > 0 {
				for i := 0; i < len(kind.Attr); i++ {
					attr[kind.Attr[i].Name.Local] = kind.Attr[i].Value
				}
			}
			if objStart {
				knot = knot.AddToNextWithAttr(knot, parent, key, &node.Node{}, attr)
			}
			objStart = true
		case xml.CharData:
			if !objStart {
				continue
			}
			val := transformer.StripSpaces(string(kind))
			if len(val) > 0 {
				if firstChild {
					knot = knot.AddToValueWithAttr(knot, parent, key, val, attr)
				} else {
					if arrCount > 0 {
						knot.SetToValue(knot, key, val)
					} else {
						knot = knot.AddToNextWithAttr(knot, parent, key, val, attr)
					}
				}
			} else {
				knot = knot.AddToNextWithAttr(knot, parent, key, &node.Node{Prev: knot}, attr)
				knot = transformer.ConvertToNode(knot.Value)
				firstChild = true
			}
			parent = knot
			objStart = false
		case xml.EndElement:
			if arrCount > 0 {
				arrCount--
			}
			parent = nil
			if knot.Parent != nil {
				knot = knot.Parent
				parent = knot
			}
		}
	}
}

// NodeToXml TODO
func NodeToXml(node *node.Node) ([]byte, error) {
	return nil, nil
}