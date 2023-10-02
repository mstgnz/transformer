package main

import (
	"fmt"
	"log"

	"gitgub.com/mstgnz/transformer/json"
)

func runJson() {

	byt, err := json.ReadJson("example/files/valid.json")
	if err != nil {
		log.Fatalln(err)
	}
	knot, _ := json.DecodeJson(byt)

	fmt.Println(knot.Key)
	//fmt.Println(knot.GetNode(nil, "test")[0].Value)
}
