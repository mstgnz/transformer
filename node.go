package main

import (
	"fmt"

	"github.com/fatih/color"
)

type Node struct {
	key    string
	value  any
	attr   map[string]string
	next   *Node
	prev   *Node
	parent *Node
}

// AddToNext data
func (*Node) AddToNext(knot *Node, parent *Node, key string, value any) *Node {
	if knot == nil {
		knot = &Node{key: key, value: value, next: nil, prev: nil, parent: nil}
		return knot
	}
	knot.next = &Node{key: key, value: value, next: nil, prev: knot, parent: parent}
	return knot.next
}

// AddToNextWithAttr data
func (*Node) AddToNextWithAttr(knot *Node, parent *Node, key string, value any, attr map[string]string) *Node {
	if knot == nil {
		knot = &Node{key: key, value: value, attr: attr, next: nil, prev: nil, parent: nil}
		return knot
	}
	knot.next = &Node{key: key, value: value, attr: attr, next: nil, prev: knot, parent: parent}
	return knot.next
}

// AddToValue data
func (*Node) AddToValue(knot *Node, parent *Node, key string, value any) *Node {
	knot.value = &Node{key: key, value: value, next: nil, prev: knot, parent: parent}
	obj, ok := knot.value.(*Node)
	if ok {
		return obj
	}
	return knot
}

// AddToValueWithAttr data
func (*Node) AddToValueWithAttr(knot *Node, parent *Node, key string, value any, attr map[string]string) *Node {
	knot.value = &Node{key: key, value: value, attr: attr, next: nil, prev: knot, parent: parent}
	obj, ok := knot.value.(*Node)
	if ok {
		return obj
	}
	return knot
}

// SetToValue data
func (*Node) SetToValue(knot *Node, key string, value any) *Node {
	knot.key = key
	knot.value = value
	return knot
}

// AddToArr data
func (*Node) AddToArr(knot *Node, value any) *Node {
	arr, ok := knot.value.([]any)
	if ok {
		knot.value = append(arr, value)
	}
	return knot
}

// AddObjToArr data
func (*Node) AddObjToArr(knot *Node) *Node {
	newObj := &Node{prev: knot, parent: knot}
	arr, ok := knot.value.([]any)
	if ok {
		knot.value = append(arr, newObj)
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
		for iter.prev != nil {
			iter = iter.prev
		}
	} else {
		iter = knot
	}
	for iter != nil {
		if iter.key == key {
			list[key] = iter
		}
		// if Node value is object
		obj, ok := iter.value.(*Node)
		if ok {
			n.GetNode(obj, key)
		}
		// if Node value is array and if value in object
		obj1, ok1 := iter.value.([]any)
		if ok1 {
			for _, v := range obj1 {
				obj2, ok2 := v.(*Node)
				if ok2 {
					n.GetNode(obj2, key)
				}
			}
		}
		iter = iter.next
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
		for iter.prev != nil {
			iter = iter.prev
		}
	} else {
		iter = knot
	}
	for iter != nil {
		if iter.parent == nil {
			fmt.Println(color.BlueString("main"))
		}
		n.print(iter)
		// if Node value is object
		obj, ok := iter.value.(*Node)
		if ok {
			name := iter.key
			if len(name) == 0 {
				name = "array object"
			}
			fmt.Println(color.BlueString(fmt.Sprintf("child for %v", name)))
			n.Print(obj)
		}
		// if Node value is array and if value in object
		arr, ok1 := iter.value.([]any)
		if ok1 {
			for _, v := range arr {
				obj2, ok2 := v.(*Node)
				if ok2 {
					fmt.Println(color.BlueString(fmt.Sprintf("child for %v %s", iter.key, "in object")))
					n.Print(obj2)
				}
			}
		}
		iter = iter.next
	}
}

func (*Node) print(iter *Node) {
	fmt.Printf("%v %v %v %v %v %v %v %v\n",
		color.YellowString("Key: "),
		iter.key,
		color.YellowString("Value:"),
		iter.value,
		color.YellowString("Attr:"),
		iter.attr,
		color.YellowString("Parent:"),
		iter.parent)
}

// Exists Node
func (n *Node) Exists() bool {
	return n != nil
}
