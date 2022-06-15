package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
)

var (
	node    ILinear
	typeMap map[string]any
	typeArr []any
)

func main() {

	// Get Flags
	file := flag.String("file", "", "choose file path")
	from := flag.String("from", "", "current file format (yaml-json-xml)")
	to := flag.String("to", "", "convert file format (yaml-json-xml)")
	flag.Parse()

	// Check Flags
	if len(*file) == 0 || len(*from) == 0 || len(*to) == 0 {
		fmt.Printf("Usage:\t--file	: choose file path\n\t--from	: current file format (yaml-json-xml)\n\t--to	: convert file format (yaml-json-xml)\n")
		os.Exit(0)
	}
	// Info
	fmt.Print(color.RedString("This %v convert from %v to %v\n", *file, *from, *to))

	// Read File
	by, err := ioutil.ReadFile(*file)
	errorHandle(err)

	switch *from {
	case "json":
		if !isJSON(by) {
			fmt.Printf("You said you gave json but this file is not json.\nYour load file: %v", file)
			os.Exit(0)
		}
	case "yaml":
		if !isYaml(by) {
			fmt.Printf("You said you gave yaml but this file is not yaml.\nYour load file: %v", file)
			os.Exit(0)
		}
	case "xml":
		if !isXml(by) {
			fmt.Printf("You said you gave xml but this file is not xml.\nYour load file: %v", file)
			os.Exit(0)
		}
	}

}
