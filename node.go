package main

import (
	"fmt"
	"reflect"

	"github.com/fatih/color"
)

type ILinear interface {
	AddToStart(key string, value any)
	AddToAfter(key string, value any, which string)
	AddToEnd(key string, value any)
	Delete(key string)
	Exists() bool
	Print()
}

type linear struct {
	key   string
	value any
	next  *linear
}

func Linear() ILinear {
	return &linear{}
}

// AddToStart data
func (node *linear) AddToStart(key string, value any) {
	temp := *node
	node.key = key
	node.value = value
	node.next = &temp
}

// AddToAfter data
func (node *linear) AddToAfter(key string, value any, whichKey string) {
	iter := node
	for iter.key != whichKey && iter.next != nil {
		iter = iter.next
	}
	if iter.key == whichKey {
		temp := *iter
		iter.next = &linear{key: key, value: value, next: nil}
		iter.next.next = temp.next
	} else {
		fmt.Println(whichKey, "not found!")
	}
}

// AddToEnd data
func (node *linear) AddToEnd(key string, value any) {
	iter := node
	if iter.value == nil {
		iter.key = key
		iter.value = value
	} else {
		for iter.next != nil {
			iter = iter.next
		}
		iter.next = &linear{key: key, value: value, next: nil}
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
	for iter.next != nil {
		value := reflect.ValueOf(iter.value).String()
		if reflect.TypeOf(iter.value) == nil {
			value = reflect.Invalid.String()
		}
		fmt.Printf("Key: %v, Value: %v\n", color.YellowString(iter.key), color.YellowString(value))
		iter = iter.next
	}
}

func (node *linear) Exists() bool {
	return node != nil && node.value != nil
}
