package main

import (
	"fmt"
	"log"

	"gitgub.com/mstgnz/transformer"
)

func json() {

	byt, err := transformer.JsonRead("example/files/valid.json")
	if err != nil {
		log.Fatalln(err)
	}
	knot, _ := transformer.JsonDecode(byt)

	fmt.Println(knot.GetNode(nil, "containerPort")[0].Value)
}
