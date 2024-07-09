package main

import (
	"log"

	"github.com/mstgnz/transformer/json"
)

func runJson() {

	byt, err := json.ReadJson("example/files/valid.json")
	if err != nil {
		log.Fatalln(err)
	}
	knot, _ := json.DecodeJson(byt)
	knot = knot.Reset()

	knot.Print()

	log.Println()

}
