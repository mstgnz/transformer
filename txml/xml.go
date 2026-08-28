// Package txml provides functionality for handling XML data using the Node structure.
// It includes functions for reading, validating, encoding, and decoding XML data.
//
// Repeated elements:
//   - Sibling elements sharing a name arrive as one array node, so
//     <item>1</item><item>2</item> keeps both values instead of the second
//     overwriting the first on the way into JSON or YAML
//   - An array node is written back out as one element per value, which is
//     what the document looked like in the first place
//
// Attributes and text:
//   - Attributes are stored as "@name" child nodes
//   - An element that carries both attributes and text keeps the text in a
//     "#text" child, so neither displaces the other
//   - An element that carries only text is a plain string, number or boolean
//
// Namespace handling:
//   - Element keys include the prefix when present: "soap:Envelope"
//   - Namespace declarations are stored as "@xmlns:prefix" or "@xmlns" attributes
//   - Attributes with a namespace prefix are stored as "@prefix:localname"
//   - Round-trip conversion preserves namespace prefixes and declarations
//
// XML comments and processing instructions are ignored during decode.
// Namespace URIs declared in parent elements are available in child elements.
//
// Mixed content (text interleaved with child elements) survives, but the text
// runs are collected into the element's "#text" node and are written back
// before the child elements rather than between them.
package txml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mstgnz/transformer/internal/scalar"
	"github.com/mstgnz/transformer/node"
)

const maxDepth = 100

// textKey names the child node holding an element's text when the element also
// carries attributes or child elements.
const textKey = "#text"

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
// Sibling elements sharing a name are collected into one array node.
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

			// A child element is about to open, so whatever text came before it
			// belongs to the element as a "#text" child rather than as its value.
			if text := strings.TrimSpace(textContent.String()); text != "" && current != nil {
				if err := addTextChild(current, text); err != nil {
					return nil, err
				}
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

			// Flush pending text. An element with nothing else in it becomes a
			// plain scalar; one that already holds attributes or children keeps
			// its text beside them.
			if text := strings.TrimSpace(textContent.String()); text != "" && current != nil {
				if current.Value != nil && current.Value.Node != nil {
					if err := addTextChild(current, text); err != nil {
						return nil, err
					}
				} else {
					current.Value = textValue(text)
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

	groupRepeated(root)
	return root, nil
}

// addTextChild hangs an element's text off the element as a "#text" node.
func addTextChild(n *node.Node, text string) error {
	return n.AddToEnd(&node.Node{Key: textKey, Value: textValue(text)})
}

// textValue types an element's text. Numbers and the two boolean spellings are
// recognised; everything else is a string.
//
// The reading is deliberately narrow. strconv.ParseBool would turn a lone "t"
// into true, strconv.ParseFloat would call "Inf" a number, and the wider
// numeric forms YAML allows (0x1F, 1_000) are not how XML writes a number, so
// only literals every format can write are treated as numeric.
func textValue(text string) *node.Value {
	if scalar.IsJSONNumber(text) {
		return &node.Value{Type: node.TypeNumber, Worth: text}
	}
	switch strings.ToLower(text) {
	case "true":
		return &node.Value{Type: node.TypeBoolean, Worth: "true"}
	case "false":
		return &node.Value{Type: node.TypeBoolean, Worth: "false"}
	}
	return &node.Value{Type: node.TypeString, Worth: text}
}

// groupRepeated walks the tree and merges runs of same-named sibling elements
// into a single array node.
//
// Without it, <port>34</port><port>55</port> stays two nodes carrying one key,
// and every writer that builds a map on the way out keeps only the last of them.
func groupRepeated(n *node.Node) {
	if n == nil || n.Value == nil {
		return
	}

	if n.Value.Type == node.TypeArray {
		for _, item := range n.Value.Array {
			if item != nil && item.Node != nil {
				groupRepeated(item.Node)
			}
		}
		return
	}

	if n.Value.Node == nil {
		return
	}

	// Group the deepest levels first, so each level sees children that are
	// already settled.
	for child := n.Value.Node; child != nil; child = child.Next {
		groupRepeated(child)
	}
	n.Value.Node = groupChildren(n, n.Value.Node)
}

// groupChildren rebuilds a sibling list with repeated element names collapsed
// into array nodes. A repeated name keeps the position of its first occurrence.
func groupChildren(parent *node.Node, first *node.Node) *node.Node {
	var order []string
	occurrences := make(map[string][]*node.Node)

	for child := first; child != nil; child = child.Next {
		if _, seen := occurrences[child.Key]; !seen {
			order = append(order, child.Key)
		}
		occurrences[child.Key] = append(occurrences[child.Key], child)
	}

	var head, last *node.Node
	link := func(n *node.Node) {
		n.Parent = parent
		n.Prev = last
		n.Next = nil
		if last == nil {
			head = n
		} else {
			last.Next = n
		}
		last = n
	}

	for _, key := range order {
		nodes := occurrences[key]

		// Attributes cannot repeat in valid XML, and collapsing them would put
		// a list where the writer expects a value.
		if len(nodes) == 1 || strings.HasPrefix(key, "@") {
			for _, n := range nodes {
				link(n)
			}
			continue
		}

		arrayNode := &node.Node{
			Key:   key,
			Value: &node.Value{Type: node.TypeArray, Array: make([]*node.Value, 0, len(nodes))},
		}
		for _, n := range nodes {
			arrayNode.Value.Array = append(arrayNode.Value.Array, elementValue(arrayNode, n))
		}
		link(arrayNode)
	}

	return head
}

// elementValue turns an element node into an array element. An element with
// attributes or children stays a node so those survive; a plain scalar is
// stored as its value.
func elementValue(owner *node.Node, n *node.Node) *node.Value {
	if n.Value == nil {
		return &node.Value{Type: node.TypeNull}
	}
	if n.Value.Type != node.TypeObject {
		return n.Value
	}
	n.Parent = owner
	n.Next = nil
	n.Prev = nil
	return &node.Value{Type: node.TypeObject, Node: n}
}

// NodeToXml converts a Node to XML bytes.
//
// Namespace prefixes in element keys (e.g. "soap:Envelope") are written as-is.
// Attributes whose keys start with "@xmlns" are written as namespace declarations.
// An array node is written as one element per value, named after the node.
func NodeToXml(n *node.Node) ([]byte, error) {
	if n == nil {
		return nil, fmt.Errorf("node is nil")
	}

	var buf bytes.Buffer
	buf.WriteString(xml.Header)

	// XML has no top-level sequence, so a root array is wrapped in the root
	// element and its values are written as <item> elements.
	if n.Value != nil && n.Value.Type == node.TypeArray {
		buf.WriteByte('<')
		buf.WriteString(n.Key)
		buf.WriteByte('>')
		for _, item := range n.Value.Array {
			if err := writeElement(&buf, "item", unwrapItem(item), 1); err != nil {
				return nil, err
			}
		}
		buf.WriteString("</")
		buf.WriteString(n.Key)
		buf.WriteByte('>')
		return buf.Bytes(), nil
	}

	if err := writeElement(&buf, n.Key, n.Value, 0); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// writeElement writes one <key>...</key> element holding v.
func writeElement(buf *bytes.Buffer, key string, v *node.Value, depth int) error {
	if depth > maxDepth {
		return fmt.Errorf("max nesting depth %d exceeded", maxDepth)
	}

	buf.WriteByte('<')
	buf.WriteString(key)
	writeAttributes(buf, v)

	if isEmptyElement(v) {
		buf.WriteString("/>")
		return nil
	}

	buf.WriteByte('>')
	if err := writeContent(buf, v, depth); err != nil {
		return err
	}
	buf.WriteString("</")
	buf.WriteString(key)
	buf.WriteByte('>')
	return nil
}

// writeAttributes writes the "@"-prefixed children of v as XML attributes.
func writeAttributes(buf *bytes.Buffer, v *node.Value) {
	if v == nil {
		return
	}
	for current := v.Node; current != nil; current = current.Next {
		if !strings.HasPrefix(current.Key, "@") {
			continue
		}
		buf.WriteByte(' ')
		buf.WriteString(strings.TrimPrefix(current.Key, "@"))
		buf.WriteString(`="`)
		if current.Value != nil {
			buf.WriteString(escapeXml(current.Value.Worth))
		}
		buf.WriteByte('"')
	}
}

// writeContent writes everything between an element's tags.
func writeContent(buf *bytes.Buffer, v *node.Value, depth int) error {
	if v == nil {
		return nil
	}

	switch v.Type {
	case node.TypeArray:
		// Only a nested list reaches here; the outer list was already split
		// into one element per value by writeChild.
		for _, item := range v.Array {
			if err := writeElement(buf, "item", unwrapItem(item), depth+1); err != nil {
				return err
			}
		}
		return nil

	case node.TypeString, node.TypeNumber, node.TypeBoolean:
		buf.WriteString(escapeXml(v.Worth))
		return writeChildren(buf, v, depth)

	case node.TypeObject:
		buf.WriteString(escapeXml(v.Worth))
		return writeChildren(buf, v, depth)
	}

	return nil
}

// writeChildren writes an element's child nodes, skipping attributes and
// turning the "#text" node back into text.
func writeChildren(buf *bytes.Buffer, v *node.Value, depth int) error {
	for current := v.Node; current != nil; current = current.Next {
		if strings.HasPrefix(current.Key, "@") {
			continue
		}
		if current.Key == textKey {
			writeText(buf, current.Value)
			continue
		}
		if err := writeChild(buf, current, depth+1); err != nil {
			return err
		}
	}
	return nil
}

// writeText writes a "#text" node. Several text runs of one element are stored
// as a list and are written one after another.
func writeText(buf *bytes.Buffer, v *node.Value) {
	if v == nil {
		return
	}
	if v.Type == node.TypeArray {
		for _, item := range v.Array {
			if item != nil {
				buf.WriteString(escapeXml(item.Worth))
			}
		}
		return
	}
	buf.WriteString(escapeXml(v.Worth))
}

// writeChild writes one child node. A list is written as one element per value,
// each named after the node, which is how the document held it before decoding.
func writeChild(buf *bytes.Buffer, n *node.Node, depth int) error {
	if n.Value != nil && n.Value.Type == node.TypeArray {
		if len(n.Value.Array) == 0 {
			// An empty list has no values to write. The tag is kept so the key
			// does not disappear from the document altogether.
			return writeElement(buf, n.Key, nil, depth)
		}
		for _, item := range n.Value.Array {
			if err := writeElement(buf, n.Key, unwrapItem(item), depth); err != nil {
				return err
			}
		}
		return nil
	}
	return writeElement(buf, n.Key, n.Value, depth)
}

// unwrapItem returns the value an array element holds. Object elements carry a
// wrapper node; everything else is stored directly.
func unwrapItem(item *node.Value) *node.Value {
	if item == nil {
		return nil
	}
	if item.Node != nil && item.Type == node.TypeObject {
		return item.Node.Value
	}
	return item
}

// isEmptyElement reports whether an element has nothing to write between its
// tags, in which case it is written self-closing.
func isEmptyElement(v *node.Value) bool {
	if v == nil || v.Type == node.TypeNull {
		return true
	}
	if v.Type == node.TypeArray {
		return len(v.Array) == 0
	}
	return v.Worth == "" && !hasContentChildren(v)
}

// hasNonAttributeChildren returns true when the node has at least one non-attribute child.
func hasNonAttributeChildren(n *node.Node) bool {
	if n == nil {
		return false
	}
	return hasContentChildren(n.Value)
}

// hasContentChildren returns true when the value has at least one child that is
// not an attribute.
func hasContentChildren(v *node.Value) bool {
	if v == nil {
		return false
	}
	for current := v.Node; current != nil; current = current.Next {
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
