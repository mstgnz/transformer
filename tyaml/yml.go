// Package tyaml provides functionality for converting between YAML and Node structures.
package tyaml

import (
	"fmt"
	"os"

	"github.com/mstgnz/transformer/node"
	"gopkg.in/yaml.v3"
)

// IsYml checks if the given bytes represent valid YAML
func IsYml(data []byte) bool {
	var out interface{}
	return yaml.Unmarshal(data, &out) == nil
}

// ReadYml reads a YAML file and returns its contents as bytes
func ReadYml(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if !IsYml(data) {
		return nil, fmt.Errorf("invalid YAML format")
	}
	return data, nil
}

// DecodeYml decodes YAML bytes into a Node
func DecodeYml(data []byte) (*node.Node, error) {
	var out interface{}
	if err := yaml.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return interfaceToNode("root", out), nil
}

// NodeToYml converts a Node to YAML bytes
func NodeToYml(n *node.Node) ([]byte, error) {
	if n == nil {
		return nil, fmt.Errorf("node is nil")
	}

	data := nodeToInterface(n)
	return yaml.Marshal(data)
}

// interfaceToNode converts an interface to a Node
func interfaceToNode(key string, v interface{}) *node.Node {
	n := &node.Node{Key: key}

	switch val := v.(type) {
	case map[string]interface{}:
		n.Value = &node.Value{Type: node.TypeObject}
		for k, v := range val {
			child := interfaceToNode(k, v)
			n.AddToEnd(child)
		}

	case []interface{}:
		n.Value = &node.Value{Type: node.TypeArray}
		for _, item := range val {
			switch v := item.(type) {
			case map[string]interface{}, []interface{}:
				child := interfaceToNode("item", v)
				n.Value.Array = append(n.Value.Array, &node.Value{Node: child})
			default:
				n.Value.Array = append(n.Value.Array, &node.Value{
					Type:  getValueType(v),
					Worth: fmt.Sprintf("%v", v),
				})
			}
		}

	case string:
		n.Value = &node.Value{Type: node.TypeString, Worth: val}

	case float64:
		n.Value = &node.Value{Type: node.TypeNumber, Worth: fmt.Sprintf("%v", val)}

	case bool:
		n.Value = &node.Value{Type: node.TypeBoolean, Worth: fmt.Sprintf("%v", val)}

	case nil:
		n.Value = &node.Value{Type: node.TypeNull}

	default:
		n.Value = &node.Value{Type: node.TypeString, Worth: fmt.Sprintf("%v", val)}
	}

	return n
}

// nodeToInterface converts a Node to an interface
func nodeToInterface(n *node.Node) interface{} {
	if n == nil || n.Value == nil {
		return nil
	}

	switch n.Value.Type {
	case node.TypeObject:
		result := make(map[string]interface{})
		current := n.Value.Node
		for current != nil {
			result[current.Key] = nodeToInterface(current)
			current = current.Next
		}
		return result

	case node.TypeArray:
		var result []interface{}
		for _, item := range n.Value.Array {
			if item == nil {
				result = append(result, nil)
			} else if item.Node != nil {
				result = append(result, nodeToInterface(item.Node))
			} else {
				result = append(result, item.Worth)
			}
		}
		return result

	case node.TypeString:
		return n.Value.Worth

	case node.TypeNumber:
		return n.Value.Worth

	case node.TypeBoolean:
		return n.Value.Worth == "true"

	case node.TypeNull:
		return nil

	default:
		return n.Value.Worth
	}
}

// getValueType determines the ValueType of an interface
func getValueType(v interface{}) node.ValueType {
	switch v.(type) {
	case string:
		return node.TypeString
	case float64:
		return node.TypeNumber
	case bool:
		return node.TypeBoolean
	case nil:
		return node.TypeNull
	default:
		return node.TypeString
	}
}
