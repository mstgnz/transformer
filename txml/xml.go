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
			// Flush any pending text content
			text := strings.TrimSpace(textContent.String())
			if text != "" && current != nil {
				current.Value.Type = node.TypeString
				current.Value.Worth = text
			}
			textContent.Reset()

			n := &node.Node{
				Key: t.Name.Local,
				Value: &node.Value{
					Type: node.TypeObject,
				},
			}

			// Handle attributes
			if len(t.Attr) > 0 {
				for _, attr := range t.Attr {
					attrNode := &node.Node{
						Key: "@" + attr.Name.Local,
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

			if t.Name.Local != "root" {
				if err := current.AddToEnd(n); err != nil {
					return nil, err
				}
			} else {
				root = n
				current = n
				continue
			}
			stack = append(stack, current)
			current = n

		case xml.EndElement:
			// Flush any pending text content
			text := strings.TrimSpace(textContent.String())
			if text != "" && current != nil {
				current.Value.Type = node.TypeString
				current.Value.Worth = text
			}
			textContent.Reset()

			if len(stack) > 0 {
				current = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}

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

	var buf bytes.Buffer
	buf.WriteString(xml.Header)

	if err := writeNodeToXml(&buf, n); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func writeNodeToXml(buf *bytes.Buffer, n *node.Node) error {
	if n == nil {
		return nil
	}

	// Start tag
	buf.WriteByte('<')
	buf.WriteString(n.Key)

	// Write attributes
	if n.Value != nil && n.Value.Node != nil {
		current := n.Value.Node
		for current != nil {
			if strings.HasPrefix(current.Key, "@") {
				buf.WriteByte(' ')
				buf.WriteString(strings.TrimPrefix(current.Key, "@"))
				buf.WriteString("=\"")
				buf.WriteString(escapeXml(current.Value.Worth))
				buf.WriteByte('"')
			}
			current = current.Next
		}
	}

	// Check if element is empty (no content and no child elements)
	isEmpty := n.Value == nil ||
		(n.Value.Type == node.TypeObject && n.Value.Worth == "" &&
			(n.Value.Node == nil || onlyHasAttributes(n.Value.Node)))

	if isEmpty {
		buf.WriteString("/>")
		return nil
	}

	buf.WriteByte('>')

	if n.Value.Type == node.TypeString {
		buf.WriteString(escapeXml(n.Value.Worth))
	} else if n.Value.Node != nil {
		current := n.Value.Node
		for current != nil {
			if !strings.HasPrefix(current.Key, "@") {
				if err := writeNodeToXml(buf, current); err != nil {
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

// onlyHasAttributes checks if a node only has attribute nodes
func onlyHasAttributes(n *node.Node) bool {
	current := n
	for current != nil {
		if !strings.HasPrefix(current.Key, "@") {
			return false
		}
		current = current.Next
	}
	return true
}

func escapeXml(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
