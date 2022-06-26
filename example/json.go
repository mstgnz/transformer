package main

import (
	"log"

	"gitgub.com/mstgnz/transformer"
)

func main() {

	byt, err := transformer.JsonRead("example/files/valid.json")
	if err != nil {
		log.Fatalln(err)
	}
	knot, _ := transformer.JsonDecode(byt)
	knot.Print(nil)

}
