package main

import (
	"log"

	xml "gitgub.com/mstgnz/transformer/xml"
)

func runXml() {

	byt, err := xml.ReadXml("example/files/valid.xml")
	if err != nil {
		log.Fatalln(err)
	}
	knot, _ := xml.DecodeXml(byt)
	knot.Print(nil)

}
