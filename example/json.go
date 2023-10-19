package main

import (
	"fmt"
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

	fmt.Println(knot.Next.Next.Value)
	//fmt.Println(knot.GetNode(nil, "test")[0].Value)
}
