package main

import (
	"log"

	"github.com/mstgnz/transformer/xml"
)

func runXml() {
	byt, err := xml.ReadXml("example/files/valid.xml")
	if err != nil {
		log.Fatalln(err)
	}

	// xml to node
	knot, _ := xml.DecodeXml(byt)
	knot.Reset()
	knot.Print()

	log.Println()

	// node to xml
	if str, err := xml.NodeToXml(knot); err != nil {
		log.Println("NodeToXml Err: ", err)
	} else {
		log.Println(str)
	}

	//result := xml.ParseXml(byt)
	//log.Println(result)
}
