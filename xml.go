package main

import (
	"encoding/xml"
)

func isXml(doc []byte) bool {
	return xml.Unmarshal(doc, new(interface{})) == nil
}

// xmlReadByte just test code
func xmlReadByte(doc []byte) error {

	return nil
}
