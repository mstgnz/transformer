// Package tjson provides functionality for handling JSON data using the Node structure.
// It includes functions for reading, validating, encoding, and decoding JSON data.
//
// Both directions are driven by tokens rather than by Go maps. A map has no
// order, so routing a document through one silently reshuffles its keys, and
// what comes out the far side is a different document every run. The Node model
// is a linked list precisely so that key order survives; decoding and encoding
// here keep it.
package tjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/mstgnz/transformer/internal/scalar"
	"github.com/mstgnz/transformer/node"
)

const maxDepth = 100

// IsJson checks if the given bytes represent valid JSON
func IsJson(data []byte) bool {
	return json.Valid(data)
}

// ReadJson reads a JSON file and returns its contents as bytes
func ReadJson(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if !IsJson(data) {
		return nil, fmt.Errorf("invalid JSON format")
	}
	return data, nil
}

// DecodeJson decodes JSON bytes into a Node.
//
// Object keys keep the order the document wrote them in, and numbers are
// preserved exactly as written: they are never routed through float64, so a
// 19-digit identifier survives.
func DecodeJson(data []byte) (*node.Node, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()

	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}

	root, err := decodeNode(dec, tok, "root", 0)
	if err != nil {
		return nil, err
	}

	// A second value after the first one means the input is not a single JSON
	// document. Decoding only the first half of a file quietly is worse than
	// saying so.
	if _, err := dec.Token(); err != io.EOF {
		return nil, fmt.Errorf("unexpected data after top-level JSON value")
	}

	return root, nil
}

// decodeNode reads one value and hangs it on a node with the given key.
func decodeNode(dec *json.Decoder, tok json.Token, key string, depth int) (*node.Node, error) {
	n := &node.Node{Key: key}
	value, err := decodeValue(dec, tok, n, depth)
	if err != nil {
		return nil, err
	}
	n.Value = value
	return n, nil
}

// decodeValue reads one value from the token stream. owner is the node the
// value will belong to, so object children can point back at their parent.
func decodeValue(dec *json.Decoder, tok json.Token, owner *node.Node, depth int) (*node.Value, error) {
	if depth > maxDepth {
		return nil, fmt.Errorf("max nesting depth %d exceeded", maxDepth)
	}

	switch t := tok.(type) {
	case json.Delim:
		switch t {
		case '{':
			return decodeObject(dec, owner, depth)
		case '[':
			return decodeArray(dec, depth)
		}
		return nil, fmt.Errorf("unexpected delimiter %q", t)

	case nil:
		return &node.Value{Type: node.TypeNull}, nil

	case string:
		return &node.Value{Type: node.TypeString, Worth: t}, nil

	case bool:
		return &node.Value{Type: node.TypeBoolean, Worth: strconv.FormatBool(t)}, nil

	case json.Number:
		return &node.Value{Type: node.TypeNumber, Worth: t.String()}, nil
	}

	return nil, fmt.Errorf("unsupported JSON token %T", tok)
}

// decodeObject reads members until the closing brace, linking them in order.
func decodeObject(dec *json.Decoder, owner *node.Node, depth int) (*node.Value, error) {
	v := &node.Value{Type: node.TypeObject}
	var last *node.Node

	for {
		keyTok, err := dec.Token()
		if err != nil {
			return nil, err
		}
		if d, ok := keyTok.(json.Delim); ok && d == '}' {
			return v, nil
		}
		key, ok := keyTok.(string)
		if !ok {
			return nil, fmt.Errorf("object key is not a string: %v", keyTok)
		}

		valTok, err := dec.Token()
		if err != nil {
			return nil, err
		}
		child, err := decodeNode(dec, valTok, key, depth+1)
		if err != nil {
			return nil, err
		}

		child.Parent = owner
		if last == nil {
			v.Node = child
		} else {
			last.Next = child
			child.Prev = last
		}
		last = child
	}
}

// decodeArray reads elements until the closing bracket.
func decodeArray(dec *json.Decoder, depth int) (*node.Value, error) {
	v := &node.Value{Type: node.TypeArray, Array: []*node.Value{}}

	for index := 0; ; index++ {
		tok, err := dec.Token()
		if err != nil {
			return nil, err
		}
		if d, ok := tok.(json.Delim); ok && d == ']' {
			return v, nil
		}

		// An object element is wrapped in a node so it keeps a key, which is
		// what the XML writer names the repeated element after. Arrays and
		// scalars are stored as the value itself.
		if d, ok := tok.(json.Delim); ok && d == '{' {
			child, err := decodeNode(dec, tok, fmt.Sprintf("item%d", index), depth+1)
			if err != nil {
				return nil, err
			}
			v.Array = append(v.Array, &node.Value{Type: node.TypeObject, Node: child})
			continue
		}

		item, err := decodeValue(dec, tok, nil, depth+1)
		if err != nil {
			return nil, err
		}
		v.Array = append(v.Array, item)
	}
}

// NodeToJson converts a Node to JSON bytes.
//
// The output is written directly rather than through map[string]any, because
// encoding/json sorts map keys and would drop the document's own order.
func NodeToJson(n *node.Node) ([]byte, error) {
	if n == nil {
		return nil, fmt.Errorf("node is nil")
	}
	if n.Value == nil {
		return nil, fmt.Errorf("node value is nil")
	}

	var buf bytes.Buffer
	if err := writeValue(&buf, n.Value, 0); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// writeValue writes one value in JSON syntax.
func writeValue(buf *bytes.Buffer, v *node.Value, depth int) error {
	if depth > maxDepth {
		return fmt.Errorf("max nesting depth %d exceeded", maxDepth)
	}
	if v == nil {
		buf.WriteString("null")
		return nil
	}

	switch v.Type {
	case node.TypeObject:
		buf.WriteByte('{')
		first := true
		for current := v.Node; current != nil; current = current.Next {
			if !first {
				buf.WriteByte(',')
			}
			first = false
			if err := writeString(buf, current.Key); err != nil {
				return err
			}
			buf.WriteByte(':')
			if err := writeValue(buf, current.Value, depth+1); err != nil {
				return err
			}
		}
		buf.WriteByte('}')

	case node.TypeArray:
		buf.WriteByte('[')
		for i, item := range v.Array {
			if i > 0 {
				buf.WriteByte(',')
			}
			if err := writeValue(buf, itemValue(item), depth+1); err != nil {
				return err
			}
		}
		buf.WriteByte(']')

	case node.TypeString:
		return writeString(buf, v.Worth)

	case node.TypeNumber:
		// A number the format cannot write is not dropped and not corrupted;
		// it goes out as the text it always was.
		if scalar.IsJSONNumber(v.Worth) {
			buf.WriteString(v.Worth)
			return nil
		}
		if normalized, ok := scalar.Number(v.Worth); ok {
			buf.WriteString(normalized)
			return nil
		}
		return writeString(buf, v.Worth)

	case node.TypeBoolean:
		if v.Worth == "true" {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}

	case node.TypeNull:
		buf.WriteString("null")

	default:
		return writeString(buf, v.Worth)
	}

	return nil
}

// itemValue unwraps an array element. Object elements carry a wrapper node
// whose key names the element; everything else is stored directly.
func itemValue(item *node.Value) *node.Value {
	if item == nil {
		return nil
	}
	if item.Node != nil && item.Type == node.TypeObject {
		return item.Node.Value
	}
	return item
}

// writeString writes s as a JSON string, escaped by the standard encoder.
func writeString(buf *bytes.Buffer, s string) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	buf.Write(b)
	return nil
}
