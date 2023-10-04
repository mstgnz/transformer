package main

import (
	"log"

	"gitgub.com/mstgnz/transformer/xml"
)

func runXml() {

	byt, err := xml.ReadXml("example/files/valid.xml")
	if err != nil {
		log.Fatalln(err)
	}
	_, _ = xml.DecodeXml(byt)

}
