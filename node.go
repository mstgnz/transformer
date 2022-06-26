package transformer

import (
	"fmt"

	"github.com/fatih/color"
)

type Node struct {
	Key    string
	Value  any
	Attr   map[string]string
	Next   *Node
	Prev   *Node
	Parent *Node
}

// AddToNext data
func (*Node) AddToNext(knot *Node, parent *Node, key string, value any) *Node {
	if knot == nil {
		knot = &Node{Key: key, Value: value, Next: nil, Prev: nil, Parent: nil}
		return knot
	}
	knot.Next = &Node{Key: key, Value: value, Next: nil, Prev: knot, Parent: parent}
	return knot.Next
}

// AddToNextWithAttr data
func (*Node) AddToNextWithAttr(knot *Node, parent *Node, key string, value any, attr map[string]string) *Node {
	if knot == nil {
		knot = &Node{Key: key, Value: value, Attr: attr, Next: nil, Prev: nil, Parent: nil}
		return knot
	}
	knot.Next = &Node{Key: key, Value: value, Attr: attr, Next: nil, Prev: knot, Parent: parent}
	return knot.Next
}

// AddToValue data
func (*Node) AddToValue(knot *Node, parent *Node, key string, value any) *Node {
	knot.Value = &Node{Key: key, Value: value, Next: nil, Prev: knot, Parent: parent}
	obj, ok := knot.Value.(*Node)
	if ok {
		return obj
	}
	return knot
}

// AddToValueWithAttr data
func (*Node) AddToValueWithAttr(knot *Node, parent *Node, key string, value any, attr map[string]string) *Node {
	knot.Value = &Node{Key: key, Value: value, Attr: attr, Next: nil, Prev: knot, Parent: parent}
	obj, ok := knot.Value.(*Node)
	if ok {
		return obj
	}
	return knot
}

// SetToValue data
func (*Node) SetToValue(knot *Node, key string, value any) *Node {
	knot.Key = key
	knot.Value = value
	return knot
}

// AddToArr data
func (*Node) AddToArr(knot *Node, value any) *Node {
	arr, ok := knot.Value.([]any)
	if ok {
		knot.Value = append(arr, value)
	}
	return knot
}

// AddObjToArr data
func (*Node) AddObjToArr(knot *Node) *Node {
	newObj := &Node{Prev: knot, Parent: knot}
	arr, ok := knot.Value.([]any)
	if ok {
		knot.Value = append(arr, newObj)
	}
	return newObj
}

// GetNode search data
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

// Print data
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
