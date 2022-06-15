package main

import (
	"encoding/xml"
)

func isXml(by []byte) bool {
	return xml.Unmarshal(by, &typeMap) == nil
}
