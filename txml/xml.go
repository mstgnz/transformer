// Package txml provides functionality for handling XML data using the Node structure.
// It includes functions for reading, validating, encoding, and decoding XML data.
//
// Namespace handling:
//   - Element keys include the prefix when present: "soap:Envelope"
//   - Namespace declarations are stored as "@xmlns:prefix" or "@xmlns" attributes
//   - Attributes with a namespace prefix are stored as "@prefix:localname"
//   - Round-trip conversion preserves namespace prefixes and declarations
//
// XML comments and processing instructions are ignored during decode.
// Namespace URIs declared in parent elements are available in child elements.
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

const maxDepth = 100

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

// DecodeXml decodes XML bytes into a Node.
//
// Namespace prefixes are preserved in element and attribute keys.
// Each element's namespace declarations are stored as attributes on that element.
func DecodeXml(data []byte) (*node.Node, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	root := &node.Node{
		Key:   "root",
		Value: &node.Value{Type: node.TypeObject},
	}

	var current *node.Node = root
	var stack []*node.Node
	var textContent strings.Builder

	// nsStack holds URI→prefix maps per element level.
	// Each frame contains only the declarations introduced at that level.
	var nsStack []map[string]string

	// resolvePrefix returns the prefix for a namespace URI by searching the stack
	// from innermost to outermost frame.
	resolvePrefix := func(uri string) (string, bool) {
		for i := len(nsStack) - 1; i >= 0; i-- {
			if prefix, ok := nsStack[i][uri]; ok {
				return prefix, true
			}
		}
		return "", false
	}

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
			if len(stack) >= maxDepth {
				return nil, fmt.Errorf("max nesting depth %d exceeded", maxDepth)
			}

			// Flush pending text
			if text := strings.TrimSpace(textContent.String()); text != "" && current != nil {
				current.Value.Type = node.TypeString
				current.Value.Worth = text
			}
			textContent.Reset()

			// Collect namespace declarations from this element's attributes
			// and push a new frame onto the namespace stack.
			frame := make(map[string]string)
			for _, attr := range t.Attr {
				switch {
				case attr.Name.Space == "xmlns":
					// xmlns:prefix="uri"
					frame[attr.Value] = attr.Name.Local
				case attr.Name.Space == "" && attr.Name.Local == "xmlns":
					// xmlns="uri"  (default namespace, empty prefix)
					frame[attr.Value] = ""
				case attr.Name.Space == "" && strings.HasPrefix(attr.Name.Local, "xmlns:"):
					// Older parsers: xmlns:prefix as a single local name
					prefix := strings.TrimPrefix(attr.Name.Local, "xmlns:")
					frame[attr.Value] = prefix
				}
			}
			nsStack = append(nsStack, frame)

			// Determine element key, including namespace prefix when present
			key := t.Name.Local
			if t.Name.Space != "" {
				if prefix, ok := resolvePrefix(t.Name.Space); ok && prefix != "" {
					key = prefix + ":" + t.Name.Local
				}
			}

			n := &node.Node{
				Key:   key,
				Value: &node.Value{Type: node.TypeObject},
			}

			// Attach attributes, translating namespace declarations to @xmlns:* keys
			for _, attr := range t.Attr {
				var attrKey string
				switch {
				case attr.Name.Space == "xmlns":
					attrKey = "@xmlns:" + attr.Name.Local
				case attr.Name.Space == "" && attr.Name.Local == "xmlns":
					attrKey = "@xmlns"
				case attr.Name.Space == "" && strings.HasPrefix(attr.Name.Local, "xmlns:"):
					attrKey = "@" + attr.Name.Local
				case attr.Name.Space != "":
					if prefix, ok := resolvePrefix(attr.Name.Space); ok && prefix != "" {
						attrKey = "@" + prefix + ":" + attr.Name.Local
					} else {
						attrKey = "@" + attr.Name.Local
					}
				default:
					attrKey = "@" + attr.Name.Local
				}
				attrNode := &node.Node{
					Key:   attrKey,
					Value: &node.Value{Type: node.TypeString, Worth: attr.Value},
				}
				if err := n.AddToEnd(attrNode); err != nil {
					return nil, err
				}
			}

			// If the top-level element is named "root", promote it to be the
			// document root instead of wrapping it in an internal root node.
			// The len(stack)==0 guard ensures nested elements named "root" are
			// treated as ordinary children.
			if t.Name.Local == "root" && len(stack) == 0 {
				root = n
				current = n
				continue
			}
			if err := current.AddToEnd(n); err != nil {
				return nil, err
			}

			stack = append(stack, current)
			current = n

		case xml.EndElement:
			// Pop the namespace frame introduced by the matching StartElement
			if len(nsStack) > 0 {
				nsStack = nsStack[:len(nsStack)-1]
			}

			// Flush pending text and infer type
			text := strings.TrimSpace(textContent.String())
			if text != "" && current != nil {
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
			if text := string(t); strings.TrimSpace(text) != "" {
				textContent.WriteString(text)
			}
		}
	}

	return root, nil
}

// NodeToXml converts a Node to XML bytes.
//
// Namespace prefixes in element keys (e.g. "soap:Envelope") are written as-is.
// Attributes whose keys start with "@xmlns" are written as namespace declarations.
func NodeToXml(n *node.Node) ([]byte, error) {
	if n == nil {
		return nil, fmt.Errorf("node is nil")
	}

	var buf bytes.Buffer
	buf.WriteString(xml.Header)

	if err := writeNodeToXml(&buf, n, 0); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func writeNodeToXml(buf *bytes.Buffer, n *node.Node, depth int) error {
	if n == nil {
		return nil
	}
	if depth > maxDepth {
		return fmt.Errorf("max nesting depth %d exceeded", maxDepth)
	}

	tagName := n.Key

	// Open tag
	buf.WriteByte('<')
	buf.WriteString(tagName)

	// Write attributes (keys starting with "@")
	if n.Value != nil && n.Value.Node != nil {
		current := n.Value.Node
		for current != nil {
			if strings.HasPrefix(current.Key, "@") {
				attrName := strings.TrimPrefix(current.Key, "@")
				buf.WriteByte(' ')
				buf.WriteString(attrName)
				buf.WriteString(`="`)
				buf.WriteString(escapeXml(current.Value.Worth))
				buf.WriteByte('"')
			}
			current = current.Next
		}
	}

	// Self-close if element has no content
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
				if item == nil {
					continue
				}
				if item.Node != nil {
					if err := writeNodeToXml(buf, item.Node, depth+1); err != nil {
						return err
					}
				} else {
					buf.WriteString("<item>")
					buf.WriteString(escapeXml(item.Worth))
					buf.WriteString("</item>")
				}
			}

		case node.TypeObject:
			if n.Value.Node != nil {
				current := n.Value.Node
				for current != nil {
					if !strings.HasPrefix(current.Key, "@") {
						if err := writeNodeToXml(buf, current, depth+1); err != nil {
							return err
						}
					}
					current = current.Next
				}
			}
		}
	}

	// Close tag
	buf.WriteString("</")
	buf.WriteString(tagName)
	buf.WriteByte('>')
	return nil
}

// onlyHasAttributes returns true when every child node is an attribute (@-prefixed).
func onlyHasAttributes(n *node.Node) bool {
	for current := n; current != nil; current = current.Next {
		if !strings.HasPrefix(current.Key, "@") {
			return false
		}
	}
	return true
}

// hasNonAttributeChildren returns true when the node has at least one non-attribute child.
func hasNonAttributeChildren(n *node.Node) bool {
	if n.Value == nil || n.Value.Node == nil {
		return false
	}
	for current := n.Value.Node; current != nil; current = current.Next {
		if !strings.HasPrefix(current.Key, "@") {
			return true
		}
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
