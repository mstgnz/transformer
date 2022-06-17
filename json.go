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
		isObj    bool
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
				knot = knot.AddToNext(knot, parent, key, &node{})
				parent = knot
				isObj = !isObj
				key = ""
			case "[": // set open array - [
				if key != "" {
					knot = knot.AddToNext(knot, parent, key, nil)
					key = ""
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
					}
				}
			}
			continue
		} else {
			// set key
			if key == "" {
				key = typeVal
				continue
			}
			// set val
			if key != "" {
				if isObj {
					knot = knot.AddToValue(knot, parent, key, typeVal)
					isObj = !isObj
				} else {
					knot = knot.AddToNext(knot, parent, key, typeVal)
				}
				key = ""
			}
		}
	}
}
