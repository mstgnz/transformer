package main

import (
	"fmt"

	"github.com/fatih/color"
)

type IJsonLinear interface {
	AddToStr(key, value string)
	AddToObj(key string)
	AddToArr(key string)
	AppendArr(arr any)
	GetNodeObj(key string) *jsonLinear
	Exists() bool
	Print()
}

type jsonLinear struct {
	key  string
	str  string
	arr  []any
	obj  *jsonLinear
	next *jsonLinear
}

func JsonLinear() IJsonLinear {
	return &jsonLinear{}
}

// AddToStr data
func (node *jsonLinear) AddToStr(key, value string) {
	iter := node
	if iter.next == nil && iter.key == "" {
		iter.key = key
		iter.str = value
	} else {
		for iter.next != nil {
			iter = iter.next
		}
		iter.next = &jsonLinear{key: key, str: value}
	}
}

// AppendArr data
func (node *jsonLinear) AppendArr(arr any) {
	iter := node
	for iter.next != nil {
		iter = iter.next
	}
	iter.arr = append(iter.arr, arr)
}

// AddToArr data
func (node *jsonLinear) AddToArr(key string) {
	iter := node
	if iter.next == nil && iter.str == "" {
		iter.key = key
	} else {
		for iter.next != nil {
			iter = iter.next
		}
		iter.next = &jsonLinear{key: key}
	}
}

// AddToObj data
func (node *jsonLinear) AddToObj(key string) {
	iter := node
	if iter.next == nil && iter.str == "" {
		iter.key = key
		iter.obj = &jsonLinear{}
	} else {
		for iter.next != nil {
			iter = iter.next
		}
		iter.next = &jsonLinear{key: key, obj: &jsonLinear{}}
	}
}

// GetNodeObj for
func (node *jsonLinear) GetNodeObj(key string) *jsonLinear {
	iter := node
	for iter.next != nil {
		if iter.key == key {
			return iter.obj
		}
		iter = iter.next
	}
	return iter.obj
}

// Print data
func (node *jsonLinear) Print() {
	iter := node
	for iter != nil {
		fmt.Printf("%v %v, %v %v, %v %v, %v %v, %v %v\n",
			color.YellowString("Key: "),
			iter.key,
			color.YellowString("Value:"),
			iter.str,
			color.YellowString("Arr:"),
			iter.arr,
			color.YellowString("Obj:"),
			iter.obj,
			color.YellowString("Next:"),
			iter.next)
		iter = iter.next
	}
}

// Exists node
func (node *jsonLinear) Exists() bool {
	return node != nil && (node.str != "" || len(node.arr) > 0 || node.obj != nil || node.next != nil)
}
