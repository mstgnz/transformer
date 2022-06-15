package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

func isYaml(by []byte) bool {
	return yaml.Unmarshal(by, &typeMap) == nil
}

func parseYaml(b []byte) error {

	// Decode File
	var node any
	err := yaml.Unmarshal(b, &node)
	errorHandle(err)

	// Encode
	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	err = enc.Encode(node)
	errorHandle(err)

	return nil
}
