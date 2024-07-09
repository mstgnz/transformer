package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"

	"github.com/mstgnz/transformer/node"
	"github.com/pkg/errors"
)

var (
	Knot *node.Node

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
	dec := json.NewDecoder(bytes.NewReader(data))
	for {
		token, err := dec.Token()
		if err == io.EOF || err != nil {
			return Knot, errors.Wrap(err, "no more")
		}
		// Get type and value of the token object in each loop
		value := fmt.Sprintf("%v", token)
		// If the type of the object is json.Delim
		if reflect.TypeOf(token).String() == "json.Delim" {
			jsonDelim(value)
		} else {
			nonJsonDelim(value)
		}
	}
}

// jsonDelim
func jsonDelim(value string) {
	if !Knot.Exists() {
		objStart = true
		Knot = &node.Node{}
		return
	}
	switch value {
	case "{": // set open object - {
		if arrStart && len(key) == 0 {
			sub := &node.Node{}
			Knot.Value.Array = append(Knot.Value.Array, &node.Value{Node: sub})
			sub.Parent = Knot
			Knot = sub
		} else {
			sub := &node.Node{}
			Knot.Next = &node.Node{Key: key, Value: &node.Value{Node: sub}, Prev: Knot, Parent: Knot.Parent}
			sub.Parent = Knot.Next
			Knot = sub
		}
		objStart = true
		arrStart = false
		key = ""
	case "[": // set open array - [
		if len(key) > 0 {
			if objStart {
				Knot.Key = key
				Knot.Value = &node.Value{Array: []*node.Value{}}
			} else {
				Knot.Next = &node.Node{Key: key, Value: &node.Value{}, Parent: Knot.Parent, Prev: Knot}
				Knot = Knot.Next
			}
			arrStart = true
			objStart = false
			arrCount++
			key = ""
		}
	case "]": // set close array
		arrStart = false
		if arrCount > 0 {
			arrCount--
		}
	case "}": // set close object and set parent Node
		objStart = false
		if arrCount > 0 {
			arrStart = true
		}
		if Knot.Parent != nil {
			Knot = Knot.Parent
		}
	default: // shouldn't go here
		log.Println("default not set")
	}
}

// nonJsonDelim
func nonJsonDelim(value string) {
	if len(key) == 0 {
		if arrStart {
			Knot.Value.Array = append(Knot.Value.Array, &node.Value{Worth: value})
		} else {
			key = value
		}
	} else {
		if objStart {
			Knot.Key = key
			Knot.Value = &node.Value{Worth: value}
		} else {
			Knot.Next = &node.Node{Key: key, Parent: Knot.Parent, Prev: Knot}
			Knot.Next.Value = &node.Value{Worth: value}
			Knot = Knot.Next
		}
		objStart = false
		key = ""
	}
}

// NodeToJson
func NodeToJson(node *node.Node) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("  ", "  ")
	err := enc.Encode(node)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
