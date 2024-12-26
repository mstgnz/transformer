package txml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
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

			// Create new node
			key := t.Name.Local
			if key == "name" {
				key = "name"
			}
			n := &node.Node{
				Key: key,
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

			// Add to parent
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
				// Try to convert the text to appropriate type
				if _, err := strconv.ParseFloat(text, 64); err == nil {
					current.Value.Type = node.TypeNumber
					current.Value.Worth = text
				} else if _, err := strconv.ParseBool(text); err == nil {
					current.Value.Type = node.TypeBoolean
					current.Value.Worth = text
				} else {
					current.Value.Type = node.TypeString
					current.Value.Worth = text
				}
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

	// Get tag name
	tagName := n.Key
	if tagName == "n" {
		tagName = "name"
	}

	// Start tag
	buf.WriteByte('<')
	buf.WriteString(tagName)

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

	// Check if element is empty
	isEmpty := n.Value == nil ||
		(n.Value.Type == node.TypeObject && n.Value.Worth == "" &&
			(n.Value.Node == nil || onlyHasAttributes(n.Value.Node)))

	if isEmpty && !hasNonAttributeChildren(n) {
		buf.WriteString("/>")
		return nil
	}

	buf.WriteByte('>')

	if n.Value != nil {
		switch n.Value.Type {
		case node.TypeString, node.TypeNumber, node.TypeBoolean:
			buf.WriteString(escapeXml(n.Value.Worth))
		case node.TypeArray:
			for _, item := range n.Value.Array {
				if item != nil {
					if item.Node != nil {
						if err := writeNodeToXml(buf, item.Node); err != nil {
							return err
						}
					} else {
						// Write array item as element
						buf.WriteString("<item>")
						buf.WriteString(escapeXml(item.Worth))
						buf.WriteString("</item>")
					}
				}
			}
		case node.TypeObject:
			if n.Value.Node != nil {
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
		}
	}

	// End tag
	buf.WriteString("</")
	buf.WriteString(tagName)
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

// hasNonAttributeChildren checks if a node has any non-attribute children
func hasNonAttributeChildren(n *node.Node) bool {
	if n.Value == nil || n.Value.Node == nil {
		return false
	}

	current := n.Value.Node
	for current != nil {
		if !strings.HasPrefix(current.Key, "@") {
			return true
		}
		current = current.Next
	}
	return false
}

func escapeXml(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
