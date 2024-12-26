package node

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// ValueType represents the type of value stored in a node
type ValueType int

const (
	TypeNull ValueType = iota
	TypeObject
	TypeArray
	TypeString
	TypeNumber
	TypeBoolean
)

// String returns the string representation of ValueType
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

// Node
// This Node structure is a double linked list structure.
// Parent Node link has been added because nested objects are in question.
// Since Json, Yaml and Xml objects are nested, the Value type of our Node structure is designed as a Value structure.
// This Node structure accommodates 3 file types, allows you to do the necessary manipulations and outputs the file type you want.
type Node struct {
	Key    string
	Value  *Value
	Next   *Node
	Prev   *Node
	Parent *Node
}

// Value struct type of Node Value.
// This Value structure is enriched to accommodate 3 file types.
// For each Node, only one of the values of the Value structure can be set (should be)
type Value struct {
	Type  ValueType
	Node  *Node
	Array []*Value
	Worth string
	Attr  map[string]string // for xml
}

// NewNode creates a new node with the given key
func NewNode(key string) *Node {
	return &Node{
		Key: key,
		Value: &Value{
			Type: TypeNull,
		},
	}
}

// AddToStart
// Adds Node to the root of our node.
// Our root Node becomes the Next of the new root Node.
// If the knot that comes as a parameter is not nil, do the operation. Node usage should be implemented accordingly.
func (n *Node) AddToStart(knot *Node) error {
	if n == nil || knot == nil {
		return fmt.Errorf("node or knot is nil")
	}
	knot.Next = n
	n.Prev = knot
	if n.Next != nil {
		n.Next.Prev = n
	}
	return nil
}

// AddToNext
// Adds a new Node next to the current Node and returns the new Node.
// Prev of the newly created Node is the Node comes with the parameter.
// The Parent comes as a parameter is the Parent of the Node to be created.
// The nil control for the knot that comes as a parameter is only for the health of the flow.
func (n *Node) AddToNext(knot *Node, parent *Node, key string) *Node {
	if parent != nil {
		parent = parent.Parent
	}
	if knot == nil {
		return &Node{Key: key, Prev: knot, Parent: parent}
	}
	knot.Next = &Node{Key: key, Prev: knot, Parent: parent}
	return knot.Next
}

// AddToValue
// Sets the Value of the current Node to the given Value
// If the Node in the given Value structure is not empty, it returns this Node.
// If the Array in the given Value structure is not empty and the last element of the Array is a Node and not nil, return this Node.
// The nil control for the Node that comes with the parameter is necessary for flow health. Node usage should be implemented accordingly.
func (n *Node) AddToValue(knot *Node, value *Value) *Node {
	if knot == nil {
		return nil
	}
	parent := knot
	knot.Value = value
	if knot.Value.Node != nil {
		knot.Value.Node.Parent = parent
		return knot.Value.Node
	}
	if size := len(knot.Value.Array); size > 0 {
		slcNode := knot.Value.Array[size-1].Node
		if slcNode != nil {
			return slcNode
		}
	}
	return knot
}

// AddToEnd adds a node to the end of the list
func (n *Node) AddToEnd(knot *Node) error {
	if n == nil || knot == nil {
		return fmt.Errorf("node or knot is nil")
	}
	n.Value.Type = TypeObject
	n.Value.Node = knot
	knot.Parent = n
	return nil
}

// Delete removes the node from the list
func (n *Node) Delete() error {
	if n == nil {
		return fmt.Errorf("node is nil")
	}
	if n.Parent == nil {
		return fmt.Errorf("cannot delete root node")
	}
	if n.Prev != nil {
		n.Prev.Next = n.Next
	}
	if n.Next != nil {
		n.Next.Prev = n.Prev
	}
	if n.Parent.Value.Node == n {
		n.Parent.Value.Node = n.Next
	}
	return nil
}

// GetNode returns all nodes with the given key
func (n *Node) GetNode(key string) []*Node {
	var list []*Node
	var search func(node *Node)
	search = func(node *Node) {
		for node != nil {
			if node.Key == key {
				list = append(list, node)
			}
			if node.Value != nil {
				if node.Value.Node != nil {
					search(node.Value.Node)
				}
				if len(node.Value.Array) > 0 {
					for _, v := range node.Value.Array {
						if v.Node != nil {
							search(v.Node)
						}
					}
				}
			}
			node = node.Next
		}
	}
	search(n.Reset())
	return list
}

// GetNodeByPath returns the node at the given path
func (n *Node) GetNodeByPath(path string) *Node {
	if path == "" {
		return nil
	}
	parts := strings.Split(path, ".")
	current := n.Reset()
	for _, part := range parts {
		found := false
		for current != nil {
			if current.Key == part {
				found = true
				if current.Value != nil && current.Value.Node != nil {
					current = current.Value.Node
				}
				break
			}
			current = current.Next
		}
		if !found {
			return nil
		}
	}
	return current
}

// FindNodes returns all nodes that match the predicate
func (n *Node) FindNodes(predicate func(*Node) bool) []*Node {
	var list []*Node
	var search func(node *Node)
	search = func(node *Node) {
		for node != nil {
			if predicate(node) {
				list = append(list, node)
			}
			if node.Value != nil {
				if node.Value.Node != nil {
					search(node.Value.Node)
				}
				if len(node.Value.Array) > 0 {
					for _, v := range node.Value.Array {
						if v.Node != nil {
							search(v.Node)
						}
					}
				}
			}
			node = node.Next
		}
	}
	search(n.Reset())
	return list
}

// Reset returns the root node
func (n *Node) Reset() *Node {
	if n == nil {
		return nil
	}
	current := n
	for current.Parent != nil {
		current = current.Parent
	}
	for current.Prev != nil {
		current = current.Prev
	}
	return current
}

// Exists returns true if the node exists
func (n *Node) Exists() bool {
	return n != nil
}

// Validate checks if the node structure is valid
func (n *Node) Validate() error {
	if n == nil {
		return fmt.Errorf("node is nil")
	}

	// Check parent-child relationships
	if n.Parent != nil {
		if n.Parent.Value == nil || n.Parent.Value.Node != n {
			return fmt.Errorf("invalid parent-child relationship")
		}
	}

	// Check next-prev relationships
	if n.Next != nil && n.Next.Prev != n {
		return fmt.Errorf("invalid next-prev relationship")
	}
	if n.Prev != nil && n.Prev.Next != n {
		return fmt.Errorf("invalid prev-next relationship")
	}

	// Recursively validate child nodes
	if n.Value != nil && n.Value.Node != nil {
		if err := n.Value.Node.Validate(); err != nil {
			return err
		}
	}

	// Validate next nodes
	if n.Next != nil {
		if err := n.Next.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Print prints the node tree
func (n *Node) Print() {
	var print func(node *Node, level int)
	print = func(node *Node, level int) {
		for node != nil {
			indent := strings.Repeat("  ", level)
			fmt.Printf("%s%s%s %s%v\n",
				indent,
				color.YellowString("Key:"),
				node.Key,
				color.YellowString("Value:"),
				node.Value)

			if node.Value != nil {
				if node.Value.Node != nil {
					print(node.Value.Node, level+1)
				}
				if len(node.Value.Array) > 0 {
					for _, v := range node.Value.Array {
						if v.Node != nil {
							print(v.Node, level+1)
						}
					}
				}
			}
			node = node.Next
		}
	}
	print(n.Reset(), 0)
}
