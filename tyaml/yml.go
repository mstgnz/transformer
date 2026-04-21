// Package tyaml provides functionality for handling YAML data using the Node structure.
// It includes functions for reading, validating, encoding, and decoding YAML data.
package tyaml

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mstgnz/transformer/node"
	"gopkg.in/yaml.v3"
)

const maxDepth = 100

// IsYaml checks if the given bytes are valid YAML
func IsYaml(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	var yamlData any
	err := yaml.Unmarshal(data, &yamlData)
	if err != nil {
		return false
	}
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
	var yamlData any
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		return nil, err
	}
	return interfaceToNode("root", yamlData, 0)
}

// IsYml is an alias for IsYaml.
//
// Deprecated: use IsYaml.
func IsYml(data []byte) bool { return IsYaml(data) }

// ReadYml is an alias for ReadYaml.
//
// Deprecated: use ReadYaml.
func ReadYml(filename string) ([]byte, error) { return ReadYaml(filename) }

// DecodeYml is an alias for DecodeYaml.
//
// Deprecated: use DecodeYaml.
func DecodeYml(data []byte) (*node.Node, error) { return DecodeYaml(data) }

// interfaceToNode converts an any to a Node with depth tracking.
func interfaceToNode(key string, data any, depth int) (*node.Node, error) {
	if depth > maxDepth {
		return nil, fmt.Errorf("max nesting depth %d exceeded", maxDepth)
	}

	if data == nil {
		return &node.Node{
			Key:   key,
			Value: &node.Value{Type: node.TypeNull},
		}, nil
	}

	switch v := data.(type) {
	case map[any]any:
		n := &node.Node{
			Key:   key,
			Value: &node.Value{Type: node.TypeObject},
		}
		var last *node.Node
		for k, val := range v {
			child, err := interfaceToNode(fmt.Sprint(k), val, depth+1)
			if err != nil {
				return nil, err
			}
			child.Parent = n
			if last == nil {
				n.Value.Node = child
			} else {
				last.Next = child
				child.Prev = last
			}
			last = child
		}
		return n, nil

	case map[string]any:
		n := &node.Node{
			Key:   key,
			Value: &node.Value{Type: node.TypeObject},
		}
		var last *node.Node
		for k, val := range v {
			child, err := interfaceToNode(k, val, depth+1)
			if err != nil {
				return nil, err
			}
			child.Parent = n
			if last == nil {
				n.Value.Node = child
			} else {
				last.Next = child
				child.Prev = last
			}
			last = child
		}
		return n, nil

	case []any:
		n := &node.Node{
			Key:   key,
			Value: &node.Value{Type: node.TypeArray, Array: make([]*node.Value, len(v))},
		}
		for i, item := range v {
			if item == nil {
				n.Value.Array[i] = &node.Value{Type: node.TypeNull}
				continue
			}
			switch val := item.(type) {
			case map[any]any, map[string]any:
				child, err := interfaceToNode(fmt.Sprintf("item%d", i), val, depth+1)
				if err != nil {
					return nil, err
				}
				n.Value.Array[i] = &node.Value{Type: node.TypeObject, Node: child}
			case []any:
				child, err := interfaceToNode(fmt.Sprintf("item%d", i), val, depth+1)
				if err != nil {
					return nil, err
				}
				n.Value.Array[i] = child.Value
			default:
				n.Value.Array[i] = valueFromInterface(val)
			}
		}
		return n, nil

	case string:
		return &node.Node{
			Key:   key,
			Value: &node.Value{Type: node.TypeString, Worth: v},
		}, nil

	case int:
		return &node.Node{
			Key:   key,
			Value: &node.Value{Type: node.TypeNumber, Worth: strconv.Itoa(v)},
		}, nil

	case float64:
		return &node.Node{
			Key:   key,
			Value: &node.Value{Type: node.TypeNumber, Worth: strconv.FormatFloat(v, 'f', -1, 64)},
		}, nil

	case bool:
		return &node.Node{
			Key:   key,
			Value: &node.Value{Type: node.TypeBoolean, Worth: strconv.FormatBool(v)},
		}, nil

	default:
		return &node.Node{
			Key:   key,
			Value: &node.Value{Type: node.TypeString, Worth: fmt.Sprint(v)},
		}, nil
	}
}

// valueFromInterface creates a Value from an any
func valueFromInterface(data any) *node.Value {
	if data == nil {
		return &node.Value{Type: node.TypeNull}
	}
	switch v := data.(type) {
	case string:
		return &node.Value{Type: node.TypeString, Worth: v}
	case int:
		return &node.Value{Type: node.TypeNumber, Worth: strconv.Itoa(v)}
	case float64:
		return &node.Value{Type: node.TypeNumber, Worth: strconv.FormatFloat(v, 'f', -1, 64)}
	case bool:
		return &node.Value{Type: node.TypeBoolean, Worth: strconv.FormatBool(v)}
	default:
		return &node.Value{Type: node.TypeString, Worth: fmt.Sprint(v)}
	}
}

// NodeToYaml converts a Node to YAML bytes
func NodeToYaml(n *node.Node) ([]byte, error) {
	if n == nil {
		return nil, fmt.Errorf("node is nil")
	}

	data := nodeToInterface(n)

	yamlBytes, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Normalize: remove empty lines and strip indentation from array items
	lines := strings.Split(string(yamlBytes), "\n")
	var formattedLines []string
	var prevWasKey bool

	for _, line := range lines {
		trimmedLine := strings.TrimRight(line, " \t\n\r")
		if trimmedLine == "" {
			continue
		}

		stripped := strings.TrimLeft(trimmedLine, " \t")
		if prevWasKey && strings.HasPrefix(stripped, "- ") {
			formattedLines = append(formattedLines, stripped)
		} else {
			formattedLines = append(formattedLines, trimmedLine)
		}

		prevWasKey = strings.HasSuffix(stripped, ":")
	}

	if len(formattedLines) == 0 {
		return []byte{}, nil
	}

	return []byte(strings.Join(formattedLines, "\n") + "\n"), nil
}

// nodeToInterface converts a Node to interface{}
func nodeToInterface(n *node.Node) any {
	if n == nil {
		return nil
	}

	if n.Key == "root" {
		if n.Value == nil {
			return nil
		}
		if n.Value.Type == node.TypeArray {
			return convertValue(n.Value)
		}
		if n.Value.Node == nil {
			return nil
		}
		return nodeToMap(n.Value.Node)
	}

	return convertValue(n.Value)
}

// nodeToMap converts a Node to map[string]any
func nodeToMap(n *node.Node) map[string]any {
	result := make(map[string]any)
	current := n

	for current != nil {
		if current.Value != nil {
			switch current.Value.Type {
			case node.TypeArray:
				result[current.Key] = convertValue(current.Value)
			case node.TypeObject:
				if current.Value.Node != nil {
					result[current.Key] = nodeToMap(current.Value.Node)
				} else {
					result[current.Key] = make(map[string]any)
				}
			default:
				result[current.Key] = convertValue(current.Value)
			}
		}
		current = current.Next
	}

	return result
}

// convertValue converts a Value to interface{}.
// Integers are tried first to avoid float notation in YAML output.
func convertValue(v *node.Value) any {
	if v == nil {
		return nil
	}

	switch v.Type {
	case node.TypeNull:
		return nil
	case node.TypeString:
		return v.Worth
	case node.TypeNumber:
		if i, err := strconv.ParseInt(v.Worth, 10, 64); err == nil {
			return i
		}
		if u, err := strconv.ParseUint(v.Worth, 10, 64); err == nil {
			return u
		}
		if f, err := strconv.ParseFloat(v.Worth, 64); err == nil {
			return f
		}
		return v.Worth
	case node.TypeBoolean:
		b, _ := strconv.ParseBool(v.Worth)
		return b
	case node.TypeArray:
		result := make([]any, 0)
		for _, item := range v.Array {
			if item != nil {
				switch item.Type {
				case node.TypeObject:
					if item.Node != nil {
						if item.Node.Value != nil {
							result = append(result, convertValue(item.Node.Value))
						} else {
							result = append(result, item.Worth)
						}
					}
				default:
					result = append(result, item.Worth)
				}
			}
		}
		return result
	case node.TypeObject:
		if v.Node == nil {
			return make(map[string]any)
		}
		return nodeToMap(v.Node)
	default:
		return v.Worth
	}
}
