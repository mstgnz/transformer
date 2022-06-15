package main

import (
	"fmt"
	"reflect"

	"github.com/fatih/color"
)

type ILinear interface {
	AddToStart(key string, value any, valueKind reflect.Type)
	AddToAfter(key string, value any, valueKind reflect.Type, which string)
	AddToEnd(key string, value any, valueKind reflect.Type)
	Delete(Key string)
	Print()
}

type linear struct {
	key       string
	value     any
	valueKind reflect.Type
	next      *linear
}

func Linear() ILinear {
	return &linear{}
}

// AddToStart data
func (node *linear) AddToStart(key string, value any, valueKind reflect.Type) {
	temp := *node
	node.key = key
	node.value = value
	node.valueKind = valueKind
	node.next = &temp
}

// AddToAfter data
func (node *linear) AddToAfter(key string, value any, valueKind reflect.Type, whichKey string) {
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
func (node *linear) AddToEnd(key string, value any, valueKind reflect.Type) {
	iter := node
	if iter.next == nil && len(iter.key) == 0 {
		iter.key = key
		iter.value = value
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
	for iter != nil && len(iter.key) > 0 {
		value := reflect.ValueOf(iter.value).String()
		if reflect.TypeOf(iter.value) == nil {
			value = reflect.Invalid.String()
		}
		fmt.Printf("Key: %v, Value: %v, ValueKind: %v\n", color.YellowString(iter.key), color.YellowString(value), color.YellowString(iter.valueKind.String()))
		iter = iter.next
	}
}
