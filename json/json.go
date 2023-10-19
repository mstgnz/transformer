package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/mstgnz/transformer/node"
	"github.com/pkg/errors"
)

var (
	Knot   *node.Node
	Parent *node.Node

	key      string
	objStart bool
	arrStart bool
	arrCount int
)

// IsJson
// Checks if the given byte slice is in JSON format.
func IsJson(data []byte) bool {
	return json.Unmarshal(data, new(map[string]any)) == nil
}

// ReadJson
// Reads the given file, returns its content as a byte slice.
func ReadJson(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return data, errors.Wrap(err, "cannot read the file")
	}
	if ok := IsJson(data); !ok {
		return data, errors.Wrap(errors.New("this file is not JSON"), "this file is not JSON")
	}
	return data, nil
}

// DecodeJson
// Converts a byte slice to a `node.Node` structure.
func DecodeJson(data []byte) (*node.Node, error) {
	dec := json.NewDecoder(strings.NewReader(string(data)))
	for {
		token, err := dec.Token()
		if err == io.EOF || err != nil {
			return Knot, errors.Wrap(err, "no more")
		}
		// Get type and value of the token object in each loop
		value := fmt.Sprintf("%v", token)
		// If the type of the object is json.Delim
		if reflect.TypeOf(token).String() == "json.Delim" {
			jsonDelim(token, value)
		} else {
			nonJsonDelim(token, value)
		}
	}
}

// jsonDelim
func jsonDelim(token json.Token, value string) {
	// If no Node has been created yet, don't enter here, skip JSON start -> {
	if !Knot.Exists() {
		return
	}
	// If the value of object token is a JSON object or array
	switch value {
	case "{": // set open object - {
		// If the key is not empty; this is an object. A new Node will be added next to the existing Node, and the newly added Node will return.
		// If the key is empty; There is an array object, and the Node will be created in the array and the newly added Node will be returned.
		if arrStart && len(key) == 0 {
			// An object only starts without a key in an array.
			Knot = Knot.AddToValue(Knot, node.Value{Node: &node.Node{Prev: Knot}})
		} else {
			Knot = Knot.AddToNext(Knot, Parent, key)
		}
		Parent = Knot
		objStart = true
		arrStart = false
		key = ""
	case "[": // set open array - [
		// If the key is not null and objStart is true; The initial value of the Node will be set.
		// If the key is not empty and objStart is false; A new Node will be created next to the Node.
		// If the key is empty; this is a nested array object. will be added directly to the current Node's array.
		if len(key) > 0 {
			// If objStart is true, then the initial value of the Node is set.
			if objStart {
				Knot = Knot.AddToValue(Knot, node.Value{Array: []node.Value{}})
			} else {
				// If objStart is false, a new Node is created next to the current Node.
				Knot = Knot.AddToNext(Knot, Parent, key)
			}
			Parent = Knot
			arrStart = true
			objStart = false
			arrCount++
			key = ""
		} else {
			// If there is no key, it is a nested array.
			_ = append(Knot.Value.Array, node.Value{Worth: value})
		}
	case "]": // set close array
		arrCount--
		arrStart = false
	case "}": // set close object and set parent Node
		Parent = nil
		// if the parent is not empty
		if Knot.Parent != nil {
			// if there is an unclosed array and there is no key
			if arrCount > 0 && len(Knot.Key) == 0 {
				arrStart = true
			}
			// use parent as a node
			Knot = Knot.Parent
			// assign the parent of the node as parent, it can be nil or node
			Parent = Knot.Parent
		}
	default: // shouldn't go here
		log.Println("default not set -> ", token)
	}
}

// nonJsonDelim
func nonJsonDelim(_ json.Token, value string) {
	// If the loop object is not a json.Delim, the key and value fields will be set.
	// Since the JSON object is a key-value pair, first the key will be set and then the value will be set.
	if len(key) == 0 {
		// If an array object is open, this key-value is essentially an array object.
		// If the array is not empty
		if arrStart {
			_ = append(Knot.Value.Array, node.Value{Worth: value})
		} else {
			key = value
		}
	} else {
		// If objStart is true, then the initial value of the Node is set.
		if objStart {
			Knot = Knot.AddToValue(Knot, node.Value{Worth: value})
		} else {
			Knot = Knot.AddToNext(Knot, Parent, key)
		}
		// Here objStart and key values are reset.
		objStart = false
		key = ""
	}
}

// NodeToJson
// TODO test
func NodeToJson(node *node.Node) ([]byte, error) {
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
