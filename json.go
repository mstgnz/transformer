package main

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

func isJSON(doc []byte) bool {
	return json.Unmarshal(doc, new(map[string]any)) == nil
}

// jsonDecode
func jsonDecode(doc []byte) (*node, error) {
	var (
		knot   *node
		parent *node
	)
	dec := json.NewDecoder(strings.NewReader(string(doc)))
	var (
		key      string
		types    string
		typeVal  string
		objStart bool
		arrCount int
	)
	for {
		t, err := dec.Token()
		if err == io.EOF && err != nil {
			return knot, errors.Wrap(err, "no more")
		}
		types = reflect.TypeOf(t).String()
		typeVal = fmt.Sprintf("%v", t)
		// json.Delim
		if types == "json.Delim" {
			if !knot.Exists() {
				continue
			}
			switch typeVal {
			case "{": // set open object - {
				if arrCount > 0 {
					knot = knot.AddObjToArr(knot)
				} else {
					knot = knot.AddToNext(knot, parent, key, &node{})
				}
				parent = knot
				objStart = true
				key = ""
			case "[": // set open array - [
				if key != "" {
					if objStart {
						knot = knot.AddToValue(knot, parent, key, []any{})
					} else {
						knot = knot.AddToNext(knot, parent, key, []any{})
					}
					key = ""
				} else {
					knot.AddToArr(knot, typeVal)
				}
				arrCount++
			default: // set close object and array - } - ]
				if typeVal == "]" {
					arrCount--
				}
				if typeVal == "}" {
					parent = nil
					if knot.parent != nil {
						knot = knot.parent
						parent = knot
					}
				}
			}
			continue
		} else {
			// set key
			if key == "" {
				if arrCount > 0 {
					knot.AddToArr(knot, typeVal)
					continue
				}
				key = typeVal
				continue
			}
			// set val
			if key != "" {
				if objStart {
					if arrCount > 0 {
						knot.SetToValue(knot, key, typeVal)
					} else {
						knot = knot.AddToValue(knot, parent, key, typeVal)
					}
					objStart = false
				} else {
					knot = knot.AddToNext(knot, parent, key, typeVal)
				}
				key = ""
			}
		}
	}
}
