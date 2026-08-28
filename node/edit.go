package node

import (
	"fmt"
	"strings"
)

// SetByPath replaces the value a path addresses and reports whether it landed.
//
// The path is written the way Match.Path writes it, opening with this node's
// own key: "root.meta.version", "root.groups[0].members[1].name".
//
// A last segment naming a key the object does not carry is added, so a path
// both changes a field and creates one. Every segment before the last has to
// exist already.
//
// An object handed to a list element is wrapped the way the decoders wrap one,
// so the element reads back as the object rather than as its first field.
func (n *Node) SetByPath(path string, value *Value) bool {
	if value == nil {
		return false
	}
	steps, ok := n.pathSteps(path)
	if !ok {
		return false
	}

	// The path names this node itself.
	if len(steps) == 0 {
		n.Value = value
		adopt(n, value)
		return true
	}

	owner, container, ok := walkSteps(n, steps[:len(steps)-1])
	if !ok || container == nil {
		return false
	}
	last := steps[len(steps)-1]

	if last.isIndex {
		if container.Type != TypeArray {
			return false
		}
		if last.index < 0 || last.index >= len(container.Array) {
			return false
		}
		container.Array[last.index] = asElement(value, last.index)
		return true
	}

	if container.Type != TypeObject {
		return false
	}
	if child := childByKey(container.Node, last.key); child != nil {
		child.Value = value
		adopt(child, value)
		return true
	}

	child := &Node{Key: last.key, Value: value, Parent: owner}
	adopt(child, value)
	appendChild(container, child)
	return true
}

// RemoveByPath removes what a path addresses, an object member or a list
// element, and reports whether it was there.
//
// The node the path is called on cannot remove itself, so a path naming only
// that node reports false.
func (n *Node) RemoveByPath(path string) bool {
	steps, ok := n.pathSteps(path)
	if !ok || len(steps) == 0 {
		return false
	}

	_, container, ok := walkSteps(n, steps[:len(steps)-1])
	if !ok || container == nil {
		return false
	}
	last := steps[len(steps)-1]

	if last.isIndex {
		if container.Type != TypeArray {
			return false
		}
		if last.index < 0 || last.index >= len(container.Array) {
			return false
		}
		container.Array = append(container.Array[:last.index], container.Array[last.index+1:]...)
		return true
	}

	if container.Type != TypeObject {
		return false
	}
	return unlinkChild(container, last.key)
}

// Replace gives every node under this one whose key is the given key a copy of
// the value, and reports how many were changed.
//
// Each target gets its own copy, so changing one of them afterwards does not
// change the rest. Where one match sits inside another, the outer one wins: its
// whole value is replaced, and what was under it goes with it.
func (n *Node) Replace(key string, value *Value) int {
	if value == nil {
		return 0
	}
	return n.ReplaceFunc(
		func(candidate *Node) bool { return candidate.Key == key },
		func(target *Node) {
			replacement := value.Clone()
			target.Value = replacement
			adopt(target, replacement)
		})
}

// ReplaceFunc hands every node the predicate accepts to replace, and reports
// how many it handed over.
//
// The matches are collected before the first change, so the walk is not
// disturbed by what replace does to the tree.
func (n *Node) ReplaceFunc(predicate func(*Node) bool, replace func(*Node)) int {
	if n == nil || predicate == nil || replace == nil {
		return 0
	}
	matches := n.SearchFunc(predicate)
	for _, match := range matches {
		replace(match.Node)
	}
	return len(matches)
}

// Rename gives every key under this one that reads oldKey the name newKey, and
// reports how many were renamed.
func (n *Node) Rename(oldKey, newKey string) int {
	if oldKey == newKey {
		return 0
	}
	return n.ReplaceFunc(
		func(candidate *Node) bool { return candidate.Key == oldKey },
		func(target *Node) { target.Key = newKey })
}

// Remove deletes every object member under this one whose key is the given key,
// and reports how many were deleted.
//
// A list element is addressed by its index rather than by a key, so elements
// are left alone; RemoveByPath removes one of those.
func (n *Node) Remove(key string) int {
	if n == nil {
		return 0
	}
	return removeUnder(n.Value, key, 0)
}

// removeUnder deletes matching members from every object under a value.
func removeUnder(v *Value, key string, depth int) int {
	if v == nil || depth > maxSearchDepth {
		return 0
	}

	removed := 0
	switch v.Type {
	case TypeObject:
		for unlinkChild(v, key) {
			removed++
		}
		for child := v.Node; child != nil; child = child.Next {
			removed += removeUnder(child.Value, key, depth+1)
		}
	case TypeArray:
		for _, item := range v.Array {
			if item == nil {
				continue
			}
			if item.Node != nil {
				removed += removeUnder(item.Node.Value, key, depth+1)
				continue
			}
			removed += removeUnder(item, key, depth+1)
		}
	}
	return removed
}

// ReplaceText rewrites every occurrence of old inside every text value under
// this node, and reports how many values it changed.
//
// Text inside lists counts. Keys do not, Rename changes those, and neither do
// numbers and booleans, so a search for "1" cannot turn a count into text.
func (n *Node) ReplaceText(old, replacement string) int {
	if n == nil || old == "" || old == replacement {
		return 0
	}
	return replaceTextIn(n.Value, old, replacement, 0)
}

func replaceTextIn(v *Value, old, replacement string, depth int) int {
	if v == nil || depth > maxSearchDepth {
		return 0
	}

	changed := 0
	switch v.Type {
	case TypeString:
		if strings.Contains(v.Worth, old) {
			v.Worth = strings.ReplaceAll(v.Worth, old, replacement)
			changed++
		}
	case TypeObject:
		for child := v.Node; child != nil; child = child.Next {
			changed += replaceTextIn(child.Value, old, replacement, depth+1)
		}
	case TypeArray:
		for _, item := range v.Array {
			if item == nil {
				continue
			}
			if item.Node != nil {
				changed += replaceTextIn(item.Node.Value, old, replacement, depth+1)
				continue
			}
			changed += replaceTextIn(item, old, replacement, depth+1)
		}
	}
	return changed
}

// adopt points a value's child chain back at the node that now holds it.
func adopt(owner *Node, value *Value) {
	if value == nil {
		return
	}
	for child := value.Node; child != nil; child = child.Next {
		child.Parent = owner
	}
}

// asElement prepares a value to sit in a list. An object needs the wrapper node
// the decoders give it, or a reader would take its first field for the whole
// element.
func asElement(value *Value, index int) *Value {
	if value.Type != TypeObject || value.Node == nil {
		return value
	}
	wrapper := &Node{Key: fmt.Sprintf("item%d", index), Value: value}
	adopt(wrapper, value)
	return &Value{Type: TypeObject, Node: wrapper}
}

// appendChild hangs a node at the end of an object's members. The caller has
// already pointed the child at the node that owns the container, so an object
// that had no members yet ends up as sound as one that did.
func appendChild(container *Value, child *Node) {
	child.Next = nil
	if container.Node == nil {
		child.Prev = nil
		container.Node = child
		return
	}
	last := container.Node
	for last.Next != nil {
		last = last.Next
	}
	last.Next = child
	child.Prev = last
}

// unlinkChild takes the first member with the given key out of an object and
// reports whether there was one.
func unlinkChild(container *Value, key string) bool {
	child := childByKey(container.Node, key)
	if child == nil {
		return false
	}

	if child.Prev != nil {
		child.Prev.Next = child.Next
	} else {
		container.Node = child.Next
	}
	if child.Next != nil {
		child.Next.Prev = child.Prev
	}

	child.Parent = nil
	child.Next = nil
	child.Prev = nil
	return true
}
