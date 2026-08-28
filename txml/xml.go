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
// Type markers:
//
// A few values JSON and YAML can write have no shape of their own in XML. An
// empty tag could be a null, an empty string or an empty object, and a single
// <port> element could be one value or a list holding one. Where the shape
// would otherwise be lost, the writer marks the element with a "type"
// attribute in this package's own namespace and the reader reads it back:
//
//	<n t:type="null"/>                              null
//	<s t:type="string"/>                            the empty string
//	<a t:type="array"/>                             the empty list
//	<a t:type="array"><item>80</item></a>           a list holding one value
//	<a t:type="array"><item t:type="array">…</item> a list of lists
//
// The marker is written only where it is needed, and the namespace is declared
// on the root element only when the document uses one. A list of two or more
// plain values still travels as plain repeated elements, and a document that
// carries no marker is written exactly as it was before markers existed.
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

// itemKey names the elements of a marked list. Their names carry no meaning:
// the marker on the parent is what says the children are list values.
const itemKey = "item"

// markerNamespace identifies this package's type markers. It is a namespace so
// that a marker can never be mistaken for one of the document's own attributes.
const markerNamespace = "https://github.com/mstgnz/transformer"

// markerPrefix is the namespace prefix the writer prefers. It gives way to
// markerPrefix + a number when the document already declares this prefix for
// something else.
const markerPrefix = "t"

const (
	markerAttr   = "type"
	markerNull   = "null"
	markerString = "string"
	markerArray  = "array"
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

// DecodeXml decodes XML bytes into a Node.
//
// Namespace prefixes are preserved in element and attribute keys.
// Each element's namespace declarations are stored as attributes on that element.
// Sibling elements sharing a name are collected into one array node, and type
// markers written by NodeToXml are read back into the values they stand for.
func DecodeXml(data []byte) (*node.Node, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	root := &node.Node{
		Key:   "root",
		Value: &node.Value{Type: node.TypeObject},
	}

	var current *node.Node = root
	var stack []*node.Node
	var textContent strings.Builder

	// markers holds the type marker each open element carries, until its end
	// tag is reached and the marker can be applied.
	markers := make(map[*node.Node]string)

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
			// and push a new frame onto the namespace stack. This package's own
			// namespace is left out: it describes the encoding, not the data.
			frame := make(map[string]string)
			for _, attr := range t.Attr {
				if attr.Value == markerNamespace {
					continue
				}
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
			key := decodeName(t.Name.Local)
			if t.Name.Space != "" {
				if prefix, ok := resolvePrefix(t.Name.Space); ok && prefix != "" {
					key = prefix + ":" + t.Name.Local
				}
			}

			n := &node.Node{
				Key:   key,
				Value: &node.Value{Type: node.TypeObject},
			}

			// Attach attributes, translating namespace declarations to @xmlns:*
			// keys. This package's marker and its declaration are consumed
			// rather than stored: they are not part of the document's data.
			var marker string
			for _, attr := range t.Attr {
				if attr.Name.Space == markerNamespace && attr.Name.Local == markerAttr {
					marker = attr.Value
					continue
				}
				if attr.Value == markerNamespace {
					continue
				}

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
				// A namespace declaration is not a data key, so it is left
				// alone; anything else may carry an escaped name.
				if !strings.HasPrefix(attrKey, "@xmlns") {
					attrKey = "@" + decodeName(strings.TrimPrefix(attrKey, "@"))
				}
				attrNode := &node.Node{
					Key:   attrKey,
					Value: &node.Value{Type: node.TypeString, Worth: attr.Value},
				}
				if err := n.AddToEnd(attrNode); err != nil {
					return nil, err
				}
			}
			if marker != "" {
				markers[n] = marker
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

			if current != nil {
				if marker, ok := markers[current]; ok {
					applyMarker(current, marker)
					delete(markers, current)
				}
			}

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

// applyMarker turns an element the writer marked back into the value the marker
// stands for. An unknown marker is ignored, so a document from elsewhere that
// happens to use this namespace is read as ordinary data.
func applyMarker(n *node.Node, marker string) {
	switch marker {
	case markerNull:
		n.Value = &node.Value{Type: node.TypeNull}
	case markerString:
		if n.Value == nil {
			n.Value = &node.Value{}
		}
		n.Value.Type = node.TypeString
	case markerArray:
		n.Value = listFromChildren(n)
	}
}

// listFromChildren turns a marked element's children into list values. Their
// names carry no meaning here; the marker on the parent is what said they are
// list values rather than named fields.
func listFromChildren(n *node.Node) *node.Value {
	list := &node.Value{Type: node.TypeArray, Array: []*node.Value{}}
	if n.Value == nil {
		return list
	}

	// Collect first: turning a child into a list value unlinks it.
	var children []*node.Node
	for child := n.Value.Node; child != nil; child = child.Next {
		if strings.HasPrefix(child.Key, "@") || child.Key == textKey {
			continue
		}
		children = append(children, child)
	}

	holder := &node.Node{Key: n.Key, Value: list}
	for _, child := range children {
		list.Array = append(list.Array, elementValue(holder, child))
	}
	return list
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
// A list of two or more plain values is written as one element per value; the
// shapes XML cannot hold on their own carry a type marker instead.
func NodeToXml(n *node.Node) ([]byte, error) {
	if n == nil {
		return nil, fmt.Errorf("node is nil")
	}

	w := newXmlWriter(n.Value)
	rootName := w.qualify(n.Key)

	var attrs, content bytes.Buffer
	w.writeAttributes(&attrs, n.Value)

	marker := valueMarker(n.Value)
	if marker == "" && n.Value != nil && n.Value.Type == node.TypeArray {
		// XML has no top-level sequence, so a root list always needs the marker
		// even when its values would otherwise travel as repeated elements.
		marker = markerArray
	}
	if marker != "" {
		w.usedMarker = true
	}

	selfClose := true
	switch marker {
	case markerNull, markerString:
	case markerArray:
		if err := w.writeListItems(&content, n.Value, 1); err != nil {
			return nil, err
		}
		selfClose = content.Len() == 0
	default:
		if selfClose = isEmptyElement(n.Value); !selfClose {
			if err := w.writeContent(&content, n.Value, 0); err != nil {
				return nil, err
			}
		}
	}

	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	buf.WriteByte('<')
	buf.WriteString(rootName)
	buf.Write(attrs.Bytes())
	if w.usedMarker {
		buf.WriteString(` xmlns:`)
		buf.WriteString(w.prefix)
		buf.WriteString(`="`)
		buf.WriteString(markerNamespace)
		buf.WriteByte('"')
	}
	if marker != "" {
		w.writeMarker(&buf, marker)
	}
	if selfClose {
		buf.WriteString("/>")
		return buf.Bytes(), nil
	}
	buf.WriteByte('>')
	buf.Write(content.Bytes())
	buf.WriteString("</")
	buf.WriteString(rootName)
	buf.WriteByte('>')
	return buf.Bytes(), nil
}

// xmlWriter carries the state the element writers share: which prefix the type
// markers use, and whether any of them was needed, which is what decides
// whether the root declares the namespace at all.
type xmlWriter struct {
	prefix     string
	usedMarker bool

	// declared holds the namespace prefixes in scope. A key like
	// "soap:Envelope" is written as a prefixed name only when its prefix is
	// one of them; "tel:+90" is a key that merely contains a colon, and
	// writing it as a prefix would produce a document no namespace-aware
	// parser accepts.
	declared map[string]bool
}

// newXmlWriter starts a writer with the prefixes XML defines on its own, plus
// the one this package uses for type markers.
func newXmlWriter(root *node.Value) *xmlWriter {
	w := &xmlWriter{
		prefix:   markerPrefixFor(root),
		declared: map[string]bool{"xml": true, "xmlns": true},
	}
	w.declared[w.prefix] = true
	w.declare(root)
	return w
}

// declare adds the prefixes a value's own xmlns attributes introduce, and
// returns them so the caller can take them back out of scope afterwards.
func (w *xmlWriter) declare(v *node.Value) []string {
	if v == nil {
		return nil
	}
	var added []string
	for current := v.Node; current != nil; current = current.Next {
		if !strings.HasPrefix(current.Key, "@xmlns:") {
			continue
		}
		prefix := strings.TrimPrefix(current.Key, "@xmlns:")
		if w.declared[prefix] {
			continue
		}
		w.declared[prefix] = true
		added = append(added, prefix)
	}
	return added
}

func (w *xmlWriter) undeclare(prefixes []string) {
	for _, prefix := range prefixes {
		delete(w.declared, prefix)
	}
}

// qualify turns a node key into a legal XML name, keeping a namespace prefix
// as a prefix when the document declares it and escaping the colon when it
// does not.
func (w *xmlWriter) qualify(key string) string {
	if i := strings.Index(key, ":"); i > 0 && w.declared[key[:i]] {
		return key[:i] + ":" + encodeName(key[i+1:])
	}
	return encodeName(key)
}

// markerPrefixFor picks a marker prefix the document does not already use for
// something of its own.
func markerPrefixFor(v *node.Value) string {
	declared := make(map[string]bool)
	if v != nil {
		for current := v.Node; current != nil; current = current.Next {
			if strings.HasPrefix(current.Key, "@xmlns:") {
				declared[strings.TrimPrefix(current.Key, "@xmlns:")] = true
			}
		}
	}
	if !declared[markerPrefix] {
		return markerPrefix
	}
	for i := 1; ; i++ {
		candidate := fmt.Sprintf("%s%d", markerPrefix, i)
		if !declared[candidate] {
			return candidate
		}
	}
}

// valueMarker returns the marker a value needs to survive XML, or the empty
// string when the element's own shape already says what it is.
func valueMarker(v *node.Value) string {
	if v == nil {
		return ""
	}
	switch v.Type {
	case node.TypeNull:
		return markerNull
	case node.TypeString:
		if v.Worth == "" && !hasContentChildren(v) {
			return markerString
		}
	case node.TypeArray:
		return markerArray
	}
	return ""
}

// listNeedsMarker reports whether a list has to be written as a marked element
// rather than as repeated elements. Two or more plain values are their own
// evidence; anything else is ambiguous once written.
func listNeedsMarker(v *node.Value) bool {
	if len(v.Array) < 2 {
		return true
	}
	for _, item := range v.Array {
		if inner := unwrapItem(item); inner != nil && inner.Type == node.TypeArray {
			return true
		}
	}
	return false
}

func (w *xmlWriter) writeMarker(buf *bytes.Buffer, marker string) {
	buf.WriteByte(' ')
	buf.WriteString(w.prefix)
	buf.WriteByte(':')
	buf.WriteString(markerAttr)
	buf.WriteString(`="`)
	buf.WriteString(marker)
	buf.WriteByte('"')
}

// writeElement writes one <key>...</key> element holding v.
func (w *xmlWriter) writeElement(buf *bytes.Buffer, key string, v *node.Value, depth int) error {
	if depth > maxDepth {
		return fmt.Errorf("max nesting depth %d exceeded", maxDepth)
	}

	marker := valueMarker(v)
	if marker != "" {
		w.usedMarker = true
	}

	// The element's own declarations come into scope first: an element may
	// use a prefix it declares on itself.
	added := w.declare(v)
	defer w.undeclare(added)
	name := w.qualify(key)

	buf.WriteByte('<')
	buf.WriteString(name)
	w.writeAttributes(buf, v)
	if marker != "" {
		w.writeMarker(buf, marker)
	}

	if isEmptyElement(v) {
		buf.WriteString("/>")
		return nil
	}

	buf.WriteByte('>')
	if err := w.writeContent(buf, v, depth); err != nil {
		return err
	}
	buf.WriteString("</")
	buf.WriteString(name)
	buf.WriteByte('>')
	return nil
}

// writeAttributes writes the "@"-prefixed children of v as XML attributes.
func (w *xmlWriter) writeAttributes(buf *bytes.Buffer, v *node.Value) {
	if v == nil {
		return
	}
	for current := v.Node; current != nil; current = current.Next {
		if !strings.HasPrefix(current.Key, "@") {
			continue
		}
		buf.WriteByte(' ')
		buf.WriteString(w.qualify(strings.TrimPrefix(current.Key, "@")))
		buf.WriteString(`="`)
		if current.Value != nil {
			buf.WriteString(escapeXml(current.Value.Worth))
		}
		buf.WriteByte('"')
	}
}

// writeContent writes everything between an element's tags.
func (w *xmlWriter) writeContent(buf *bytes.Buffer, v *node.Value, depth int) error {
	if v == nil {
		return nil
	}

	switch v.Type {
	case node.TypeArray:
		return w.writeListItems(buf, v, depth+1)

	case node.TypeString, node.TypeNumber, node.TypeBoolean:
		buf.WriteString(escapeXml(v.Worth))
		return w.writeChildren(buf, v, depth)

	case node.TypeObject:
		buf.WriteString(escapeXml(v.Worth))
		return w.writeChildren(buf, v, depth)
	}

	return nil
}

// writeListItems writes the values of a marked list as <item> elements. Their
// names carry no meaning; the marker on the parent is what says they are list
// values rather than named fields.
func (w *xmlWriter) writeListItems(buf *bytes.Buffer, v *node.Value, depth int) error {
	for _, item := range v.Array {
		if err := w.writeElement(buf, itemKey, unwrapItem(item), depth); err != nil {
			return err
		}
	}
	return nil
}

// writeChildren writes an element's child nodes, skipping attributes and
// turning the "#text" node back into text.
func (w *xmlWriter) writeChildren(buf *bytes.Buffer, v *node.Value, depth int) error {
	for current := v.Node; current != nil; current = current.Next {
		if strings.HasPrefix(current.Key, "@") {
			continue
		}
		if current.Key == textKey {
			writeText(buf, current.Value)
			continue
		}
		if err := w.writeChild(buf, current, depth+1); err != nil {
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

// writeChild writes one child node. A list of two or more plain values is
// written as one element per value, which is how the document held it before
// decoding; anything else carries a marker so its shape survives.
func (w *xmlWriter) writeChild(buf *bytes.Buffer, n *node.Node, depth int) error {
	if n.Value != nil && n.Value.Type == node.TypeArray && !listNeedsMarker(n.Value) {
		for _, item := range n.Value.Array {
			if err := w.writeElement(buf, n.Key, unwrapItem(item), depth); err != nil {
				return err
			}
		}
		return nil
	}
	return w.writeElement(buf, n.Key, n.Value, depth)
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

// XML names are not free text: an element cannot be called "first name" or
// "1st", so a JSON key that is not a legal name would produce a document no
// parser accepts. Such a name is escaped on the way out and read back on the
// way in, the way XmlConvert.EncodeName has long done it, so the key survives
// and the document stays valid.
//
// A name that is already legal is written untouched, which is nearly all of
// them.

// emptyNameEscape stands for the empty key. No character escapes to it, since
// a NUL cannot appear in XML text at all.
const emptyNameEscape = "_x0000_"

// encodeName returns a legal XML name for a node key.
func encodeName(name string) string {
	if name == "" {
		return emptyNameEscape
	}
	if isLegalName(name) && !strings.Contains(name, "_x") {
		return name
	}

	var b strings.Builder
	for i, r := range name {
		switch {
		case r == '_' && looksEscaped(name[i:]):
			// Escape the underscore so the sequence cannot be read back as an
			// escape that was never written as one.
			b.WriteString("_x005F_")
		case i == 0 && isNameStart(r), i > 0 && isNameChar(r):
			b.WriteRune(r)
		default:
			fmt.Fprintf(&b, "_x%04X_", r)
		}
	}
	return b.String()
}

// decodeName reverses encodeName.
func decodeName(name string) string {
	if name == emptyNameEscape {
		return ""
	}
	if !strings.Contains(name, "_x") {
		return name
	}

	var b strings.Builder
	for i := 0; i < len(name); {
		if name[i] == '_' && looksEscaped(name[i:]) {
			var r rune
			if _, err := fmt.Sscanf(name[i+2:i+6], "%04X", &r); err == nil {
				b.WriteRune(r)
				i += 7
				continue
			}
		}
		b.WriteByte(name[i])
		i++
	}
	return b.String()
}

// looksEscaped reports whether s opens with an _xHHHH_ escape.
func looksEscaped(s string) bool {
	if len(s) < 7 || s[0] != '_' || s[1] != 'x' || s[6] != '_' {
		return false
	}
	for i := 2; i < 6; i++ {
		c := s[i]
		if (c < '0' || c > '9') && (c < 'A' || c > 'F') && (c < 'a' || c > 'f') {
			return false
		}
	}
	return true
}

// isLegalName reports whether a name can be written as an XML name as it is.
func isLegalName(name string) bool {
	for i, r := range name {
		if i == 0 {
			if !isNameStart(r) {
				return false
			}
			continue
		}
		if !isNameChar(r) {
			return false
		}
	}
	return name != ""
}

// isNameStart follows the NameStartChar production of XML 1.0, minus ':'.
// A colon is not an ordinary name character here: it separates a namespace
// prefix from the local name, and qualify decides which of the two a given
// key holds.
func isNameStart(r rune) bool {
	switch {
	case r == '_':
		return true
	case r >= 'A' && r <= 'Z', r >= 'a' && r <= 'z':
		return true
	case r >= 0xC0 && r <= 0xD6, r >= 0xD8 && r <= 0xF6, r >= 0xF8 && r <= 0x2FF:
		return true
	case r >= 0x370 && r <= 0x37D, r >= 0x37F && r <= 0x1FFF:
		return true
	case r >= 0x200C && r <= 0x200D, r >= 0x2070 && r <= 0x218F:
		return true
	case r >= 0x2C00 && r <= 0x2FEF, r >= 0x3001 && r <= 0xD7FF:
		return true
	case r >= 0xF900 && r <= 0xFDCF, r >= 0xFDF0 && r <= 0xFFFD:
		return true
	case r >= 0x10000 && r <= 0xEFFFF:
		return true
	}
	return false
}

// isNameChar follows the NameChar production of XML 1.0.
func isNameChar(r rune) bool {
	switch {
	case isNameStart(r):
		return true
	case r == '-' || r == '.' || r == 0xB7:
		return true
	case r >= '0' && r <= '9':
		return true
	case r >= 0x300 && r <= 0x36F, r >= 0x203F && r <= 0x2040:
		return true
	}
	return false
}
