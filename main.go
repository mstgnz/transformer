package main

import (
	"fmt"
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

	// jsonDecode
	knot, err := jsonDecode(doc)

	// yamlDecode
	//knot, err := yamlDecode(doc)

	// xmlDecode
	//knot, err := xmlDecode(doc)
	if err != nil {
		log.Println(err.Error())
	}
	knot.Print(nil)

	fmt.Println(convertToNode(convertToNode(convertToSlice(convertToNode(knot.value).value)[1]).value).next.next)

	//test := convertToNode(knot.GetNode(nil, "development")["development"])

	//fmt.Println(test.next)

}
