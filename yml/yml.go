package yml

import (
	"os"

	"gitgub.com/mstgnz/transformer/node"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// IsYml
// Checks if the given file is in xml format.
func IsYml(byt []byte) bool {
	return yaml.Unmarshal(byt, &yaml.Node{}) == nil
}

// ReadYml
// Reads the given file, returns as byt
func ReadYml(filename string) ([]byte, error) {
	byt, err := os.ReadFile(filename)
	if err != nil {
		return byt, errors.Wrap(err, "cannot read the file")
	}
	if ok := IsYml(byt); !ok {
		return byt, errors.Wrap(errors.New("this file is not yaml"), "this file is not yaml")
	}
	return byt, nil
}

// DecodeYml
// Converts a byte array to a key value struct.
func DecodeYml(byt []byte) (*node.Node, error) {
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

// NodeToYml
// TODO implement
func NodeToYml(node *node.Node) ([]byte, error) {
	return nil, nil
}
