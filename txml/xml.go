package txml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mstgnz/transformer/node"
)

// IsXml checks if the given data is valid XML
func IsXml(data []byte) bool {
	return xml.Unmarshal(data, new(any)) == nil
}

// ReadXml reads XML data from a file
func ReadXml(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if !IsXml(data) {
		return nil, fmt.Errorf("invalid XML format")
	}
	return data, nil
}

// DecodeXml decodes XML bytes into a Node
func DecodeXml(data []byte) (*node.Node, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	root := &node.Node{
		Key: "root",
		Value: &node.Value{
			Type: node.TypeObject,
		},
	}

	var current *node.Node = root
	var stack []*node.Node
	var textContent strings.Builder

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch t := token.(type) {
		case xml.StartElement:
			// Create new node
			n := &node.Node{
				Key: t.Name.Local,
				Value: &node.Value{
					Type: node.TypeObject,
				},
			}

			// Handle attributes
			if len(t.Attr) > 0 {
				for _, attr := range t.Attr {
					if current == root && n.Key == "root" {
						// Root attributes are added directly to root node
						attrNode := &node.Node{
							Key: attr.Name.Local,
							Value: &node.Value{
								Type:  node.TypeString,
								Worth: attr.Value,
							},
						}
						if err := root.AddToEnd(attrNode); err != nil {
							return nil, err
						}
					} else {
						// Non-root attributes are added to the current node
						attrNode := &node.Node{
							Key: attr.Name.Local,
							Value: &node.Value{
								Type:  node.TypeString,
								Worth: attr.Value,
							},
						}
						if err := n.AddToEnd(attrNode); err != nil {
							return nil, err
						}
					}
				}
			}

			if err := current.AddToEnd(n); err != nil {
				return nil, err
			}

			stack = append(stack, current)
			current = n
			textContent.Reset()

		case xml.EndElement:
			text := strings.TrimSpace(textContent.String())
			if text != "" && current != nil && current != root {
				current.Value = &node.Value{
					Type:  node.TypeString,
					Worth: text,
				}
			}

			if len(stack) > 0 {
				current = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}
			textContent.Reset()

		case xml.CharData:
			text := string(t)
			if strings.TrimSpace(text) != "" {
				textContent.WriteString(text)
			}
		}
	}

	return root, nil
}

// NodeToXml converts a Node to XML bytes
func NodeToXml(n *node.Node) ([]byte, error) {
	if n == nil {
		return nil, fmt.Errorf("node is nil")
	}

	if n.Value == nil {
		return nil, fmt.Errorf("node value is nil")
	}

	var buf bytes.Buffer
	buf.WriteString(xml.Header)

	if err := writeNode(&buf, n); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// writeNode writes a Node to the buffer
func writeNode(buf *bytes.Buffer, n *node.Node) error {
	if n == nil || n.Value == nil {
		return nil
	}

	// Handle root node
	if n.Key == "root" {
		buf.WriteString("<root")

		// Write root attributes
		if n.Value.Type == node.TypeObject && n.Value.Node != nil {
			current := n.Value.Node
			for current != nil {
				if current.Value != nil && current.Value.Type == node.TypeString {
					buf.WriteByte(' ')
					buf.WriteString(current.Key)
					buf.WriteString("=\"")
					buf.WriteString(escapeXml(current.Value.Worth))
					buf.WriteByte('"')
				}
				current = current.Next
			}
		}

		buf.WriteByte('>')

		// Write child nodes
		if n.Value.Type == node.TypeObject && n.Value.Node != nil {
			current := n.Value.Node
			for current != nil {
				if current.Value != nil && current.Value.Type != node.TypeString {
					if err := writeChildNode(buf, current); err != nil {
						return err
					}
				}
				current = current.Next
			}
		}

		buf.WriteString("</root>")
		return nil
	}

	return writeChildNode(buf, n)
}

// writeChildNode writes a non-root node to the buffer
func writeChildNode(buf *bytes.Buffer, n *node.Node) error {
	if n == nil || n.Value == nil {
		return nil
	}

	// Start tag
	buf.WriteByte('<')
	buf.WriteString(n.Key)

	// Write attributes
	if n.Value.Type == node.TypeObject && n.Value.Node != nil {
		current := n.Value.Node
		for current != nil {
			if current.Value != nil && current.Value.Type == node.TypeString {
				buf.WriteByte(' ')
				buf.WriteString(current.Key)
				buf.WriteString("=\"")
				buf.WriteString(escapeXml(current.Value.Worth))
				buf.WriteByte('"')
			}
			current = current.Next
		}
	}

	// Handle different value types
	switch n.Value.Type {
	case node.TypeString:
		buf.WriteByte('>')
		buf.WriteString(escapeXml(n.Value.Worth))
		buf.WriteString("</")
		buf.WriteString(n.Key)
		buf.WriteByte('>')
		return nil

	case node.TypeNumber, node.TypeBoolean:
		buf.WriteByte('>')
		buf.WriteString(n.Value.Worth)
		buf.WriteString("</")
		buf.WriteString(n.Key)
		buf.WriteByte('>')
		return nil

	case node.TypeObject:
		buf.WriteByte('>')

		// Write child nodes
		if n.Value.Node != nil {
			current := n.Value.Node
			for current != nil {
				if current.Value != nil && current.Value.Type != node.TypeString {
					if err := writeChildNode(buf, current); err != nil {
						return err
					}
				}
				current = current.Next
			}
		}

		buf.WriteString("</")
		buf.WriteString(n.Key)
		buf.WriteByte('>')
		return nil
	}

	return nil
}

// escapeXml escapes special XML characters
func escapeXml(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
