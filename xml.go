package main

import (
	"encoding/xml"
)

func isXml(doc []byte) bool {
	return xml.Unmarshal(doc, &typeMap) == nil
}
