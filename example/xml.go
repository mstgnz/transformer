package main

import (
	"fmt"
	"log"

	"gitgub.com/mstgnz/transformer"
)

func main() {

	byt, err := transformer.XmlRead("example/files/valid.xml")
	if err != nil {
		log.Fatalln(err)
	}
	knot, err := transformer.XmlDecode(byt)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(knot)

}
