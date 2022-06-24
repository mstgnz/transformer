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

	//test := convertToNode(knot.GetNode("development")["development"])

	//fmt.Println(test.next)

}
