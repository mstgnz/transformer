package main

import (
	"fmt"
	"log"

	json2 "gitgub.com/mstgnz/transformer/json"
)

func json() {

	byt, err := json2.ReadJson("example/files/valid.json")
	if err != nil {
		log.Fatalln(err)
	}
	knot, _ := json2.DecodeJson(byt)

	fmt.Println(knot.GetNode(nil, "test")[0].Value)
}
