// Package tjson provides functionality for handling JSON data using the Node structure.
// It includes functions for reading, validating, encoding, and decoding JSON data.
package tjson

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/mstgnz/transformer/node"
)

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

// DecodeJson decodes JSON bytes into a Node
func DecodeJson(data []byte) (*node.Node, error) {
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, err
	}

	return interfaceToNode("root", jsonData), nil
}

// interfaceToNode converts an interface{} to a Node
func interfaceToNode(key string, data interface{}) *node.Node {
	if data == nil {
		return &node.Node{
			Key: key,
			Value: &node.Value{
				Type: node.TypeNull,
			},
		}
	}

	switch v := data.(type) {
	case map[string]interface{}:
		n := &node.Node{
			Key: key,
			Value: &node.Value{
				Type: node.TypeObject,
			},
		}

		var last *node.Node
		for k, val := range v {
			child := interfaceToNode(k, val)
			if child != nil {
				child.Parent = n
				if last == nil {
					n.Value.Node = child
				} else {
					last.Next = child
					child.Prev = last
				}
				last = child
			}
		}
		return n

	case []interface{}:
		n := &node.Node{
			Key: key,
			Value: &node.Value{
				Type:  node.TypeArray,
				Array: make([]*node.Value, len(v)),
			},
		}

		for i, item := range v {
			if item == nil {
				n.Value.Array[i] = &node.Value{Type: node.TypeNull}
				continue
			}

			switch val := item.(type) {
			case map[string]interface{}:
				child := interfaceToNode(fmt.Sprintf("item%d", i), val)
				if child != nil {
					n.Value.Array[i] = &node.Value{
						Type: node.TypeObject,
						Node: child,
					}
				}
			case []interface{}:
				child := interfaceToNode(fmt.Sprintf("item%d", i), val)
				if child != nil {
					n.Value.Array[i] = child.Value
				}
			default:
				n.Value.Array[i] = valueFromInterface(val)
			}
		}
		return n

	case string:
		return &node.Node{
			Key: key,
			Value: &node.Value{
				Type:  node.TypeString,
				Worth: v,
			},
		}

	case float64:
		return &node.Node{
			Key: key,
			Value: &node.Value{
				Type:  node.TypeNumber,
				Worth: strconv.FormatFloat(v, 'f', -1, 64),
			},
		}

	case bool:
		return &node.Node{
			Key: key,
			Value: &node.Value{
				Type:  node.TypeBoolean,
				Worth: strconv.FormatBool(v),
			},
		}

	default:
		return &node.Node{
			Key: key,
			Value: &node.Value{
				Type:  node.TypeString,
				Worth: fmt.Sprint(v),
			},
		}
	}
}

// valueFromInterface creates a Value from an interface{}
func valueFromInterface(data interface{}) *node.Value {
	if data == nil {
		return &node.Value{Type: node.TypeNull}
	}

	switch v := data.(type) {
	case string:
		return &node.Value{
			Type:  node.TypeString,
			Worth: v,
		}
	case float64:
		return &node.Value{
			Type:  node.TypeNumber,
			Worth: strconv.FormatFloat(v, 'f', -1, 64),
		}
	case bool:
		return &node.Value{
			Type:  node.TypeBoolean,
			Worth: strconv.FormatBool(v),
		}
	default:
		return &node.Value{
			Type:  node.TypeString,
			Worth: fmt.Sprint(v),
		}
	}
}

// NodeToJson converts a Node to JSON bytes
func NodeToJson(n *node.Node) ([]byte, error) {
	if n == nil {
		return nil, fmt.Errorf("node is nil")
	}

	if n.Value == nil {
		return nil, fmt.Errorf("node value is nil")
	}

	data := nodeToInterface(n)
	return json.Marshal(data)
}

// nodeToInterface converts a Node to a generic interface{}
func nodeToInterface(n *node.Node) interface{} {
	if n == nil || n.Value == nil {
		return nil
	}

	if n.Key == "root" {
		switch n.Value.Type {
		case node.TypeObject:
			result := make(map[string]interface{})
			current := n.Value.Node
			for current != nil {
				if current.Value != nil {
					if current.Value.Type == node.TypeObject || current.Value.Type == node.TypeArray {
						result[current.Key] = nodeToInterface(current)
					} else {
						result[current.Key] = convertValue(current.Value)
					}
				}
				current = current.Next
			}
			return result

		case node.TypeArray:
			result := make([]interface{}, 0)
			for _, item := range n.Value.Array {
				if item != nil {
					if item.Node != nil {
						result = append(result, nodeToInterface(item.Node))
					} else {
						result = append(result, convertValue(item))
					}
				}
			}
			return result

		default:
			return convertValue(n.Value)
		}
	}

	switch n.Value.Type {
	case node.TypeObject:
		result := make(map[string]interface{})
		current := n.Value.Node
		for current != nil {
			if current.Value != nil {
				if current.Value.Type == node.TypeObject || current.Value.Type == node.TypeArray {
					result[current.Key] = nodeToInterface(current)
				} else {
					result[current.Key] = convertValue(current.Value)
				}
			}
			current = current.Next
		}
		return result

	case node.TypeArray:
		result := make([]interface{}, 0)
		for _, item := range n.Value.Array {
			if item != nil {
				if item.Node != nil {
					result = append(result, nodeToInterface(item.Node))
				} else {
					result = append(result, convertValue(item))
				}
			}
		}
		return result

	default:
		return convertValue(n.Value)
	}
}

// convertValue converts a Value to a suitable interface{} type
func convertValue(v *node.Value) interface{} {
	if v == nil {
		return nil
	}

	switch v.Type {
	case node.TypeString:
		return v.Worth
	case node.TypeNumber:
		if f, err := strconv.ParseFloat(v.Worth, 64); err == nil {
			return f
		}
		return v.Worth
	case node.TypeBoolean:
		return v.Worth == "true"
	case node.TypeNull:
		return nil
	default:
		return v.Worth
	}
}
