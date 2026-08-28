// Package tyaml provides functionality for handling YAML data using the Node structure.
// It includes functions for reading, validating, encoding, and decoding YAML data.
//
// Both directions go through yaml.Node rather than through map[string]any.
// A Go map has no order, and the output side used to reindent the encoder's
// text by hand to compensate, which produced documents YAML itself could not
// read back. yaml.Node keeps the document's own key order and lets the encoder
// do the indenting.
//
// YAML anchors are resolved on the way in: an alias arrives as the value it
// points at, and a merge key (<<) is expanded in place, with keys written
// explicitly in the mapping winning over merged ones. The anchor syntax itself
// is not preserved.
package tyaml

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/mstgnz/transformer/internal/scalar"
	"github.com/mstgnz/transformer/node"
	"gopkg.in/yaml.v3"
)

const maxDepth = 100

// maxAliasHops bounds alias chasing. yaml.v3 does not produce cyclic aliases,
// but a hand-built node could, and an unbounded loop would hang rather than
// fail.
const maxAliasHops = 100

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

// DecodeYaml decodes YAML bytes into a Node.
//
// Mapping keys keep the order the document wrote them in.
func DecodeYaml(data []byte) (*node.Node, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}

	// An empty document parses without error and without content.
	if doc.Kind == 0 || len(doc.Content) == 0 {
		return &node.Node{Key: "root", Value: &node.Value{Type: node.TypeNull}}, nil
	}

	return yamlToNode("root", doc.Content[0], 0)
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

// yamlToNode converts one YAML node into a Node with the given key.
func yamlToNode(key string, y *yaml.Node, depth int) (*node.Node, error) {
	n := &node.Node{Key: key}
	value, err := yamlToValue(y, n, depth)
	if err != nil {
		return nil, err
	}
	n.Value = value
	return n, nil
}

// yamlToValue converts one YAML node into a Value. owner is the node the value
// will belong to, so mapping children can point back at their parent.
func yamlToValue(y *yaml.Node, owner *node.Node, depth int) (*node.Value, error) {
	if depth > maxDepth {
		return nil, fmt.Errorf("max nesting depth %d exceeded", maxDepth)
	}

	y = resolveAlias(y)
	if y == nil {
		return &node.Value{Type: node.TypeNull}, nil
	}

	switch y.Kind {
	case yaml.MappingNode:
		return mappingToValue(y, owner, depth)
	case yaml.SequenceNode:
		return sequenceToValue(y, depth)
	case yaml.ScalarNode:
		return scalarToValue(y), nil
	case yaml.DocumentNode:
		if len(y.Content) == 0 {
			return &node.Value{Type: node.TypeNull}, nil
		}
		return yamlToValue(y.Content[0], owner, depth)
	}

	return &node.Value{Type: node.TypeNull}, nil
}

// mappingToValue builds an object value, linking its members in document order.
func mappingToValue(y *yaml.Node, owner *node.Node, depth int) (*node.Value, error) {
	pairs, err := mappingPairs(y, depth)
	if err != nil {
		return nil, err
	}

	v := &node.Value{Type: node.TypeObject}
	var last *node.Node
	for _, p := range pairs {
		child, err := yamlToNode(p.key, p.value, depth+1)
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
	return v, nil
}

// sequenceToValue builds an array value.
func sequenceToValue(y *yaml.Node, depth int) (*node.Value, error) {
	v := &node.Value{Type: node.TypeArray, Array: make([]*node.Value, 0, len(y.Content))}

	for i, item := range y.Content {
		resolved := resolveAlias(item)
		if resolved == nil {
			v.Array = append(v.Array, &node.Value{Type: node.TypeNull})
			continue
		}

		// A mapping element is wrapped in a node so it keeps a key, which is
		// what the XML writer names the repeated element after. Sequences and
		// scalars are stored as the value itself.
		if resolved.Kind == yaml.MappingNode {
			child, err := yamlToNode(fmt.Sprintf("item%d", i), resolved, depth+1)
			if err != nil {
				return nil, err
			}
			v.Array = append(v.Array, &node.Value{Type: node.TypeObject, Node: child})
			continue
		}

		value, err := yamlToValue(resolved, nil, depth+1)
		if err != nil {
			return nil, err
		}
		v.Array = append(v.Array, value)
	}

	return v, nil
}

// scalarToValue maps a YAML scalar onto a Value using the node's resolved tag.
func scalarToValue(y *yaml.Node) *node.Value {
	switch y.Tag {
	case "!!null":
		return &node.Value{Type: node.TypeNull}

	case "!!bool":
		if b, ok := scalar.Bool(y.Value); ok {
			return &node.Value{Type: node.TypeBoolean, Worth: strconv.FormatBool(b)}
		}
		return &node.Value{Type: node.TypeString, Worth: y.Value}

	case "!!int", "!!float":
		if normalized, ok := scalar.Number(y.Value); ok {
			return &node.Value{Type: node.TypeNumber, Worth: normalized}
		}
		return &node.Value{Type: node.TypeString, Worth: y.Value}
	}

	return &node.Value{Type: node.TypeString, Worth: y.Value}
}

// pair is one mapping member after merge keys have been expanded.
type pair struct {
	key   string
	value *yaml.Node
}

// mappingPairs flattens a mapping into ordered members, expanding YAML merge
// keys (<<) where they appear. A key written explicitly in the mapping wins
// over the same key arriving through a merge, and among several merged sources
// the first one wins, which is what the merge specification asks for.
func mappingPairs(y *yaml.Node, depth int) ([]pair, error) {
	if depth > maxDepth {
		return nil, fmt.Errorf("max nesting depth %d exceeded", maxDepth)
	}

	explicit := make(map[string]bool)
	for i := 0; i+1 < len(y.Content); i += 2 {
		if y.Content[i].Tag != "!!merge" {
			explicit[y.Content[i].Value] = true
		}
	}

	var pairs []pair
	added := make(map[string]bool)

	for i := 0; i+1 < len(y.Content); i += 2 {
		key, value := y.Content[i], y.Content[i+1]

		if key.Tag == "!!merge" {
			for _, source := range mergeSources(value) {
				merged, err := mappingPairs(source, depth+1)
				if err != nil {
					return nil, err
				}
				for _, m := range merged {
					if explicit[m.key] || added[m.key] {
						continue
					}
					added[m.key] = true
					pairs = append(pairs, m)
				}
			}
			continue
		}

		// A duplicate key is invalid YAML. Keeping the first occurrence beats
		// emitting the same key twice into JSON or XML further down the line.
		if added[key.Value] {
			continue
		}
		added[key.Value] = true
		pairs = append(pairs, pair{key: key.Value, value: value})
	}

	return pairs, nil
}

// mergeSources returns the mappings a merge key points at. The value is either
// one mapping (usually reached through an alias) or a sequence of them.
func mergeSources(v *yaml.Node) []*yaml.Node {
	v = resolveAlias(v)
	if v == nil {
		return nil
	}

	switch v.Kind {
	case yaml.MappingNode:
		return []*yaml.Node{v}
	case yaml.SequenceNode:
		var sources []*yaml.Node
		for _, item := range v.Content {
			if resolved := resolveAlias(item); resolved != nil && resolved.Kind == yaml.MappingNode {
				sources = append(sources, resolved)
			}
		}
		return sources
	}
	return nil
}

// resolveAlias follows an alias to the node it names.
func resolveAlias(y *yaml.Node) *yaml.Node {
	for hops := 0; y != nil && y.Kind == yaml.AliasNode; hops++ {
		if hops >= maxAliasHops {
			return nil
		}
		y = y.Alias
	}
	return y
}

// NodeToYaml converts a Node to YAML bytes.
func NodeToYaml(n *node.Node) ([]byte, error) {
	if n == nil {
		return nil, fmt.Errorf("node is nil")
	}

	doc, err := valueToYaml(n.Value, 0)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(doc); err != nil {
		enc.Close()
		return nil, err
	}
	if err := enc.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// valueToYaml builds the yaml.Node the encoder writes. Tags are set explicitly
// so a number stays a number and a string that looks like one stays quoted.
func valueToYaml(v *node.Value, depth int) (*yaml.Node, error) {
	if depth > maxDepth {
		return nil, fmt.Errorf("max nesting depth %d exceeded", maxDepth)
	}
	if v == nil {
		return nullNode(), nil
	}

	switch v.Type {
	case node.TypeObject:
		y := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		for current := v.Node; current != nil; current = current.Next {
			value, err := valueToYaml(current.Value, depth+1)
			if err != nil {
				return nil, err
			}
			y.Content = append(y.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: current.Key},
				value)
		}
		return y, nil

	case node.TypeArray:
		y := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		for _, item := range v.Array {
			value, err := valueToYaml(itemValue(item), depth+1)
			if err != nil {
				return nil, err
			}
			y.Content = append(y.Content, value)
		}
		return y, nil

	case node.TypeString:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: v.Worth}, nil

	case node.TypeNumber:
		normalized, ok := scalar.Number(v.Worth)
		if !ok {
			return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: v.Worth}, nil
		}
		tag := "!!float"
		if scalar.IsInteger(normalized) {
			tag = "!!int"
		}
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: tag, Value: normalized}, nil

	case node.TypeBoolean:
		value := "false"
		if b, ok := scalar.Bool(v.Worth); ok && b {
			value = "true"
		}
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: value}, nil

	case node.TypeNull:
		return nullNode(), nil
	}

	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: v.Worth}, nil
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

func nullNode() *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null", Value: "null"}
}
