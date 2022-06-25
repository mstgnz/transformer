package main

import (
	"fmt"
	"log"
)

func main() {

	bytes, err := JsonRead("files/small.json")
	if err != nil {
		log.Fatalln(err.Error())
	}
	knot, err := JsonDecode(bytes)
	knot.Print(nil)

	fmt.Println(ConvertToNode(ConvertToNode(ConvertToSlice(ConvertToNode(knot.value).value)[1]).value).next.next)

}
