package main

import (
	"bytes"
	"encoding/xml"
)

var (
	xmlMap interface{}
)

func isXml(doc []byte) bool {
	return xml.Unmarshal(doc, &xmlMap) == nil
}

// xmlReadByte just test code
func xmlReadByte(doc []byte) error {

	buf := bytes.NewBuffer(doc)
	dec := xml.NewDecoder(buf)
	err := dec.Decode(node)
	if err != nil {
		return err
	}

	//lessThan := byte('<') // 60
	//greaterThan := byte('>') // 62
	//slash := byte('/') // 47

	/*for _, v := range doc {
		if v != 9 && v != 10 && v != 32 {

			fmt.Printf("%v %v\n", v, string(v))
		}
	}*/
	//fmt.Println(string([]byte{114, 111, 111, 116}))
	return nil
}
