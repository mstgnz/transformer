package main

import (
	"log"

	"gitgub.com/mstgnz/transformer"
)

func main() {

	bytes, err := transformer.JsonRead("files/small.json")
	if err != nil {
		log.Fatalln(err)
	}
	knot, _ := transformer.JsonDecode(bytes)
	knot.Print(nil)

}
