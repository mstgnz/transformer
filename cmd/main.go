package main

import (
	"fmt"
	"log"

	"gitgub.com/mstgnz/transformer"
)

func main() {

	// Json
	jsonDecode()

	// Yaml
	//yamlDecode()

	// Xml
	//xmlDecode()

}

func jsonDecode() {
	byt, err := transformer.JsonRead("files/small.json")
	if err != nil {
		log.Fatalln(err)
	}
	knot, _ := transformer.JsonDecode(byt)
	knot.Print(nil)

	fmt.Println(transformer.ConvertToNode(transformer.ConvertToNode(transformer.ConvertToSlice(transformer.ConvertToNode(knot.Value).Value)[1]).Value).Next.Next)
}

func yamlDecode() {
	byt, err := transformer.YamlRead("files/valid.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	knot, err := transformer.YamlDecode(byt)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(knot)
}

func xmlDecode() {
	byt, err := transformer.XmlRead("files/valid.xml")
	if err != nil {
		log.Fatalln(err)
	}
	knot, err := transformer.XmlDecode(byt)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(knot)
}
