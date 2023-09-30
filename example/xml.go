package main

import (
	"log"

	xml2 "gitgub.com/mstgnz/transformer/xml"
)

func xml() {

	byt, err := xml2.ReadXml("example/files/valid.xml")
	if err != nil {
		log.Fatalln(err)
	}
	knot, _ := xml2.DecodeXml(byt)
	knot.Print(nil)

}
