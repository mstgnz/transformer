package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

func isYaml(doc []byte) bool {
	return yaml.Unmarshal(doc, &typeMap) == nil
}

func parseYaml(doc []byte) error {

	// Decode File
	var test any
	err := yaml.Unmarshal(doc, &test)
	errorHandle(err)

	// Encode
	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	err = enc.Encode(test)
	errorHandle(err)

	return nil
}
