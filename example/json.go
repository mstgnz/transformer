package main

import (
	"fmt"
	"log"

	"gitgub.com/mstgnz/transformer"
)

func main() {

	byt, err := transformer.JsonRead("example/files/small.json")
	if err != nil {
		log.Fatalln(err)
	}
	knot, _ := transformer.JsonDecode(byt)
	knot.Print(nil)

	fmt.Println(transformer.ConvertToNode(transformer.ConvertToNode(transformer.ConvertToSlice(transformer.ConvertToNode(knot.Value).Value)[1]).Value).Next.Next)

}
