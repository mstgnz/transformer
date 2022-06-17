package main

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var (
	valFormat reflect.Type
	typeArr   []any
	typeMap   map[string]any
)

func isJSON(doc []byte) bool {
	return json.Unmarshal(doc, &typeMap) == nil
}

func jsonDecode(doc []byte) error {
	dec := json.NewDecoder(strings.NewReader(string(doc)))
	var (
		key     string
		types   string
		typeVal string
	)
	for {
		t, err := dec.Token()
		if err == io.EOF && err != nil {
			return errors.Wrap(err, "no more")
		}
		types = reflect.TypeOf(t).String()
		typeVal = fmt.Sprintf("%v", t)
		//fmt.Printf("%T -> %v\n", t, t)
		// set key
		if types == "string" && key == "" {
			key = typeVal
			continue
		}
		// set val
		if types == "string" && key != "" {
			node.AddToEnd(key, reflect.ValueOf(t).String())
			key = ""
			continue
		}
		// json.Delim
		if types == "json.Delim" && node.Exists() {
			// set open object - {
			if typeVal == "{" {
				node.AddToEnd(key, typeVal)
				key = ""
				continue
			}
			// set close object - }
			if typeVal == "}" {
				node.AddToEnd(key, typeVal)
				key = ""
				continue
			}
			// set open array - [
			if typeVal == "[" && key != "" {
				node.AddToEnd(key, typeVal)
				key = ""
				continue
			}
			// set close array - ]
			if typeVal == "]" {
				node.AddToEnd(key, typeVal)
				key = ""
				continue
			}
		}
	}
}

// json Unmarshal for map but map is unordered
func jsonRecursive(typeMap map[string]any) {
	valStr := ""
	for key, val := range typeMap {
		valFormat = reflect.TypeOf(val)
		valStr = fmt.Sprintf("%v", val)
		if valFormat.Kind().String() != "map" && valFormat.Kind().String() != "slice" {
			node.AddToEnd(key, valStr)
		} else if valFormat.String() == "map[string]interface {}" {
			convert, err := json.Marshal(val)
			errorHandle(err)
			typeMap = make(map[string]any)
			err = json.Unmarshal(convert, &typeMap)
			errorHandle(err)
			node.AddToEnd(key, nil)
			jsonRecursive(typeMap)
		} else if valFormat.String() == "[]interface {}" {
			convert, err := json.Marshal(val)
			errorHandle(err)
			err = json.Unmarshal(convert, &typeArr)
			errorHandle(err)
			node.AddToEnd(key, nil)
			jsonSubRecursive(typeArr)
		}
	}
}

func jsonSubRecursive(typeArr []any) {
	for key, val := range typeArr {
		valFormat = reflect.TypeOf(val)
		if valFormat.String() == "map[string]interface {}" {
			convert, err := json.Marshal(val)
			typeMap = make(map[string]any)
			err = json.Unmarshal(convert, &typeMap)
			errorHandle(err)
			node.AddToEnd(strconv.Itoa(key), nil)
			jsonRecursive(typeMap)
		}
	}
}

// jsonReadByte just test code for byte read
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
				node.AddToEnd(keySlice[i], nil)
			} else if valSlice[i] == "object" {
				node.AddToEnd(keySlice[i], nil)
			} else {
				node.AddToEnd(keySlice[i], valSlice[i])
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
