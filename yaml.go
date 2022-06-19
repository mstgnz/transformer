package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

func isYaml(doc []byte) bool {
	return yaml.Unmarshal(doc, new(map[string]any)) == nil
}

func yamlDecode(doc []byte) (*node, error) {
	var (
		knot *node
		//parent *node
	)

	// Decode File
	var test yaml.Node
	err := yaml.Unmarshal(doc, &test)
	errorHandle(err)

	// Encode
	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	err = enc.Encode(test.Content[0])
	errorHandle(err)

	return knot, nil
}
