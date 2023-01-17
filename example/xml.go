package main

import (
	"log"

	"gitgub.com/mstgnz/transformer"
)

func xml() {

	byt, err := transformer.XmlRead("example/files/valid.xml")
	if err != nil {
		log.Fatalln(err)
	}
	knot, _ := transformer.XmlDecode(byt)
	knot.Print(nil)

}
