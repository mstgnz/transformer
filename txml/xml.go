// Package txml provides functionality for handling XML data using the Node structure.
// It includes functions for reading, validating, encoding, and decoding XML data.
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

// IsXml checks if the given bytes represent valid XML
func IsXml(data []byte) bool {
	return xml.Unmarshal(data, new(interface{})) == nil
}

// ReadXml reads an XML file and returns its contents as bytes
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

	if n.Key == "root" && n.Value.Type == node.TypeObject {
		buf.WriteString("<root")

		// Write root attributes
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

		buf.WriteByte('>')

		// Write root children
		current = n.Value.Node
		for current != nil {
			if current.Value != nil && current.Value.Type != node.TypeString {
				if err := writeNode(&buf, current, 0); err != nil {
					return nil, err
				}
			}
			current = current.Next
		}

		buf.WriteString("</root>")
		return buf.Bytes(), nil
	}

	return nil, fmt.Errorf("root node must be an object")
}

// writeNode writes a Node to the buffer with proper indentation
func writeNode(buf *bytes.Buffer, n *node.Node, indent int) error {
	if n == nil || n.Value == nil {
		return nil
	}

	writeIndent(buf, indent)
	buf.WriteByte('<')
	buf.WriteString(n.Key)

	switch n.Value.Type {
	case node.TypeObject:
		// Write attributes first
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

		// Write non-attribute children
		hasChildren := false
		current = n.Value.Node
		for current != nil {
			if current.Value != nil && current.Value.Type != node.TypeString {
				hasChildren = true
				break
			}
			current = current.Next
		}

		if !hasChildren {
			buf.WriteString("/>")
			buf.WriteByte('\n')
			return nil
		}

		buf.WriteByte('>')
		buf.WriteByte('\n')

		current = n.Value.Node
		for current != nil {
			if current.Value != nil && current.Value.Type != node.TypeString {
				if err := writeNode(buf, current, indent+2); err != nil {
					return err
				}
			}
			current = current.Next
		}

		writeIndent(buf, indent)
		buf.WriteString("</")
		buf.WriteString(n.Key)
		buf.WriteByte('>')
		buf.WriteByte('\n')

	case node.TypeArray:
		buf.WriteByte('>')
		buf.WriteByte('\n')
		for _, item := range n.Value.Array {
			if item == nil {
				continue
			}
			if item.Node != nil {
				if err := writeNode(buf, item.Node, indent+2); err != nil {
					return err
				}
			} else {
				writeIndent(buf, indent+2)
				buf.WriteString(escapeXml(item.Worth))
				buf.WriteByte('\n')
			}
		}
		writeIndent(buf, indent)
		buf.WriteString("</")
		buf.WriteString(n.Key)
		buf.WriteByte('>')
		buf.WriteByte('\n')

	case node.TypeString, node.TypeNumber, node.TypeBoolean:
		buf.WriteByte('>')
		buf.WriteString(escapeXml(n.Value.Worth))
		buf.WriteString("</")
		buf.WriteString(n.Key)
		buf.WriteByte('>')
		buf.WriteByte('\n')

	case node.TypeNull:
		buf.WriteString("/>")
		buf.WriteByte('\n')

	default:
		return fmt.Errorf("unknown type: %v", n.Value.Type)
	}

	return nil
}

// writeIndent writes indentation spaces
func writeIndent(buf *bytes.Buffer, n int) {
	for i := 0; i < n; i++ {
		buf.WriteByte(' ')
	}
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
	var textContent string

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
			n := &node.Node{
				Key: t.Name.Local,
				Value: &node.Value{
					Type: node.TypeObject,
				},
			}

			// Handle attributes
			if len(t.Attr) > 0 {
				n.Value.Type = node.TypeObject
				for _, attr := range t.Attr {
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

			if err := current.AddToEnd(n); err != nil {
				return nil, err
			}

			stack = append(stack, current)
			current = n

		case xml.EndElement:
			if textContent != "" {
				current.Value = &node.Value{
					Type:  node.TypeString,
					Worth: textContent,
				}
				textContent = ""
			} else if current.Value.Type == node.TypeObject && current.Value.Node == nil {
				current.Value = &node.Value{
					Type: node.TypeNull,
				}
			}

			if len(stack) > 0 {
				current = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}

		case xml.CharData:
			text := strings.TrimSpace(string(t))
			if text != "" {
				textContent = text
			}
		}
	}

	return root, nil
}

// ParseXml parses XML data into a map
func ParseXml(data []byte) (map[string]interface{}, error) {
	var result interface{}
	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.Strict = false
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}

	if m, ok := result.(map[string]interface{}); ok {
		return m, nil
	}

	// If the root is not a map, wrap it in a map with a "root" key
	return map[string]interface{}{"root": result}, nil
}
