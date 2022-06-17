package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func getArgsAndCheckFile() error {
	// Get Flags
	file := flag.String("file", "", "choose file path")
	from := flag.String("from", "", "current file format (yaml-json-xml)")
	to := flag.String("to", "", "convert file format (yaml-json-xml)")
	flag.Parse()

	// Check Flags
	if len(*file) == 0 || len(*from) == 0 || len(*to) == 0 {
		log.Fatalf("\nUsage:\t--file: choose file path\n\t--from: current file format (yaml-json-xml)\n\t--to: convert file format (yaml-json-xml)\n")
	}
	// Info message
	fmt.Printf("%v %v %v %v %v %v\n", color.RedString("This"), color.YellowString(*file), color.RedString("file are starting to convert from"), color.YellowString(*from), color.RedString("to"), color.YellowString(*to))

	// Read File
	bytes, err := ioutil.ReadFile(*file)
	doc = bytes
	if err != nil {
		return errors.Wrap(err, "cannot read the file")
	}

	switch *from {
	case "json":
		if !isJSON(doc) {
			log.Fatalf("\nYou said you gave json but this file is not json.\nYour load file: %v", *file)
		}
	case "yaml":
		if !isYaml(doc) {
			log.Fatalf("\nYou said you gave yaml but this file is not yaml.\nYour load file: %v", *file)
		}
	case "xml":
		if !isXml(doc) {
			log.Fatalf("\nYou said you gave xml but this file is not xml.\nYour load file: %v", *file)
		}
	}
	return nil
}
