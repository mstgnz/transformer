// Package tyaml provides functionality for handling YAML data using the Node structure.
// It includes functions for reading, validating, encoding, and decoding YAML data.
package tyaml

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mstgnz/transformer/node"
	"gopkg.in/yaml.v3"
)

// IsYaml checks if the given bytes are valid YAML
func IsYaml(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	var yamlData interface{}
	err := yaml.Unmarshal(data, &yamlData)
	if err != nil {
		return false
	}

	// Check if the unmarshaled data is not nil
	return yamlData != nil
}

// ReadYaml reads a YAML file and returns its contents as bytes
func ReadYaml(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if !IsYaml(data) {
		return nil, fmt.Errorf("invalid YAML format")
	}
	return data, nil
}

// DecodeYaml decodes YAML bytes into a Node
func DecodeYaml(data []byte) (*node.Node, error) {
	var yamlData interface{}
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		return nil, err
	}

	return interfaceToNode("root", yamlData), nil
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
	case map[interface{}]interface{}:
		n := &node.Node{
			Key: key,
			Value: &node.Value{
				Type: node.TypeObject,
			},
		}

		var last *node.Node
		for k, val := range v {
			child := interfaceToNode(fmt.Sprint(k), val)
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
			case map[interface{}]interface{}, map[string]interface{}:
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

	case int:
		return &node.Node{
			Key: key,
			Value: &node.Value{
				Type:  node.TypeNumber,
				Worth: strconv.Itoa(v),
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
	case int:
		return &node.Value{
			Type:  node.TypeNumber,
			Worth: strconv.Itoa(v),
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

// NodeToYaml converts a Node to YAML bytes
func NodeToYaml(n *node.Node) ([]byte, error) {
	if n == nil {
		return nil, fmt.Errorf("node is nil")
	}

	if n.Value == nil {
		return nil, fmt.Errorf("node value is nil")
	}

	data := nodeToInterface(n)
	return yaml.Marshal(data)
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
