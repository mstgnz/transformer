// Package node provides a tree-based data structure implementation for handling hierarchical data.
// It supports various data types including objects, arrays, strings, numbers, and booleans.
// The package is designed to be flexible and can be used for parsing and manipulating
// different data formats like JSON, XML, and YAML.
package node

import (
	"fmt"
	"strings"
)

// ValueType represents the type of a value in the node tree.
// It is implemented as an integer constant to efficiently identify
// the type of data stored in a node's value.
type ValueType int

// Constants representing different value types that can be stored in a node.
// These types align with common data format types (JSON, XML, YAML) for easy conversion.
const (
	TypeNull    ValueType = iota // Represents a null/nil value
	TypeObject                   // Represents an object/map type that can contain child nodes
	TypeArray                    // Represents an array/slice type that can contain multiple values
	TypeString                   // Represents a string value
	TypeNumber                   // Represents a numeric value
	TypeBoolean                  // Represents a boolean value
)

// String returns the string representation of a ValueType.
// This method is useful for debugging and logging purposes.
// It converts the internal ValueType constant to a human-readable string.
func (t ValueType) String() string {
	switch t {
	case TypeNull:
		return "null"
	case TypeObject:
		return "object"
	case TypeArray:
		return "array"
	case TypeString:
		return "string"
	case TypeNumber:
		return "number"
	case TypeBoolean:
		return "boolean"
	default:
		return "unknown"
	}
}

// Value represents a value stored in a node of the tree.
// It contains the type of the value and the actual data,
// which can be a primitive value, an object (node), or an array.
type Value struct {
	Type  ValueType // The type of the value (null, object, array, string, number, boolean)
	Worth string    // The actual value as a string (for primitive types)
	Node  *Node     // Reference to a child node (for object types)
	Array []*Value  // Array of values (for array types)
}

// Node represents a single node in the tree structure.
// Each node can have a key, a value, a parent node, and next/previous sibling nodes.
// This structure allows for building complex hierarchical data structures.
type Node struct {
	Key    string // The key/name of the node
	Value  *Value // The value stored in the node
	Parent *Node  // Reference to the parent node
	Next   *Node  // Reference to the next sibling node
	Prev   *Node  // Reference to the previous sibling node
}

// NewNode creates a new Node with the given key.
// The new node is initialized with a null value type.
// This is the primary constructor for creating new nodes in the tree.
func NewNode(key string) *Node {
	return &Node{
		Key: key,
		Value: &Value{
			Type: TypeNull,
		},
	}
}

// AddToStart adds a node to the start of the current node's children.
// If the current node has no children, the added node becomes the first child.
// If there are existing children, the added node becomes the first child,
// and the existing children are shifted to the right.
// Returns an error if either the current node or the node to add is nil.
func (n *Node) AddToStart(node *Node) error {
	if n == nil {
		return fmt.Errorf("node is nil")
	}

	if node == nil {
		return fmt.Errorf("node to add is nil")
	}

	if n.Value == nil {
		n.Value = &Value{
			Type: TypeObject,
		}
	} else if n.Value.Type != TypeObject {
		n.Value.Type = TypeObject
	}

	node.Parent = n
	if n.Value.Node == nil {
		n.Value.Node = node
	} else {
		node.Next = n.Value.Node
		n.Value.Node.Prev = node
		n.Value.Node = node
	}

	return nil
}

// AddToEnd adds a node to the end of the current node's children.
// If the current node has no children, the added node becomes the first child.
// If there are existing children, the added node is appended after the last child.
// Returns an error if either the current node or the node to add is nil.
func (n *Node) AddToEnd(node *Node) error {
	if n == nil {
		return fmt.Errorf("node is nil")
	}

	if node == nil {
		return fmt.Errorf("node to add is nil")
	}

	if n.Value == nil {
		n.Value = &Value{
			Type: TypeObject,
		}
	} else if n.Value.Type != TypeObject {
		n.Value.Type = TypeObject
	}

	node.Parent = n
	if n.Value.Node == nil {
		n.Value.Node = node
	} else {
		current := n.Value.Node
		for current.Next != nil {
			current = current.Next
		}
		current.Next = node
		node.Prev = current
	}

	return nil
}

// AddToValue adds a value to the current node.
// This method is used to set or update the value of a node.
// If the value contains a node, the parent reference of that node is updated.
// Returns an error if either the current node or the value to add is nil.
func (n *Node) AddToValue(value *Value) error {
	if n == nil {
		return fmt.Errorf("node is nil")
	}

	if value == nil {
		return fmt.Errorf("value to add is nil")
	}

	n.Value = value
	if value.Node != nil {
		value.Node.Parent = n
	}

	return nil
}

// Delete removes the current node from its parent.
// This method handles updating the linked list of siblings and parent references.
// Returns an error if the node is nil, is a root node, or if the parent is invalid.
func (n *Node) Delete() error {
	if n == nil {
		return fmt.Errorf("node is nil")
	}

	if n.Parent == nil {
		return fmt.Errorf("cannot delete root node")
	}

	if n.Parent.Value == nil || n.Parent.Value.Type != TypeObject {
		return fmt.Errorf("parent is not an object")
	}

	if n.Prev != nil {
		n.Prev.Next = n.Next
	} else {
		n.Parent.Value.Node = n.Next
	}

	if n.Next != nil {
		n.Next.Prev = n.Prev
	}

	n.Parent = nil
	n.Next = nil
	n.Prev = nil

	return nil
}

// Print displays the node tree structure in a human-readable format.
// This is a wrapper around the private print method that handles indentation.
// The output format is similar to JSON but with custom formatting for different types.
func (n *Node) Print() {
	n.print(0)
}

// print is a helper function for Print that handles indentation.
// It recursively prints the node tree with proper indentation for nested structures.
// The indentation parameter determines the number of spaces to add before each line.
func (n *Node) print(indent int) {
	if n == nil {
		return
	}

	indentStr := strings.Repeat("  ", indent)
	fmt.Printf("%s%s: ", indentStr, n.Key)

	if n.Value == nil {
		fmt.Println("null")
		return
	}

	switch n.Value.Type {
	case TypeObject:
		fmt.Println("{")
		current := n.Value.Node
		for current != nil {
			current.print(indent + 1)
			current = current.Next
		}
		fmt.Printf("%s}\n", indentStr)

	case TypeArray:
		fmt.Println("[")
		for i, item := range n.Value.Array {
			if item == nil {
				fmt.Printf("%s  null", indentStr)
			} else if item.Node != nil {
				item.Node.print(indent + 1)
			} else {
				switch item.Type {
				case TypeString:
					fmt.Printf("%s  \"%s\"", indentStr, item.Worth)
				case TypeNull:
					fmt.Printf("%s  null", indentStr)
				default:
					fmt.Printf("%s  %v", indentStr, item.Worth)
				}
			}
			if i < len(n.Value.Array)-1 {
				fmt.Printf("\n%s  ,\n", indentStr)
			}
		}
		fmt.Printf("\n%s]\n", indentStr)

	case TypeString:
		fmt.Printf("\"%s\"\n", n.Value.Worth)

	case TypeNumber:
		fmt.Println(n.Value.Worth)

	case TypeBoolean:
		fmt.Println(n.Value.Worth)

	case TypeNull:
		fmt.Println("null")

	default:
		fmt.Printf("unknown type: %v\n", n.Value.Type)
	}
}

// Validate checks if the node tree structure is valid.
// It verifies:
// - Parent-child relationships are properly set
// - Next-prev sibling relationships are consistent
// - All nodes have valid values
// - Array values are valid
// Returns an error if any validation fails.
func (n *Node) Validate() error {
	if n == nil {
		return fmt.Errorf("node is nil")
	}

	if n.Value == nil {
		return fmt.Errorf("value is nil for node %s", n.Key)
	}

	// Check parent-child relationships
	if n.Parent != nil {
		if n.Parent.Value == nil || n.Parent.Value.Type != TypeObject {
			return fmt.Errorf("parent of node %s is not an object", n.Key)
		}

		found := false
		current := n.Parent.Value.Node
		for current != nil {
			if current == n {
				found = true
				break
			}
			current = current.Next
		}
		if !found {
			return fmt.Errorf("node %s is not in parent's children", n.Key)
		}
	}

	// Check next-prev relationships
	if n.Next != nil {
		if n.Next.Prev != n {
			return fmt.Errorf("next node's prev pointer does not point back to node %s", n.Key)
		}
	}
	if n.Prev != nil {
		if n.Prev.Next != n {
			return fmt.Errorf("prev node's next pointer does not point back to node %s", n.Key)
		}
	}

	// Validate children
	if n.Value.Type == TypeObject {
		current := n.Value.Node
		for current != nil {
			if err := current.Validate(); err != nil {
				return fmt.Errorf("invalid child node %s: %v", current.Key, err)
			}
			current = current.Next
		}
	}

	// Validate array values
	if n.Value.Type == TypeArray {
		for i, item := range n.Value.Array {
			if item == nil {
				continue
			}
			if item.Node != nil {
				if err := item.Node.Validate(); err != nil {
					return fmt.Errorf("invalid array item %d: %v", i, err)
				}
			}
		}
	}

	return nil
}

// AddToNext adds a node as the next node
func (n *Node) AddToNext(next *Node) error {
	if n == nil {
		return fmt.Errorf("node is nil")
	}
	if next == nil {
		return fmt.Errorf("next node is nil")
	}
	n.Next = next
	next.Prev = n
	return nil
}

// AddToPrev adds a node as the previous node
func (n *Node) AddToPrev(prev *Node) error {
	if n == nil {
		return fmt.Errorf("node is nil")
	}
	if prev == nil {
		return fmt.Errorf("previous node is nil")
	}
	n.Prev = prev
	prev.Next = n
	return nil
}

// Exists reports whether this node or any node under it carries the given key.
//
// Keys inside arrays count. The earlier version of this only walked object
// children, so a key that lived inside a list, which is where most of a
// document's keys live, was reported as absent.
func (n *Node) Exists(key string) bool {
	if n == nil {
		return false
	}
	if n.Key == key {
		return true
	}
	return existsUnder(n.Value, key, 0)
}

// existsUnder walks a value looking for the key, stopping at the first hit.
func existsUnder(v *Value, key string, depth int) bool {
	if v == nil || depth > maxSearchDepth {
		return false
	}

	switch v.Type {
	case TypeObject:
		for child := v.Node; child != nil; child = child.Next {
			if child.Key == key || existsUnder(child.Value, key, depth+1) {
				return true
			}
		}
	case TypeArray:
		for _, item := range v.Array {
			if item == nil {
				continue
			}
			if item.Node != nil {
				// The element's own key belongs to the list above it, so only
				// what is inside it counts.
				if existsUnder(item.Node.Value, key, depth+1) {
					return true
				}
				continue
			}
			if existsUnder(item, key, depth+1) {
				return true
			}
		}
	}
	return false
}

// Clone returns a deep copy of the node and of every sibling after it.
//
// The copy shares nothing with the original, so a value replaced in one is not
// replaced in the other.
func (n *Node) Clone() *Node {
	if n == nil {
		return nil
	}
	clone := &Node{
		Key:   n.Key,
		Value: n.Value.Clone(),
	}
	if clone.Value != nil && clone.Value.Node != nil {
		clone.Value.Node.Parent = clone
	}
	if n.Next != nil {
		clone.Next = n.Next.Clone()
		clone.Next.Prev = clone
	}
	return clone
}

// Clone returns a deep copy of the value.
//
// Lists inside lists are copied too. The earlier version copied only an
// element's type, text and node, so a list of lists came back with its inner
// lists emptied.
func (v *Value) Clone() *Value {
	if v == nil {
		return nil
	}
	clone := &Value{
		Type:  v.Type,
		Worth: v.Worth,
	}
	if v.Node != nil {
		clone.Node = v.Node.Clone()
	}
	if v.Array != nil {
		clone.Array = make([]*Value, len(v.Array))
		for i, item := range v.Array {
			clone.Array[i] = item.Clone()
			if clone.Array[i] != nil && clone.Array[i].Node != nil {
				clone.Array[i].Node.Parent = nil
			}
		}
	}
	return clone
}

// Equal reports whether two nodes carry the same key and the same value.
//
// Siblings after the node are not part of the comparison, but the children
// under it are, all of them: the earlier version compared only an object's
// first child, which made {"a":1,"b":2} equal to {"a":1}.
func (n *Node) Equal(other *Node) bool {
	if n == nil || other == nil {
		return n == other
	}
	if n.Key != other.Key {
		return false
	}
	return n.Value.Equal(other.Value)
}

// Equal reports whether two values hold the same data, all the way down.
func (v *Value) Equal(other *Value) bool {
	if v == nil || other == nil {
		return v == other
	}
	if v.Type != other.Type || v.Worth != other.Worth {
		return false
	}
	if !siblingsEqual(v.Node, other.Node) {
		return false
	}
	if len(v.Array) != len(other.Array) {
		return false
	}
	for i := range v.Array {
		if !v.Array[i].Equal(other.Array[i]) {
			return false
		}
	}
	return true
}

// siblingsEqual compares two chains of nodes, both what they hold and how many
// of them there are.
func siblingsEqual(a, b *Node) bool {
	for a != nil && b != nil {
		if !a.Equal(b) {
			return false
		}
		a, b = a.Next, b.Next
	}
	return a == nil && b == nil
}

// String returns a string representation of the node
func (n *Node) String() string {
	if n == nil {
		return ""
	}
	if n.Value == nil {
		return n.Key
	}
	switch n.Value.Type {
	case TypeString:
		return fmt.Sprintf("%s: %v", n.Key, n.Value.Worth)
	case TypeObject:
		var b strings.Builder
		b.WriteString(n.Key)
		b.WriteString(": {")
		if n.Value.Node != nil {
			b.WriteString(n.Value.Node.String())
		}
		b.WriteByte('}')
		return b.String()
	case TypeArray:
		var b strings.Builder
		b.WriteString(n.Key)
		b.WriteString(": [")
		for i, v := range n.Value.Array {
			if i > 0 {
				b.WriteString(", ")
			}
			if v != nil && v.Node != nil {
				b.WriteString(v.Node.String())
			} else if v != nil {
				b.WriteString(fmt.Sprintf("%v", v.Worth))
			}
		}
		b.WriteByte(']')
		return b.String()
	default:
		return fmt.Sprintf("%s: %v", n.Key, n.Value.Worth)
	}
}

// Type returns the type of the node's value
func (n *Node) Type() ValueType {
	if n == nil || n.Value == nil {
		return TypeNull
	}
	return n.Value.Type
}
