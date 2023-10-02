package main

import (
	"fmt"
	"log"

	"gitgub.com/mstgnz/transformer/yml"
)

func runYml() {

	byt, err := yml.ReadYml("example/files/valid.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	knot, err := yml.DecodeYml(byt)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(knot.Key)

}
