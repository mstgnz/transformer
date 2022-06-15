package main

import (
	"encoding/json"
)

func isJSON(by []byte) bool {
	return json.Unmarshal(by, &typeMap) == nil
}
