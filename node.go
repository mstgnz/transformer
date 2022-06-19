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
func (n *node) AddToNext(knot *node, parent *node, key string, value any) *node {
	if knot == nil {
		knot = &node{key: key, value: value, next: nil, prev: nil, parent: nil}
		return knot
	}
	knot.next = &node{key: key, value: value, next: nil, prev: knot, parent: parent}
	return knot.next
}

// AddToNextWithAttr data
func (n *node) AddToNextWithAttr(knot *node, parent *node, key string, value any, attr map[string]string) *node {
	if knot == nil {
		knot = &node{key: key, value: value, attr: attr, next: nil, prev: nil, parent: nil}
		return knot
	}
	knot.next = &node{key: key, value: value, attr: attr, next: nil, prev: knot, parent: parent}
	return knot.next
}

// AddToValue data
func (n *node) AddToValue(knot *node, parent *node, key string, value any) *node {
	knot.value = &node{key: key, value: value, next: nil, prev: knot, parent: parent}
	obj, ok := knot.value.(*node)
	if ok {
		return obj
	}
	return knot
}

// SetToValue data
func (n *node) SetToValue(knot *node, key string, value any) *node {
	knot.key = key
	knot.value = value
	return knot
}

// AddToArr data
func (n *node) AddToArr(knot *node, value any) *node {
	arr, ok := knot.value.([]any)
	if ok {
		knot.value = append(arr, value)
	}
	return knot
}

// AddObjToArr data
func (n *node) AddObjToArr(knot *node) *node {
	newObj := &node{prev: knot, parent: knot}
	arr, ok := knot.value.([]any)
	if ok {
		knot.value = append(arr, newObj)
	}
	return newObj
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
		n.print(iter)
		// if node value is object
		obj, ok := iter.value.(*node)
		if ok {
			fmt.Println("child for", iter.key)
			n.Print(obj)
		}
		// TODO if the node value is slice and one of the slice value is an object
		iter = iter.next
	}
}

func (n *node) print(iter *node) {
	fmt.Printf("%v %v, %v %v, %v %v\n",
		color.YellowString("Key: "),
		iter.key,
		color.YellowString("Value:"),
		iter.value,
		color.YellowString("Attr:"),
		iter.attr)
}

// Exists node
func (n *node) Exists() bool {
	return n != nil
}
