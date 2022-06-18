package main

import (
	"encoding/xml"
	"fmt"
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
		knot *node
		//parent *node
	)
	dec := xml.NewDecoder(strings.NewReader(string(doc)))
	var (
		key string
		//objStart bool
		//arrCount int
		startEl bool
	)

	for {
		t, err := dec.Token()
		if err == io.EOF && err != nil {
			return knot, errors.Wrap(err, "no more")
		}

		switch kind := t.(type) {
		case xml.StartElement:
			startEl = !startEl
			fmt.Println("start", kind.Name.Local)
		case xml.CharData:
			str := stripSpaces(string([]byte(kind)))
			if len(str) > 0 {
				fmt.Println("char", str)
			}
		case xml.EndElement:
			fmt.Println("end", kind.Name)
		}
		// set key
		if key == "" {
			/*if arrCount > 0 {
				knot.AddToArr(knot, typeVal)
				continue
			}
			key = typeVal
			continue*/
		}
		// set val
		if key != "" {
			/*if objStart {
				if arrCount > 0 {
					knot.SetToValue(knot, key, typeVal)
				} else {
					knot = knot.AddToValue(knot, parent, key, typeVal)
				}
				objStart = false
			} else {
				knot = knot.AddToNext(knot, parent, key, typeVal)
			}*/
			key = ""
		}
	}
}
