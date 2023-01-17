package main

import (
	"fmt"
	"log"

	"gitgub.com/mstgnz/transformer"
)

func yaml() {

	byt, err := transformer.YamlRead("example/files/valid.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	knot, err := transformer.YamlDecode(byt)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(knot)

}
