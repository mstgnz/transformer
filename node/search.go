package node

import (
	"fmt"
	"strconv"
	"strings"
)

// maxSearchDepth bounds a walk. A tree built by the decoders cannot be cyclic,
// but one assembled by hand can be, and a search should fail rather than hang.
const maxSearchDepth = 1000

// Match is one node a search found, together with the route that reaches it.
//
// The route is what tells two hits on the same key apart: a document can carry
// "name" at the top level and again inside every member of every group, and
// which one is meant is a question the path answers.
type Match struct {
	// Node is the node that matched.
	Node *Node

	// Path addresses the node from the one that was searched, written the way
	// a reader would write it and opening with the searched node's own key:
	// "root.groups[0].members[1].profile.name". GetNodeByPath reads the same
	// form back.
	Path string

	// Depth is how many steps that route takes. A direct child of the searched
	// node is at depth 1, and an array index is a step of its own, so
	// "groups[0].name" is at depth 3.
	Depth int
}

// String renders a match as its path, its depth and its value.
func (m Match) String() string {
	if m.Node == nil {
		return m.Path
	}
	if m.Node.Value != nil && m.Node.Value.Worth != "" {
		return fmt.Sprintf("%s (depth %d) = %s", m.Path, m.Depth, m.Node.Value.Worth)
	}
	return fmt.Sprintf("%s (depth %d) [%s]", m.Path, m.Depth, m.Node.Type())
}

// Search returns every node under this one whose key is the given key, in
// document order, each with the path and depth that reach it.
//
// The searched node itself is not a candidate: the search looks inside it.
// Keys inside arrays, and inside arrays of arrays, are reached like any other.
func (n *Node) Search(key string) []Match {
	return n.SearchFunc(func(candidate *Node) bool { return candidate.Key == key })
}

// SearchFunc returns every node under this one that the predicate accepts.
func (n *Node) SearchFunc(predicate func(*Node) bool) []Match {
	if n == nil || predicate == nil {
		return nil
	}
	var matches []Match
	walkNode(n, n.Key, 0, predicate, &matches)
	return matches
}

// Depths returns the depths at which a key appears, in ascending order and
// without repeats. It answers "which levels is this key on" without having to
// read every match.
func (n *Node) Depths(key string) []int {
	seen := make(map[int]bool)
	var depths []int
	for _, match := range n.Search(key) {
		if seen[match.Depth] {
			continue
		}
		seen[match.Depth] = true
		depths = append(depths, match.Depth)
	}
	for i := 1; i < len(depths); i++ {
		for j := i; j > 0 && depths[j] < depths[j-1]; j-- {
			depths[j], depths[j-1] = depths[j-1], depths[j]
		}
	}
	return depths
}

// walkNode visits a node and everything under it.
func walkNode(n *Node, path string, depth int, predicate func(*Node) bool, out *[]Match) {
	if n == nil || depth > maxSearchDepth {
		return
	}

	if depth > 0 && predicate(n) {
		*out = append(*out, Match{Node: n, Path: path, Depth: depth})
	}

	if n.Value == nil {
		return
	}

	switch n.Value.Type {
	case TypeObject:
		for child := n.Value.Node; child != nil; child = child.Next {
			walkNode(child, joinPath(path, child.Key), depth+1, predicate, out)
		}
	case TypeArray:
		walkValues(n.Value.Array, path, depth, predicate, out)
	}
}

// walkValues visits the elements of an array. An element that is a plain value
// has no node of its own and so cannot match, but one that holds an object or
// another array is walked like anything else. The element's index takes the
// place of its key in the route, because the index is how a reader addresses
// it; the node an element holds keeps whatever key it was given.
func walkValues(values []*Value, path string, depth int, predicate func(*Node) bool, out *[]Match) {
	if depth > maxSearchDepth {
		return
	}
	for i, item := range values {
		if item == nil {
			continue
		}
		indexed := fmt.Sprintf("%s[%d]", path, i)
		switch {
		case item.Node != nil:
			walkNode(item.Node, indexed, depth+1, predicate, out)
		case item.Type == TypeArray:
			walkValues(item.Array, indexed, depth+1, predicate, out)
		}
	}
}

func joinPath(path, key string) string {
	if path == "" {
		return key
	}
	return path + "." + key
}

// GetNode returns this node and every node under it whose key is the given key.
//
// Search is the fuller answer: it reports where each hit sits. This one is kept
// for callers that only want the nodes.
func (n *Node) GetNode(key string) []*Node {
	if n == nil {
		return nil
	}

	var nodes []*Node
	if n.Key == key {
		nodes = append(nodes, n)
	}
	for _, match := range n.Search(key) {
		nodes = append(nodes, match.Node)
	}
	return nodes
}

// FindNodes returns this node and every node under it that the predicate
// accepts.
//
// A plain value inside an array has no node of its own, so it is not a
// candidate; search its list instead.
func (n *Node) FindNodes(predicate func(*Node) bool) []*Node {
	if n == nil || predicate == nil {
		return nil
	}

	var nodes []*Node
	if predicate(n) {
		nodes = append(nodes, n)
	}
	for _, match := range n.SearchFunc(predicate) {
		nodes = append(nodes, match.Node)
	}
	return nodes
}

// GetNodeByPath returns the node a path addresses, or nil when the path leads
// nowhere.
//
// The path is written the way Match.Path writes it, opening with this node's
// own key: "root.meta.version", "root.groups[0].members[1].name". A path that
// ends at a plain value inside a list addresses no node and returns nil;
// GetValueByPath answers for those.
//
// A key that itself contains a dot or a bracket cannot be addressed this way.
func (n *Node) GetNodeByPath(path string) *Node {
	node, _, ok := n.resolvePath(path)
	if !ok {
		return nil
	}
	return node
}

// GetValueByPath returns the value a path addresses, and whether the path led
// anywhere. Unlike GetNodeByPath it also reaches plain values inside a list,
// which have no node of their own: "root.matrix[0][2]", "root.meta.tags[1]".
func (n *Node) GetValueByPath(path string) (*Value, bool) {
	_, value, ok := n.resolvePath(path)
	return value, ok
}

// pathSteps parses a path and strips the opening segment, which names this
// node itself so that a path taken from a Match reads back unchanged.
func (n *Node) pathSteps(path string) ([]pathStep, bool) {
	if n == nil {
		return nil, false
	}
	steps, ok := parsePath(path)
	if !ok || len(steps) == 0 {
		return nil, false
	}
	if steps[0].isIndex || steps[0].key != n.Key {
		return nil, false
	}
	return steps[1:], true
}

// walkSteps follows a path from a starting node and reports where it ended:
// the node (nil when the path lands on a plain value inside a list), the value,
// and whether every step led somewhere.
func walkSteps(start *Node, steps []pathStep) (*Node, *Value, bool) {
	if start == nil {
		return nil, nil, false
	}

	current := start
	value := start.Value

	for _, step := range steps {
		if step.isIndex {
			item, ok := elementAt(value, step.index)
			if !ok {
				return nil, nil, false
			}
			if item.Node != nil {
				current, value = item.Node, item.Node.Value
			} else {
				current, value = nil, item
			}
			continue
		}

		if value == nil || value.Type != TypeObject {
			return nil, nil, false
		}
		child := childByKey(value.Node, step.key)
		if child == nil {
			return nil, nil, false
		}
		current, value = child, child.Value
	}

	return current, value, true
}

// elementAt returns the element a list holds at an index.
func elementAt(v *Value, index int) (*Value, bool) {
	if v == nil || v.Type != TypeArray {
		return nil, false
	}
	if index < 0 || index >= len(v.Array) {
		return nil, false
	}
	item := v.Array[index]
	if item == nil {
		return nil, false
	}
	return item, true
}

// resolvePath walks a path and reports the node it ended on (nil for a plain
// value inside a list), the value it ended on, and whether it led anywhere.
func (n *Node) resolvePath(path string) (*Node, *Value, bool) {
	steps, ok := n.pathSteps(path)
	if !ok {
		return nil, nil, false
	}
	return walkSteps(n, steps)
}

func childByKey(first *Node, key string) *Node {
	for current := first; current != nil; current = current.Next {
		if current.Key == key {
			return current
		}
	}
	return nil
}

// pathStep is one segment of a path: a key, or an index into a list.
type pathStep struct {
	key     string
	index   int
	isIndex bool
}

// parsePath splits "groups[0].members[1].name" into its steps.
func parsePath(path string) ([]pathStep, bool) {
	if path == "" {
		return nil, false
	}

	var steps []pathStep
	rest := path

	for rest != "" {
		if strings.HasPrefix(rest, "[") {
			end := strings.Index(rest, "]")
			if end < 0 {
				return nil, false
			}
			index, err := strconv.Atoi(rest[1:end])
			if err != nil {
				return nil, false
			}
			steps = append(steps, pathStep{index: index, isIndex: true})
			rest = rest[end+1:]
			if strings.HasPrefix(rest, ".") {
				rest = rest[1:]
				if rest == "" {
					return nil, false
				}
			}
			continue
		}

		end := strings.IndexAny(rest, ".[")
		if end < 0 {
			steps = append(steps, pathStep{key: rest})
			break
		}
		if end == 0 {
			return nil, false
		}
		steps = append(steps, pathStep{key: rest[:end]})
		if rest[end] == '.' {
			rest = rest[end+1:]
			if rest == "" {
				return nil, false
			}
			continue
		}
		rest = rest[end:]
	}

	return steps, len(steps) > 0
}
