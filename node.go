package main

import (
	"fmt"

	"github.com/fatih/color"
)

type ILinear interface {
	AddToStart(key, value, valueKind string)
	AddToAfter(key, value, valueKind, which string)
	AddToEnd(key, value, valueKind string)
	Delete(Key string)
	Print()
}

type linear struct {
	key       string
	value     string
	valueKind string
	next      *linear
}

func Linear() ILinear {
	return &linear{}
}

// AddToStart data
func (node *linear) AddToStart(key, value, valueKind string) {
	temp := *node
	node.key = key
	node.valueKind = value
	node.valueKind = valueKind
	node.next = &temp
}

// AddToAfter data
func (node *linear) AddToAfter(key, value, valueKind, whichKey string) {
	iter := node
	for iter.key != whichKey && iter.next != nil {
		iter = iter.next
	}
	if iter.key == whichKey {
		temp := *iter
		iter.next = &linear{key: key, value: value, valueKind: valueKind, next: nil}
		iter.next.next = temp.next
	} else {
		fmt.Println(whichKey, "not found!")
	}
}

// AddToEnd data
func (node *linear) AddToEnd(key, value, valueKind string) {
	iter := node
	if iter.next == nil && len(iter.key) == 0 && len(iter.valueKind) == 0 {
		iter.key = key
		iter.valueKind = value
		iter.valueKind = valueKind
	} else {
		for iter.next != nil {
			iter = iter.next
		}
		iter.next = &linear{key: key, value: value, valueKind: valueKind, next: nil}
	}
}

// Delete data
func (node *linear) Delete(key string) {
	iter := node
	if iter.key == key {
		if iter.next != nil {
			node.key = iter.next.key
			node.next = iter.next.next
		} else {
			fmt.Println(key, "is set to zero because it is the last element.")
			node.key = ""
		}
	} else {
		for iter.next != nil && iter.next.key != key {
			iter = iter.next
		}
		if iter.next == nil {
			fmt.Println(key, "not found!")
		} else {
			node.next = iter.next.next
		}
	}
}

// Print data
func (node *linear) Print() {
	iter := node
	for iter != nil {
		fmt.Printf("Key: %v, Value: %v, ValueKind: %v\n", color.YellowString(iter.key), color.YellowString(iter.value), color.YellowString(iter.valueKind))
		iter = iter.next
	}
}
