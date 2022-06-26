package transformer

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
func (*Node) AddToNext(knot *Node, parent *Node, key string, value any) *Node {
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
func (*Node) AddToNextWithAttr(knot *Node, parent *Node, key string, value any, attr map[string]string) *Node {
	knot = knot.AddToNext(knot, parent, key, value)
	knot.Attr = attr
	return knot
}

// AddToValue
// It is the method used if the value of a node is a node.
// Parameters are the same as AddToNext.
func (*Node) AddToValue(knot *Node, parent *Node, key string, value any) *Node {
	if knot == nil {
		return &Node{Key: key, Value: value, Prev: knot, Parent: parent}
	}
	knot.Value = &Node{Key: key, Value: value, Prev: knot, Parent: parent}
	obj, ok := knot.Value.(*Node)
	if ok {
		return obj
	}
	return knot
}

// AddToValueWithAttr
// It is the method used if the value of a node is a node.
// Same as AddToNextWithAttr
func (*Node) AddToValueWithAttr(knot *Node, parent *Node, key string, value any, attr map[string]string) *Node {
	knot = knot.AddToValue(knot, parent, key, value)
	knot.Attr = attr
	return knot
}

// SetToValue
// It only allows the key and value to be set.
func (*Node) SetToValue(knot *Node, key string, value any) *Node {
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
func (*Node) AddToArr(knot *Node, value any) *Node {
	if knot == nil {
		knot = &Node{Value: []any{value}}
		return knot
	}
	arr, ok := knot.Value.([]any)
	if ok {
		knot.Value = append(arr, value)
	}
	return knot
}

// AddObjToArr
// If the value of the node is a slice, it appends *Node.
func (*Node) AddObjToArr(knot *Node) *Node {
	if knot == nil {
		knot = &Node{}
		return knot
	}
	newObj := &Node{Prev: knot, Parent: knot}
	arr, ok := knot.Value.([]any)
	if ok {
		knot.Value = append(arr, newObj)
	}
	return newObj
}

// GetNode search Node
func (n *Node) GetNode(knot *Node, key string) map[string]any {
	iter := n
	list := make(map[string]any)
	if iter == nil && knot == nil {
		return list
	}
	if knot == nil {
		for iter.Prev != nil {
			iter = iter.Prev
		}
	} else {
		iter = knot
	}
	for iter != nil {
		if iter.Key == key {
			list[key] = iter
		}
		// if Node value is object
		obj, ok := iter.Value.(*Node)
		if ok {
			n.GetNode(obj, key)
		}
		// if Node value is array and if value in object
		obj1, ok1 := iter.Value.([]any)
		if ok1 {
			for _, v := range obj1 {
				obj2, ok2 := v.(*Node)
				if ok2 {
					n.GetNode(obj2, key)
				}
			}
		}
		iter = iter.Next
	}
	return list
}

// Print Node
func (n *Node) Print(knot *Node) {
	iter := n
	if iter == nil && knot == nil {
		return
	}
	if knot == nil {
		for iter.Prev != nil {
			iter = iter.Prev
		}
	} else {
		iter = knot
	}
	for iter != nil {
		if iter.Parent == nil {
			fmt.Println(color.BlueString("main"))
		}
		n.print(iter)
		// if Node value is object
		obj, ok := iter.Value.(*Node)
		if ok {
			name := iter.Key
			if len(name) == 0 {
				name = "array object"
			}
			fmt.Println(color.BlueString(fmt.Sprintf("child for %v", name)))
			n.Print(obj)
		}
		// if Node value is array and if value in object
		arr, ok1 := iter.Value.([]any)
		if ok1 {
			for _, v := range arr {
				obj2, ok2 := v.(*Node)
				if ok2 {
					fmt.Println(color.BlueString(fmt.Sprintf("child for %v %s", iter.Key, "in object")))
					n.Print(obj2)
				}
			}
		}
		iter = iter.Next
	}
}

func (*Node) print(iter *Node) {
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
