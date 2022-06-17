package main

import (
	"log"
)

//Since it is known from which format the file will be converted
//to which format when the program is run, base definitions are made here.
var (
	doc []byte
)

func main() {

	// Get args and check file
	err := getArgsAndCheckFile()
	if err != nil {
		log.Println(err.Error())
	}

	// Fill Node
	knot, err := jsonDecode(doc)
	if err != nil {
		log.Println(err.Error())
	}
	knot.Print(nil)
}
