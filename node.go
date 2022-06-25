package main

import (
	"fmt"

	"github.com/fatih/color"
)

type node struct {
	key    string
	value  any
	attr   map[string]string
	next   *node
	prev   *node
	parent *node
}

// AddToNext data
func (*node) AddToNext(knot *node, parent *node, key string, value any) *node {
	if knot == nil {
		knot = &node{key: key, value: value, next: nil, prev: nil, parent: nil}
		return knot
	}
	knot.next = &node{key: key, value: value, next: nil, prev: knot, parent: parent}
	return knot.next
}

// AddToNextWithAttr data
func (*node) AddToNextWithAttr(knot *node, parent *node, key string, value any, attr map[string]string) *node {
	if knot == nil {
		knot = &node{key: key, value: value, attr: attr, next: nil, prev: nil, parent: nil}
		return knot
	}
	knot.next = &node{key: key, value: value, attr: attr, next: nil, prev: knot, parent: parent}
	return knot.next
}

// AddToValue data
func (*node) AddToValue(knot *node, parent *node, key string, value any) *node {
	knot.value = &node{key: key, value: value, next: nil, prev: knot, parent: parent}
	obj, ok := knot.value.(*node)
	if ok {
		return obj
	}
	return knot
}

// AddToValueWithAttr data
func (*node) AddToValueWithAttr(knot *node, parent *node, key string, value any, attr map[string]string) *node {
	knot.value = &node{key: key, value: value, attr: attr, next: nil, prev: knot, parent: parent}
	obj, ok := knot.value.(*node)
	if ok {
		return obj
	}
	return knot
}

// SetToValue data
func (*node) SetToValue(knot *node, key string, value any) *node {
	knot.key = key
	knot.value = value
	return knot
}

// AddToArr data
func (*node) AddToArr(knot *node, value any) *node {
	arr, ok := knot.value.([]any)
	if ok {
		knot.value = append(arr, value)
	}
	return knot
}

// AddObjToArr data
func (*node) AddObjToArr(knot *node) *node {
	newObj := &node{prev: knot, parent: knot}
	arr, ok := knot.value.([]any)
	if ok {
		knot.value = append(arr, newObj)
	}
	return newObj
}

// GetNode search data
func (n *node) GetNode(knot *node, key string) map[string]any {
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
		// if node value is object
		obj, ok := iter.value.(*node)
		if ok {
			n.GetNode(obj, key)
		}
		// if node value is array and if value in object
		obj1, ok1 := iter.value.([]any)
		if ok1 {
			for _, v := range obj1 {
				obj2, ok2 := v.(*node)
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
func (n *node) Print(knot *node) {
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
		// if node value is object
		obj, ok := iter.value.(*node)
		if ok {
			name := iter.key
			if len(name) == 0 {
				name = "array object"
			}
			fmt.Println(color.BlueString(fmt.Sprintf("child for %v", name)))
			n.Print(obj)
			// if node value is array and if value in object
			arr, ok1 := obj.value.([]any)
			if ok1 {
				for _, v := range arr {
					obj2, ok2 := v.(*node)
					if ok2 {
						fmt.Println(color.BlueString(fmt.Sprintf("child for %v %s", obj.key, "in object")))
						n.Print(obj2)
					}
				}
			}
		}
		iter = iter.next
	}
}

func (*node) print(iter *node) {
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

// Exists node
func (n *node) Exists() bool {
	return n != nil
}
