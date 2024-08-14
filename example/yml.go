package main

import (
	"log"

	"github.com/mstgnz/transformer/yml"
)

func runYml() {

	byt, err := yml.ReadYml("example/files/valid.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	knot, _ := yml.DecodeYml(byt)
	knot = knot.Reset()

	knot.Print()

	log.Println(knot.Parent)

	/* for knot.Next != nil {
		log.Println(knot.Key)
		knot = knot.Next
	} */

}
