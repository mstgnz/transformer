package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"

	"github.com/fatih/color"
)

var (
	node      ILinear
	typeMap   map[string]any
	typeArr   []any
	valFormat reflect.Type
)

func main() {

	getArgsAndCheckFile()

	// Node
	node = Linear()

	// Fill Node
	recursive(typeMap)

	node.Print()

}

func recursive(typeMap map[string]any) {
	valStr := ""
	for key, val := range typeMap {
		valFormat = reflect.TypeOf(val)
		valStr = fmt.Sprintf("%v", val)
		if valFormat.Kind().String() != "map" && valFormat.Kind().String() != "slice" {
			node.AddToEnd(key, valStr, valFormat)
		} else if valFormat.String() == "map[string]interface {}" {
			convert, err := json.Marshal(val)
			errorHandle(err)
			typeMap = make(map[string]any)
			err = json.Unmarshal(convert, &typeMap)
			errorHandle(err)
			node.AddToEnd(key, nil, valFormat)
			recursive(typeMap)
		} else if valFormat.String() == "[]interface {}" {
			convert, err := json.Marshal(val)
			errorHandle(err)
			err = json.Unmarshal(convert, &typeArr)
			errorHandle(err)
			node.AddToEnd(key, nil, valFormat)
			recursive1(typeArr)
		}
	}
}

func recursive1(typeArr []any) {
	for key, val := range typeArr {
		valFormat = reflect.TypeOf(val)
		if valFormat.String() == "map[string]interface {}" {
			convert, err := json.Marshal(val)
			typeMap = make(map[string]any)
			err = json.Unmarshal(convert, &typeMap)
			errorHandle(err)
			node.AddToEnd(strconv.Itoa(key), nil, valFormat)
			recursive(typeMap)
		}
	}
}

func getArgsAndCheckFile() {
	// Get Flags
	file := flag.String("file", "", "choose file path")
	from := flag.String("from", "", "current file format (yaml-json-xml)")
	to := flag.String("to", "", "convert file format (yaml-json-xml)")
	flag.Parse()

	// Check Flags
	if len(*file) == 0 || len(*from) == 0 || len(*to) == 0 {
		fmt.Printf("Usage:\t--file: choose file path\n\t--from: current file format (yaml-json-xml)\n\t--to: convert file format (yaml-json-xml)\n")
		os.Exit(0)
	}
	// Info message
	fmt.Printf("%v %v %v %v %v %v\n", color.RedString("This"), color.YellowString(*file), color.RedString("file are starting to convert from"), color.YellowString(*from), color.RedString("to"), color.YellowString(*to))

	// Read File
	by, err := ioutil.ReadFile(*file)
	errorHandle(err)

	switch *from {
	case "json":
		if !isJSON(by) {
			fmt.Printf("You said you gave json but this file is not json.\nYour load file: %v", *file)
			os.Exit(0)
		}
	case "yaml":
		if !isYaml(by) {
			fmt.Printf("You said you gave yaml but this file is not yaml.\nYour load file: %v", *file)
			os.Exit(0)
		}
	case "xml":
		if !isXml(by) {
			fmt.Printf("You said you gave xml but this file is not xml.\nYour load file: %v", *file)
			os.Exit(0)
		}
	}
}
