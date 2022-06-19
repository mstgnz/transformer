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
		key string
		//objStart bool
		arrCount int
		objCount int
		attr     map[string]string
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
			objCount++
		case xml.CharData:
			if objCount <= 0 {
				continue
			}
			val := stripSpaces(string(kind))
			if arrCount > 0 {
				knot.SetToValue(knot, key, val)
			} else {
				if len(attr) > 0 {
					knot = knot.AddToNextWithAttr(knot, parent, key, val, attr)
				} else {
					knot = knot.AddToNext(knot, parent, key, val)
				}
			}
			parent = knot
		case xml.EndElement:
			if arrCount > 0 {
				arrCount--
			}
			if objCount > 0 {
				objCount--
			}
			parent = nil
			if knot.parent != nil {
				knot = knot.parent
				parent = knot
			}
		}
	}
}
