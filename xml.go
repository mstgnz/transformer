package main

import (
	"encoding/xml"
	"io"
	"strings"

	"github.com/pkg/errors"
)

func isXml(doc []byte) bool {
	return xml.Unmarshal(doc, new(interface{})) == nil
}

// xmlDecode
func xmlDecode(doc []byte) (*node, error) {
	var (
		knot   *node
		parent *node
	)
	dec := xml.NewDecoder(strings.NewReader(string(doc)))
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
				knot = knot.AddToNextWithAttr(knot, parent, key, &node{}, attr)
			}
			objStart = true
		case xml.CharData:
			if !objStart {
				continue
			}
			val := stripSpaces(string(kind))
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
				knot = knot.AddToNextWithAttr(knot, parent, key, &node{prev: knot}, attr)
				knot = convertToNode(knot.value)
				firstChild = true
			}
			parent = knot
			objStart = false
		case xml.EndElement:
			if arrCount > 0 {
				arrCount--
			}
			parent = nil
			if knot.parent != nil {
				knot = knot.parent
				parent = knot
			}
			//fmt.Println("end", kind.Name.Local)
		}
	}
}
