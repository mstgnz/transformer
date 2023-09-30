package main

import (
	"fmt"
	"log"

	"gitgub.com/mstgnz/transformer/yml"
)

func yaml() {

	byt, err := yml.ReadYaml("example/files/valid.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	knot, err := yml.DecodeYaml(byt)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(knot)

}
