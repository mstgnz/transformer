package transformer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// IsJSON Checks if the given file is in json format.
func IsJSON(byt []byte) bool {
	return json.Unmarshal(byt, new(map[string]any)) == nil
}

// JsonRead Reads the given file, returns as byt
func JsonRead(filename string) ([]byte, error) {
	byt, err := os.ReadFile(filename)
	if err != nil {
		return byt, errors.Wrap(err, "cannot read the file")
	}
	if isJ := IsJSON(byt); !isJ {
		return byt, errors.Wrap(errors.New("this file is not json"), "this file is not json")
	}
	return byt, nil
}

// JsonDecode Converts a byte array to a key value struct.
func JsonDecode(byt []byte) (*Node, error) {
	var (
		knot   *Node
		parent *Node
	)
	dec := json.NewDecoder(strings.NewReader(string(byt)))
	var (
		key      string
		value    string
		objStart bool
		arrStart bool
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
			// If no Node has been created yet, don't enter here, skip json start -> {
			if !knot.Exists() {
				continue
			}
			// If the value of object t is json object or array
			switch value {
			case "{": // set open object - {
				/*
					1- If the key is not empty; this is an object. A new Node will be added next to the existing Node and the newly added Node will return.
					2- If the key is empty; There is an array object, and the Node will be created in the array and the newly added Node will be returned.
				*/
				if arrStart && len(key) == 0 {
					// An object only starts without a key in an array.
					knot = knot.AddObjToArr(knot)
				} else {
					knot = knot.AddToNext(knot, parent, key, &Node{})
				}
				parent = knot
				objStart = true
				arrStart = false
				key = ""
			case "[": // set open array - [
				/*
					1- If key is not null and objStart is true; The initial value of the Node will be set.
					2- If key is not empty and objStart is false; A new Node will be created next to the Node.
					3- If the key is empty; this is a nested array object. will be added directly to the current Node's array.
				*/
				if len(key) > 0 {
					// If objStart is true then the initial value of the Node is set.
					if objStart {
						knot = knot.AddToValue(knot, parent, key, []any{})
					} else {
						// If objStart is false, a new Node is created next to the current Node.
						knot = knot.AddToNext(knot, parent, key, []any{})
					}
					parent = knot
					arrStart = true
					objStart = false
					arrCount++
					key = ""
				} else {
					// If there is no key, it is a nested array.
					//knot = knot.AddToArr(knot, value)
				}
			case "]": // set close array
				arrCount--
				arrStart = false
			case "}": // set close object and set parent Node
				parent = nil
				if knot.Parent != nil {
					knot = knot.Parent
					parent = knot.Parent
					if arrCount > 0 && len(knot.Key) == 0 {
						arrStart = true
						knot = knot.Parent
						parent = knot.Parent
					}
				}
			default: // shouldn't go here
				log.Println("default not set -> ", t)
			}
		} else {
			// If the loop object is not a json.Delim, the key and value fields will be set.
			// Since the json object is a key value pair, first the key will be set and then the value will be set.
			if len(key) == 0 {
				// If an array object is open, this key value is essentially an array object.
				// If the array is not empty
				if arrStart {
					knot = knot.AddToArr(knot, value)
				} else {
					key = value
				}
			} else {
				// If objStart is true then the initial value of the Node is set.
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

// NodeToJson TODO test
func NodeToJson(node *Node) ([]byte, error) {
	var (
		buf bytes.Buffer
		enc *json.Encoder
	)
	enc = json.NewEncoder(&buf)
	enc.SetIndent("  ", "  ")
	err := enc.Encode(node)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
