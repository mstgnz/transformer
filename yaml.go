package transformer

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// IsYaml Checks if the given file is in xml format.
func IsYaml(byt []byte) bool {
	return yaml.Unmarshal(byt, &yaml.Node{}) == nil
}

// YamlRead Reads the given file, returns as byt
func YamlRead(filename string) ([]byte, error) {
	byt, err := ioutil.ReadFile(filename)
	if err != nil {
		return byt, errors.Wrap(err, "cannot read the file")
	}
	if isY := IsYaml(byt); !isY {
		return byt, errors.Wrap(errors.New("this file is not yaml"), "this file is not yaml")
	}
	return byt, nil
}

// YamlDecode Converts a byte array to a key value struct.
func YamlDecode(byt []byte) (*Node, error) {
	var (
		knot *Node
		yam  yaml.Node
		//parent *Node
	)

	// Decode File
	err := yaml.Unmarshal(byt, &yam)
	if err != nil {
		return knot, err
	}

	/*// Encode
	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	err = enc.Encode(yam.Content[0])
	ErrorHandle(err)*/

	return knot, nil
}
