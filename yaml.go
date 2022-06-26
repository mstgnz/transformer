package transformer

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// IsYaml Checks if the given file is in xml format.
func isYaml(bytes []byte) bool {
	return yaml.Unmarshal(bytes, new(map[string]any)) == nil
}

// YamlRead Reads the given file, returns as bytes
func YamlRead(filename string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return bytes, errors.Wrap(err, "cannot read the file")
	}
	if isY := IsXml(bytes); !isY {
		return bytes, errors.Wrap(errors.New("this file is not json"), "this file is not yaml")
	}
	return bytes, nil
}

// YamlDecode Converts a byte array to a key value struct.
func YamlDecode(bytes []byte) (*Node, error) {
	var (
		knot *Node
		//parent *Node
	)

	// Decode File
	var test yaml.Node
	err := yaml.Unmarshal(bytes, &test)
	ErrorHandle(err)

	// Encode
	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	err = enc.Encode(test.Content[0])
	ErrorHandle(err)

	return knot, nil
}
