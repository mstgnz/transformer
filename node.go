package main

import (
	"fmt"

	"github.com/fatih/color"
)

type node struct {
	key    string
	value  any
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
	if knot == nil {
		for iter.prev != nil {
			iter = iter.prev
		}
	} else {
		iter = knot
	}
	for iter != nil {
		n.print(iter)
		obj, ok := iter.value.(*node)
		if ok {
			fmt.Println("child for", iter.key)
			n.Print(obj)
		}
		iter = iter.next
	}
}

func (n *node) print(iter *node) {
	fmt.Printf("%v %v, %v %v, %v %v\n",
		color.YellowString("Key: "),
		iter.key,
		color.YellowString("Value:"),
		iter.value,
		color.YellowString("Parent:"),
		iter.parent)
}

// Exists node
func (n *node) Exists() bool {
	return n != nil
}
