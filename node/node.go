package node

import (
	"fmt"

	"github.com/fatih/color"
)

// Node We decode json, xml and yaml formats with a single node.
type Node struct {
	Key    string
	Value  any
	Attr   map[string]string
	Next   *Node
	Prev   *Node
	Parent *Node
}

// AddToNext
// If the given node is zero, it overwrites and returns.
// If knot is not nil, it adds a new node to the next one and returns the newly formed node.
// Params
// knot is the working node.
// parent is the parent of the current node.
// key is the key name
// value supports multiple types. The ones we use are slice, node, string, float etc.
func (n *Node) AddToNext(knot *Node, parent *Node, key string, value any) *Node {
	if knot == nil {
		knot = &Node{Key: key, Value: value}
		return knot
	}
	knot.Next = &Node{Key: key, Value: value, Prev: knot, Parent: parent}
	return knot.Next
}

// AddToNextWithAttr
// It is the method used for the attribute feature in xml format, since we are running things through a single node.
// Since the order of the attribute does not matter, it will be kept as a map.
// Parameters are the same as AddToNext.
func (n *Node) AddToNextWithAttr(knot *Node, parent *Node, key string, value any, attr map[string]string) *Node {
	knot = knot.AddToNext(knot, parent, key, value)
	knot.Attr = attr
	return knot
}

// AddToValue
// It is the method used if the value of a node is a node.
// Parameters are the same as AddToNext.
func (n *Node) AddToValue(knot *Node, parent *Node, key string, value any) *Node {
	if knot == nil {
		return &Node{Key: key, Value: value, Prev: knot, Parent: parent}
	}
	knot.Value = &Node{Key: key, Value: value, Prev: knot, Parent: parent}
	if obj, ok := knot.Value.(*Node); ok {
		return obj
	}
	return knot
}

// AddToValueWithAttr
// It is the method used if the value of a node is a node.
// Same as AddToNextWithAttr
func (n *Node) AddToValueWithAttr(knot *Node, parent *Node, key string, value any, attr map[string]string) *Node {
	knot = knot.AddToValue(knot, parent, key, value)
	knot.Attr = attr
	return knot
}

// SetToValue
// It only allows the key and value to be set.
func (n *Node) SetToValue(knot *Node, key string, value any) *Node {
	if knot == nil {
		return &Node{Key: key, Value: value}
	}
	knot.Key = key
	knot.Value = value
	return knot
}

// AddToArr
// If the value of the node is a slice, it appends value.
// Value supports multiple types. The ones we use are slice, node, string, float etc.
func (n *Node) AddToArr(knot *Node, value any) *Node {
	if knot == nil {
		knot = &Node{Value: []any{value}}
		return knot
	}
	if arr, ok := knot.Value.([]any); ok {
		knot.Value = append(arr, value)
	}
	return knot
}

// AddObjToArr
// If the value of the node is a slice, it appends *Node.
func (n *Node) AddObjToArr(knot *Node) *Node {
	if knot == nil {
		knot = &Node{}
		return knot
	}
	newObj := &Node{Prev: knot, Parent: knot}
	if arr, ok := knot.Value.([]any); ok {
		knot.Value = append(arr, newObj)
	}
	return newObj
}

// GetNode search Node
func (n *Node) GetNode(knot *Node, key string) []*Node {
	var list []*Node
	var search func(node *Node)
	search = func(node *Node) {
		if node == nil {
			return
		}
		for node != nil {
			if node.Key == key {
				list = append(list, node)
			}
			// if Node value is object
			if obj, ok := node.Value.(*Node); ok {
				search(obj)
			}
			// if Node value is array and if value in object
			if arr, ok := node.Value.([]any); ok {
				for _, v := range arr {
					if obj, ok := v.(*Node); ok {
						search(obj)
					}
				}
			}
			node = node.Next
		}
	}
	if knot == nil {
		knot = n.Reset()
	}
	search(knot)
	return list
}

// Print Node
func (n *Node) Print(knot *Node) {
	var write func(node *Node)
	write = func(node *Node) {
		for node != nil {
			if node.Parent == nil {
				fmt.Println(color.BlueString("main"))
			}
			write(node)
			// if Node value is object
			if obj, ok := node.Value.(*Node); ok {
				name := node.Key
				if len(name) == 0 {
					name = "array object"
				}
				fmt.Println(color.BlueString(fmt.Sprintf("child for %v", name)))
				write(obj)
			}
			// if Node value is array and if value in object
			if arr, ok := node.Value.([]any); ok {
				for _, v := range arr {
					if obj, ok := v.(*Node); ok {
						fmt.Println(color.BlueString(fmt.Sprintf("child for %v %s", node.Key, "in object")))
						write(obj)
					}
				}
			}
			node = node.Next
		}
	}
	if knot == nil {
		knot = n.Reset()
	}
	write(knot)
}

func (n *Node) print(iter *Node) {
	fmt.Printf("%v %v %v %v %v %v %v %v\n",
		color.YellowString("Key: "),
		iter.Key,
		color.YellowString("Value:"),
		iter.Value,
		color.YellowString("Attr:"),
		iter.Attr,
		color.YellowString("Parent:"),
		iter.Parent)
}

// Exists Node
func (n *Node) Exists() bool {
	return n != nil
}

func (n *Node) Reset() *Node {
	iter := n
	for iter.Prev != nil {
		iter = iter.Prev
	}
	return iter
}
