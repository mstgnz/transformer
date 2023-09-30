package yml

import (
	"os"

	"gitgub.com/mstgnz/transformer/node"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// IsYaml Checks if the given file is in xml format.
func IsYaml(byt []byte) bool {
	return yaml.Unmarshal(byt, &yaml.Node{}) == nil
}

// ReadYaml Reads the given file, returns as byt
func ReadYaml(filename string) ([]byte, error) {
	byt, err := os.ReadFile(filename)
	if err != nil {
		return byt, errors.Wrap(err, "cannot read the file")
	}
	if ok := IsYaml(byt); !ok {
		return byt, errors.Wrap(errors.New("this file is not yaml"), "this file is not yaml")
	}
	return byt, nil
}

// DecodeYaml TODO
// DecodeYaml Converts a byte array to a key value struct.
func DecodeYaml(byt []byte) (*node.Node, error) {
	var (
		knot *node.Node
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

// NodeToYml TODO
func NodeToYml(node *node.Node) ([]byte, error) {
	return nil, nil
}
