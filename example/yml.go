package main

import (
	"fmt"
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
	for knot.Next != nil {
		fmt.Println(knot.Key)
		knot = knot.Next
	}

}
