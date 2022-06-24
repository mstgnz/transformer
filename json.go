package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
		value    string
		objStart bool
		arrStart bool
		objCount int
		arrCount int
	)
	for {
		t, err := dec.Token()
		if err == io.EOF || err != nil {
			return knot, errors.Wrap(err, "no more")
		}
		// Get type and value of t object in each loop
		value = fmt.Sprintf("%v", t)
		// If the type of the object is json.Delim
		if reflect.TypeOf(t).String() == "json.Delim" {
			// If no node has been created yet, don't enter here, skip json start -> {
			if !knot.Exists() {
				continue
			}
			// If the value of object t is json object or array
			switch value {
			case "{": // set open object - {
				/*
					There are two types of options here.
					1- If the key is not empty; this is an object. A new node will be added next to the existing node and the newly added node will return.
					2- If the key is empty; There is an array object, and the node will be created in the array and the newly added node will be returned.
				*/
				if arrStart {
					// bu aslında olmazsa olmazdır çünkü bir obje sadece ve sadece array içersinde keysiz başlar.
					knot = knot.AddObjToArr(knot)
				} else {
					knot = knot.AddToNext(knot, parent, key, &node{})
				}
				parent = knot
				objStart = true
				arrStart = false
				objCount++
				key = ""
			case "[": // set open array - [
				/*
					there are three types of options here
					1- If key is not null and objStart is true; The initial value of the node will be set.
					2- If key is not empty and objStart is false; A new node will be created next to the node.
					3- If the key is empty; this is a nested array object. will be added directly to the current node's array.
				*/
				if len(key) > 0 {
					// If objStart is true then the initial value of the node is set.
					if objStart {
						knot = knot.AddToValue(knot, parent, key, []any{})
					} else {
						// If objStart is false, a new node is created next to the current node.
						knot = knot.AddToNext(knot, parent, key, []any{})
					}
					parent = knot
				} else {
					// If there is no key, it is a nested array.
					// TODO nested array için indis tutulacak ve array kapatılana kadar bu arraye append edilecektir.
					//knot = knot.AddToArr(knot, value)
				}
				arrStart = true
				objStart = false
				arrCount++
				key = ""
			case "]": // set close array
				arrCount--
				arrStart = false
				if knot.parent != nil {
					parent = knot.parent
				}
			case "}": // set close object and set parent node
				objCount--
				objStart = false
				if arrCount > 0 {
					arrStart = true
					if knot.parent.parent != nil {
						_, ok := knot.parent.parent.value.([]any)
						if ok {
							knot = knot.parent.parent
						} else {
							knot = knot.parent
						}
					}
				} else {
					parent = nil
					if knot.parent != nil {
						knot = knot.parent
						parent = knot.parent
					}
				}
			default: // shouldn't go here
				log.Println("default not set -> ", t)
			}
		} else {
			// If the loop object is not a json.Delim, the key and value fields will be set.
			// Since the json object is a key value value pair, first the key will be set and then the value will be set.
			if len(key) == 0 {
				// If an array object is open, this key value is essentially an array object.
				// If the array is not empty
				if arrStart {
					knot = knot.AddToArr(knot, value)
				} else {
					key = value
				}
			} else {
				// If objStart is true then the initial value of the node is set.
				if objStart {
					knot = knot.AddToValue(knot, parent, key, value)
				} else {
					knot = knot.AddToNext(knot, parent, key, value)
				}
				// Here objStart and key values are reset.
				objStart = false
				key = ""
			}
		}
	}
}
