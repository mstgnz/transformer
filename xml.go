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
		objStart bool
		attr     map[string]string
	)

	for {
		t, err := dec.Token()
		if err == io.EOF || err != nil {
			return knot, errors.Wrap(err, "no more")
		}
		switch kind := t.(type) {
		case xml.StartElement:
			key = kind.Name.Local
			// Attr
			if len(kind.Attr) > 0 {
				attr = map[string]string{}
				for i := 0; i < len(kind.Attr); i++ {
					attr[kind.Attr[i].Name.Local] = kind.Attr[i].Value
				}
				knot = knot.AddToNextWithAttr(knot, parent, key, &node{}, attr)
			} else {
				knot = knot.AddToNext(knot, parent, key, &node{})
			}
			parent = knot
			objStart = true
		case xml.CharData:
			val := stripSpaces(string([]byte(kind)))
			if objStart {
				if arrCount > 0 {
					knot.SetToValue(knot, key, val)
				} else {
					knot = knot.AddToValue(knot, parent, key, val)
				}
				objStart = false
			} else {
				knot = knot.AddToNext(knot, parent, key, val)
			}
			objStart = false
		case xml.EndElement:
			arrCount--
			parent = nil
			if knot.parent != nil {
				knot = knot.parent
				parent = knot
			}
			//fmt.Println("end", kind.Name.Local)
		}
	}
}
