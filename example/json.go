package main

import (
	"fmt"
	"log"

	json "gitgub.com/mstgnz/transformer/json"
)

func runJson() {

	byt, err := json.ReadJson("example/files/valid.json")
	if err != nil {
		log.Fatalln(err)
	}
	knot, _ := json.DecodeJson(byt)

	fmt.Println(knot.GetNode(nil, "test")[0].Value)
}
