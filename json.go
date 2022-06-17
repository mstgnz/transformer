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

func jsonDecode(doc []byte) (IJsonLinear, error) {
	node := JsonLinear()
	dec := json.NewDecoder(strings.NewReader(string(doc)))
	var (
		key        string
		types      string
		typeVal    string
		object     []string
		objectName string
		arrCount   int
	)
	for {
		t, err := dec.Token()
		if err == io.EOF && err != nil {
			return node, errors.Wrap(err, "no more")
		}
		types = reflect.TypeOf(t).String()
		typeVal = fmt.Sprintf("%v", t)
		// get object name
		if len(object) > 0 {
			objectName = object[len(object)-1:][0]
		}
		// json.Delim
		if types == "json.Delim" {
			switch typeVal {
			case "{": // set open object - {
				object = append(object, key)
				if arrCount > 0 {
					node.GetNodeObj(objectName).AppendArr(&jsonLinear{})
					continue
				}
				if objectName == "" {
					node.AddToObj(key)
				} else {
					node.GetNodeObj(objectName).AddToObj(key)
				}
				key = ""
			case "[": // set open array - [
				if key != "" {
					object = append(object, key)
					if objectName == "" {
						node.AddToArr(key)
					} else {
						node.GetNodeObj(objectName).AddToArr(key)
					}
					key = ""
				} else {
					if arrCount > 0 {
						node.GetNodeObj(objectName).AppendArr([]any{})
					}
				}
				arrCount++
			default: // set close object and array - } - ]
				if len(object) > 0 {
					object = object[:len(object)-1]
				}
				if typeVal == "]" {
					arrCount--
				}
			}
			continue
		}
		// isArr
		if arrCount > 0 {
			node.GetNodeObj(objectName).AppendArr(t)
			continue
		}
		// set key
		if key == "" {
			key = typeVal
			continue
		}
		// set val
		if key != "" {
			if objectName == "" {
				node.AddToStr(key, typeVal)
			} else {
				node.GetNodeObj(objectName).AddToStr(key, typeVal)
			}
			key = ""
		}
	}
}
