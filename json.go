package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func isJSON(doc []byte) bool {
	return json.Unmarshal(doc, &typeMap) == nil
}

// jsonReadByte just test code
func jsonReadByte(doc []byte) error {
	var j json.RawMessage
	err := json.Unmarshal(doc, &j)
	if err != nil {
		return err
	}
	colon := false
	object := false
	array := false
	str := false
	keyDone := false
	intVal := false
	key, val := "", ""
	// We will render the json object as key and value by looping over the byte array with a single loop.
	// Because I couldn't find an example decoded as key value in go programming language.
	for _, v := range j {
		// except for the spaces
		if v != 32 && v != 10 {
			// key string
			if str && !colon {
				if v == 34 {
					str = false
					key += ", "
					keyDone = true
					continue
				}
				key += string(v)
			}
			// val string
			if str && colon {
				if v == 34 {
					colon, str = false, false
					val += ", "
					continue
				}
				val += string(v)
			}
			// int value
			if contains([]byte{46, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57}, v) {
				val += string(v)
				if !str {
					intVal = true
				}
			}
			// quote -> "
			if v == 34 {
				str = true
				continue
			}
			// colon -> :
			if !str && v == 58 {
				colon = true
				str = false
				continue
			}
			// array start -> [
			if v == 91 {
				array = true
				str, colon = false, false
				if keyDone {
					val += "array, "
					keyDone = false
				}
				continue
			}
			// array end -> ]
			if v == 93 && array {
				array = false
				str = false
				continue
			}
			// object start -> {
			if v == 123 && len(val) > 0 {
				object = true
				str, colon = false, false
				continue
			}
			// object end -> }
			if v == 125 && object {
				object = false
				str = false
				if intVal {
					val += ", "
					intVal = false
				}
				if keyDone {
					val += "object, "
					keyDone = false
				}
				continue
			}
		}
	}
	// trim last comma
	key = strings.Trim(key, ", ")
	val = strings.Trim(val, ", ")

	// Convert to slice
	keySlice := strings.Split(key, ", ")
	valSlice := strings.Split(val, ", ")

	// If key and value lengths are equal, return one by one and add nodes
	if len(keySlice) == len(valSlice) {
		for i := 0; i < len(keySlice); i++ {
			if valSlice[i] == "array" {
				node.AddToEnd(keySlice[i], nil, reflect.TypeOf([]any{}))
			} else if valSlice[i] == "object" {
				node.AddToEnd(keySlice[i], nil, reflect.TypeOf(map[string]any{}))
			} else {
				node.AddToEnd(keySlice[i], valSlice[i], reflect.TypeOf("string"))
			}
		}
	} else {
		fmt.Println(key)
		fmt.Println(val)
		fmt.Println(len(keySlice), len(valSlice))
		return errors.New("key and val not equal")
	}
	return nil
}
