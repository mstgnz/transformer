package node

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// Node
// This Node structure is a double linked list structure.
// Parent Node link has been added because nested objects are in question.
// Since Json, Yaml and Xml objects are nested, the Value type of our Node structure is designed as a Value structure.
// This Node structure accommodates 3 file types, allows you to do the necessary manipulations and outputs the file type you want.
type Node struct {
	Key    string
	Value  Value
	Next   *Node
	Prev   *Node
	Parent *Node
}

// Value struct type of Node Value.
// This Value structure is enriched to accommodate 3 file types.
// For each Node, only one of the values of the Value structure can be set (should be)
type Value struct {
	Node  *Node
	Array []Value
	Worth string
	Attr  map[string]string // for xml
}

// AddToStart
// Adds Node to the root of our node.
// Our root Node becomes the Next of the new root Node.
// If the knot that comes as a parameter is not nil, do the operation. Node usage should be implemented accordingly.
func (n *Node) AddToStart(knot *Node) *Node {
	if knot == nil {
		return nil
	}
	temp := *n
	*n = *knot
	n.Next = &temp
	n.Next.Prev = n
	if n.Next.Next != nil {
		n.Next.Next.Prev = n.Next
	}
	return n
}

// AddToNext
// Adds a new Node next to the current Node and returns the new Node.
// Prev of the newly created Node is the Node comes with the parameter.
// The Parent comes as a parameter is the Parent of the Node to be created.
// The nil control for the knot that comes as a parameter is only for the health of the flow.
func (n *Node) AddToNext(knot *Node, parent *Node, key string) *Node {
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
func (n *Node) AddToValue(knot *Node, value Value) *Node {
	if knot == nil {
		return nil
	}
	knot.Value = value
	if knot.Value.Node != nil {
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

// AddToEnd
// It adds the Node comes with the parameter to the Next of the last Node of our Node.
// nil control for knot coming with the parameter is not required as it will not cause an error.
func (n *Node) AddToEnd(knot *Node) *Node {
	if n == nil {
		return nil
	}
	iter := n
	for iter.Next != nil {
		iter = iter.Next
	}
	knot.Prev = iter
	iter.Next = knot
	return knot
}

// Delete
// If the Node to be deleted is root, root's next is assigned as the new root.
func (n *Node) Delete(knot *Node) *Node {
	if n == knot {
		*n = *n.Next
		return n
	}
	// If the node to be deleted is in between or at the end, we move our iter object to the previous node of the node to be deleted.
	for n.Next != nil && n.Next != knot {
		*n = *n.Next
	}
	if n.Next != nil {
		if n.Next.Next != nil {
			n.Next = n.Next.Next
			n.Next.Prev = n
		} else {
			n.Next = nil
		}
	}
	return n
}

// GetNode
// It searches nested according to the Key comes as a parameter.
func (n *Node) GetNode(key string) []*Node {
	var list []*Node
	var search func(node *Node)
	search = func(node *Node) {
		for node != nil {
			if node.Key == key {
				list = append(list, node)
			}
			// if Node Value.Node exists
			if node.Value.Node != nil {
				search(node.Value.Node)
			}
			// if Node Value.Array exists
			if len(node.Value.Array) > 0 {
				for _, slc := range node.Value.Array {
					// if Array.Value.Node exists
					if slc.Node != nil {
						search(slc.Node)
					}
				}
			}
			node = node.Next
		}
	}
	search(n.Reset())
	return list
}

// Print
// You can get output to see the node tree structure.
func (n *Node) Print() {
	var write func(node *Node)
	level := 0
	write = func(node *Node) {
		for node != nil {
			node.print(node, level)

			// if Node Value.Node exists
			if node.Value.Node != nil {
				if len(node.Key) == 0 {
					node.Key = "array"
				}
				level++
				write(node.Value.Node)
			}

			// if Node Value.Array exists
			if len(node.Value.Array) > 0 {
				for _, slc := range node.Value.Array {
					// if Array.Value.Node exists
					if slc.Node != nil {
						level++
						write(slc.Node)
					}
				}
			}
			level--
			node = node.Next
		}
	}
	write(n.Reset())
}

// print It belongs to Print
func (n *Node) print(knot *Node, level int) {
	indent := strings.Repeat("  ", level)
	fmt.Printf("%s%v %v %v %+v\n",
		indent,
		color.YellowString("Key: "),
		knot.Key,
		color.YellowString("Value:"),
		knot.Value)
}

// Reset
// It goes to the root of our Node.
func (n *Node) Reset() *Node {
	iter := n
	for iter.Parent != nil {
		iter = iter.Parent
	}
	for iter.Prev != nil {
		iter = iter.Prev
	}
	return iter
}

// Exists Node
func (n *Node) Exists() bool {
	return n != nil
}
