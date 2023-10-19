package yml

import (
	"os"

	"github.com/mstgnz/transformer/node"
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
func DecodeYml(data []byte) (*node.Node, error) {
	var (
		Knot   *node.Node
		Parent *node.Node
		dec    yaml.Node
		parser func(node *yaml.Node)
	)

	// Decode File
	err := yaml.Unmarshal(data, &dec)
	if err != nil {
		return Knot, err
	}

	// recursive
	parser = func(yam *yaml.Node) {
		indent := 0
		current := 0
		for k, child := range yam.Content {
			// mod 2 = Key
			if k%2 == 0 {
				indent = child.Column - 1
				Knot = Knot.AddToNext(Knot, Parent, child.Value)
			} else {
				if child.Kind == yaml.MappingNode {
					Knot = Knot.AddToValue(Knot, node.Value{})
				} else {
					Knot.Value.Worth = child.Value
				}
			}
			if current > indent {
				// end objet
				Parent = Knot.Parent
			}
			parser(child)
		}
	}
	parser(dec.Content[0])

	return Knot, nil
}

// NodeToYml
// TODO implement
func NodeToYml(node *node.Node) ([]byte, error) {
	return nil, nil
}
