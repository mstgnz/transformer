package node

import (
	"testing"
)

func str(s string) *Value { return &Value{Type: TypeString, Worth: s} }
func num(s string) *Value { return &Value{Type: TypeNumber, Worth: s} }

// editable builds a tree with the same key at several depths, inside objects,
// inside a list of objects, and a list of lists to reach past.
//
//	root
//	├── a     = 1
//	├── b     { a: 2, c: "staging.example.com" }
//	├── list  [ {a: 3}, {a: 4} ]
//	├── m     [[1, 2]]
//	└── empty {}
func editable() *Node {
	b := &Node{Key: "b", Value: &Value{Type: TypeObject}}
	b.Value.Node = &Node{Key: "a", Value: num("2"), Parent: b,
		Next: &Node{Key: "c", Value: str("staging.example.com"), Parent: b}}
	b.Value.Node.Next.Prev = b.Value.Node

	item0 := &Node{Key: "item0", Value: &Value{Type: TypeObject}}
	item0.Value.Node = &Node{Key: "a", Value: num("3"), Parent: item0}
	item1 := &Node{Key: "item1", Value: &Value{Type: TypeObject}}
	item1.Value.Node = &Node{Key: "a", Value: num("4"), Parent: item1}

	root := &Node{Key: "root", Value: &Value{Type: TypeObject}}

	a := &Node{Key: "a", Value: num("1"), Parent: root}
	list := &Node{Key: "list", Parent: root, Value: &Value{Type: TypeArray, Array: []*Value{
		{Type: TypeObject, Node: item0},
		{Type: TypeObject, Node: item1},
	}}}
	m := &Node{Key: "m", Parent: root, Value: &Value{Type: TypeArray, Array: []*Value{
		{Type: TypeArray, Array: []*Value{num("1"), num("2")}},
	}}}
	empty := &Node{Key: "empty", Value: &Value{Type: TypeObject}, Parent: root}

	b.Parent = root
	a.Next, b.Prev = b, a
	b.Next, list.Prev = list, b
	list.Next, m.Prev = m, list
	m.Next, empty.Prev = empty, m
	root.Value.Node = a

	return root
}

// worthAt reads the text a path leads to, for compact assertions.
func worthAt(t *testing.T, root *Node, path string) string {
	t.Helper()
	v, ok := root.GetValueByPath(path)
	if !ok {
		return "<none>"
	}
	return v.Worth
}

// ---- SetByPath ----

func TestSetByPath(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		value *Value
		want  bool
		check func(*testing.T, *Node)
	}{
		{
			name: "existing member", path: "root.b.c", value: str("new"), want: true,
			check: func(t *testing.T, n *Node) {
				if got := worthAt(t, n, "root.b.c"); got != "new" {
					t.Errorf("b.c = %q, want new", got)
				}
			},
		},
		{
			name: "member that was not there", path: "root.b.added", value: str("yes"), want: true,
			check: func(t *testing.T, n *Node) {
				if got := worthAt(t, n, "root.b.added"); got != "yes" {
					t.Errorf("b.added = %q, want yes", got)
				}
				// It goes at the end, so the members that were there keep their order.
				var keys []string
				for c := n.GetNodeByPath("root.b").Value.Node; c != nil; c = c.Next {
					keys = append(keys, c.Key)
				}
				if len(keys) != 3 || keys[0] != "a" || keys[2] != "added" {
					t.Errorf("members = %v, want a, c, added", keys)
				}
			},
		},
		{
			name: "first member of an empty object", path: "root.empty.first", value: num("1"), want: true,
			check: func(t *testing.T, n *Node) {
				first := n.GetNodeByPath("root.empty.first")
				if first == nil {
					t.Fatal("the member was not added")
				}
				if first.Parent == nil || first.Parent.Key != "empty" {
					t.Errorf("the new member was not wired to its parent: %v", first.Parent)
				}
			},
		},
		{
			name: "list element", path: "root.list[1]", value: str("flat"), want: true,
			check: func(t *testing.T, n *Node) {
				if got := worthAt(t, n, "root.list[1]"); got != "flat" {
					t.Errorf("list[1] = %q, want flat", got)
				}
			},
		},
		{
			name: "element of a list inside a list", path: "root.m[0][1]", value: num("42"), want: true,
			check: func(t *testing.T, n *Node) {
				if got := worthAt(t, n, "root.m[0][1]"); got != "42" {
					t.Errorf("m[0][1] = %q, want 42", got)
				}
			},
		},
		{
			name: "member inside a list element", path: "root.list[0].a", value: num("30"), want: true,
			check: func(t *testing.T, n *Node) {
				if got := worthAt(t, n, "root.list[0].a"); got != "30" {
					t.Errorf("list[0].a = %q, want 30", got)
				}
			},
		},
		{
			name: "the node itself", path: "root", value: str("flattened"), want: true,
			check: func(t *testing.T, n *Node) {
				if n.Value.Type != TypeString || n.Value.Worth != "flattened" {
					t.Errorf("root = %+v, want the string flattened", n.Value)
				}
			},
		},
		{name: "missing step on the way", path: "root.nope.deeper", value: str("x"), want: false},
		{name: "index past the end", path: "root.list[9]", value: str("x"), want: false},
		{name: "index into an object", path: "root.b[0]", value: str("x"), want: false},
		{name: "member of a scalar", path: "root.a.deeper", value: str("x"), want: false},
		{name: "path that names another node", path: "other.a", value: str("x"), want: false},
		{name: "empty path", path: "", value: str("x"), want: false},
		{name: "nil value", path: "root.a", value: nil, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := editable()
			if got := root.SetByPath(tt.path, tt.value); got != tt.want {
				t.Fatalf("SetByPath(%q) = %v, want %v", tt.path, got, tt.want)
			}
			if tt.check != nil {
				tt.check(t, root)
			}
			if err := root.Validate(); err != nil {
				t.Errorf("the tree is no longer sound after the edit: %v", err)
			}
		})
	}
}

// TestSetByPath_ObjectIntoList covers the shape a list element has to have. An
// object stored raw would be read back as its own first field.
func TestSetByPath_ObjectIntoList(t *testing.T) {
	root := editable()
	object := &Value{Type: TypeObject, Node: &Node{Key: "z", Value: num("7")}}

	if !root.SetByPath("root.list[0]", object) {
		t.Fatal("SetByPath did not land")
	}
	if got := worthAt(t, root, "root.list[0].z"); got != "7" {
		t.Errorf("list[0].z = %q, want 7", got)
	}
	if got := root.GetNodeByPath("root.list[0]"); got == nil {
		t.Error("the element no longer addresses a node")
	}
}

// ---- RemoveByPath ----

func TestRemoveByPath(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		want  bool
		check func(*testing.T, *Node)
	}{
		{
			name: "member in the middle", path: "root.b.a", want: true,
			check: func(t *testing.T, n *Node) {
				if n.GetNodeByPath("root.b.a") != nil {
					t.Error("b.a is still there")
				}
				if worthAt(t, n, "root.b.c") != "staging.example.com" {
					t.Error("the member beside it was disturbed")
				}
			},
		},
		{
			name: "first member", path: "root.a", want: true,
			check: func(t *testing.T, n *Node) {
				if n.Value.Node == nil || n.Value.Node.Key != "b" {
					t.Errorf("the object's first member is %v, want b", n.Value.Node)
				}
				if n.Value.Node.Prev != nil {
					t.Error("the new first member still points back at the old one")
				}
			},
		},
		{
			name: "last member", path: "root.empty", want: true,
			check: func(t *testing.T, n *Node) {
				last := n.Value.Node
				for last.Next != nil {
					last = last.Next
				}
				if last.Key != "m" {
					t.Errorf("the object's last member is %v, want m", last.Key)
				}
			},
		},
		{
			name: "list element", path: "root.list[0]", want: true,
			check: func(t *testing.T, n *Node) {
				list := n.GetNodeByPath("root.list")
				if len(list.Value.Array) != 1 {
					t.Fatalf("the list holds %d values, want 1", len(list.Value.Array))
				}
				if worthAt(t, n, "root.list[0].a") != "4" {
					t.Error("the wrong element was removed")
				}
			},
		},
		{
			name: "element of a list inside a list", path: "root.m[0][0]", want: true,
			check: func(t *testing.T, n *Node) {
				if worthAt(t, n, "root.m[0][0]") != "2" {
					t.Error("the wrong inner value was removed")
				}
			},
		},
		{name: "member that is not there", path: "root.nope", want: false},
		{name: "index past the end", path: "root.list[9]", want: false},
		{name: "the node itself", path: "root", want: false},
		{name: "empty path", path: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := editable()
			if got := root.RemoveByPath(tt.path); got != tt.want {
				t.Fatalf("RemoveByPath(%q) = %v, want %v", tt.path, got, tt.want)
			}
			if tt.check != nil {
				tt.check(t, root)
			}
			if err := root.Validate(); err != nil {
				t.Errorf("the tree is no longer sound after the removal: %v", err)
			}
		})
	}
}

// ---- Replace ----

func TestReplace(t *testing.T) {
	root := editable()

	if got := root.Replace("a", num("99")); got != 4 {
		t.Fatalf("Replace changed %d nodes, want 4", got)
	}
	for _, path := range []string{"root.a", "root.b.a", "root.list[0].a", "root.list[1].a"} {
		if got := worthAt(t, root, path); got != "99" {
			t.Errorf("%s = %q, want 99", path, got)
		}
	}
	if err := root.Validate(); err != nil {
		t.Errorf("the tree is no longer sound: %v", err)
	}
}

// TestReplace_EachTargetGetsItsOwnCopy is what keeps a replacement from turning
// into an alias: changing one of the results must not change the rest.
func TestReplace_EachTargetGetsItsOwnCopy(t *testing.T) {
	root := editable()
	root.Replace("a", num("5"))

	first := root.GetNodeByPath("root.a")
	first.Value.Worth = "changed"

	if got := worthAt(t, root, "root.b.a"); got != "5" {
		t.Errorf("b.a = %q, want 5: the targets are sharing one value", got)
	}
	if got := worthAt(t, root, "root.list[0].a"); got != "5" {
		t.Errorf("list[0].a = %q, want 5: the targets are sharing one value", got)
	}
}

func TestReplace_CopiesNestedLists(t *testing.T) {
	root := editable()
	nested := &Value{Type: TypeArray, Array: []*Value{
		{Type: TypeArray, Array: []*Value{num("7"), num("8")}},
	}}

	root.Replace("a", nested)
	if got := worthAt(t, root, "root.b.a[0][1]"); got != "8" {
		t.Errorf("b.a[0][1] = %q, want 8: the copy dropped the inner list", got)
	}
}

func TestReplace_NoMatchesAndNilValue(t *testing.T) {
	root := editable()
	if got := root.Replace("nope", num("1")); got != 0 {
		t.Errorf("Replace on a key that is not there = %d, want 0", got)
	}
	if got := root.Replace("a", nil); got != 0 {
		t.Errorf("Replace with no value = %d, want 0", got)
	}
}

func TestReplaceFunc(t *testing.T) {
	root := editable()

	got := root.ReplaceFunc(
		func(n *Node) bool { return n.Type() == TypeNumber },
		func(n *Node) { n.Value = str("N") })
	if got != 4 {
		t.Errorf("ReplaceFunc changed %d nodes, want 4", got)
	}

	// A plain value inside a list has no node, so it is left alone.
	if worthAt(t, root, "root.m[0][0]") != "1" {
		t.Error("a plain value inside a list was changed through a node walk")
	}

	var nilNode *Node
	if got := nilNode.ReplaceFunc(func(*Node) bool { return true }, func(*Node) {}); got != 0 {
		t.Errorf("ReplaceFunc on nil = %d, want 0", got)
	}
	if got := root.ReplaceFunc(nil, func(*Node) {}); got != 0 {
		t.Errorf("ReplaceFunc without a predicate = %d, want 0", got)
	}
	if got := root.ReplaceFunc(func(*Node) bool { return true }, nil); got != 0 {
		t.Errorf("ReplaceFunc without a replacement = %d, want 0", got)
	}
}

// ---- Rename ----

func TestRename(t *testing.T) {
	root := editable()

	if got := root.Rename("a", "id"); got != 4 {
		t.Fatalf("Rename changed %d keys, want 4", got)
	}
	if root.GetNodeByPath("root.a") != nil {
		t.Error("the old key is still addressable")
	}
	if worthAt(t, root, "root.id") != "1" {
		t.Error("the renamed key does not carry the old value")
	}
	if worthAt(t, root, "root.list[1].id") != "4" {
		t.Error("a key inside a list element was not renamed")
	}
	if err := root.Validate(); err != nil {
		t.Errorf("the tree is no longer sound: %v", err)
	}

	if got := root.Rename("id", "id"); got != 0 {
		t.Errorf("renaming a key to itself changed %d keys, want 0", got)
	}
	if got := root.Rename("nope", "x"); got != 0 {
		t.Errorf("renaming a key that is not there changed %d keys, want 0", got)
	}
}

// ---- Remove ----

func TestRemove(t *testing.T) {
	root := editable()

	if got := root.Remove("a"); got != 4 {
		t.Fatalf("Remove deleted %d members, want 4", got)
	}
	for _, path := range []string{"root.a", "root.b.a", "root.list[0].a", "root.list[1].a"} {
		if root.GetNodeByPath(path) != nil {
			t.Errorf("%s is still there", path)
		}
	}
	if worthAt(t, root, "root.b.c") != "staging.example.com" {
		t.Error("a member that should have been left alone is gone")
	}
	if err := root.Validate(); err != nil {
		t.Errorf("the tree is no longer sound: %v", err)
	}
}

// TestRemove_RepeatedMembers covers an object that carries the same key twice,
// which XML produces before the elements are grouped.
func TestRemove_RepeatedMembers(t *testing.T) {
	root := &Node{Key: "root", Value: &Value{Type: TypeObject}}
	first := &Node{Key: "dup", Value: num("1"), Parent: root}
	second := &Node{Key: "dup", Value: num("2"), Parent: root}
	keep := &Node{Key: "keep", Value: num("3"), Parent: root}
	first.Next, second.Prev = second, first
	second.Next, keep.Prev = keep, second
	root.Value.Node = first

	if got := root.Remove("dup"); got != 2 {
		t.Fatalf("Remove deleted %d members, want 2", got)
	}
	if root.Value.Node == nil || root.Value.Node.Key != "keep" {
		t.Errorf("what is left is %v, want keep", root.Value.Node)
	}
	if err := root.Validate(); err != nil {
		t.Errorf("the tree is no longer sound: %v", err)
	}
}

func TestRemove_NoMatches(t *testing.T) {
	root := editable()
	if got := root.Remove("nope"); got != 0 {
		t.Errorf("Remove = %d, want 0", got)
	}
	var nilNode *Node
	if got := nilNode.Remove("a"); got != 0 {
		t.Errorf("Remove on nil = %d, want 0", got)
	}
}

// ---- ReplaceText ----

func TestReplaceText(t *testing.T) {
	root := editable()
	root.SetByPath("root.tags", &Value{Type: TypeArray, Array: []*Value{
		str("staging.example.com"),
		str("keep"),
	}})
	root.SetByPath("root.deep", &Value{Type: TypeArray, Array: []*Value{
		{Type: TypeArray, Array: []*Value{str("staging.example.com")}},
	}})

	if got := root.ReplaceText("staging.", ""); got != 3 {
		t.Fatalf("ReplaceText changed %d values, want 3", got)
	}
	if got := worthAt(t, root, "root.b.c"); got != "example.com" {
		t.Errorf("b.c = %q, want example.com", got)
	}
	if got := worthAt(t, root, "root.tags[0]"); got != "example.com" {
		t.Errorf("tags[0] = %q, want example.com: text inside a list was skipped", got)
	}
	if got := worthAt(t, root, "root.deep[0][0]"); got != "example.com" {
		t.Errorf("deep[0][0] = %q: text inside a list of lists was skipped", got)
	}
	if got := worthAt(t, root, "root.tags[1]"); got != "keep" {
		t.Errorf("tags[1] = %q, want keep", got)
	}
}

// TestReplaceText_LeavesEverythingElseAlone keeps the substitution to text: a
// key, a number or a boolean that happens to contain the text is not rewritten.
func TestReplaceText_LeavesEverythingElseAlone(t *testing.T) {
	root := &Node{Key: "root", Value: &Value{Type: TypeObject}}
	root.Value.Node = &Node{Key: "1", Value: num("1"), Parent: root,
		Next: &Node{Key: "flag", Value: &Value{Type: TypeBoolean, Worth: "true"}, Parent: root}}

	if got := root.ReplaceText("1", "one"); got != 0 {
		t.Errorf("ReplaceText changed %d values, want 0", got)
	}
	if root.Value.Node.Key != "1" {
		t.Error("a key was rewritten")
	}
	if root.Value.Node.Value.Worth != "1" {
		t.Error("a number was rewritten")
	}
	if root.Value.Node.Next.Value.Worth != "true" {
		t.Error("a boolean was rewritten")
	}
}

func TestReplaceText_Guards(t *testing.T) {
	root := editable()
	if got := root.ReplaceText("", "x"); got != 0 {
		t.Errorf("ReplaceText with no text to find = %d, want 0", got)
	}
	if got := root.ReplaceText("same", "same"); got != 0 {
		t.Errorf("ReplaceText with the same text on both sides = %d, want 0", got)
	}
	var nilNode *Node
	if got := nilNode.ReplaceText("a", "b"); got != 0 {
		t.Errorf("ReplaceText on nil = %d, want 0", got)
	}
}

// ---- Copying and comparing ----
//
// Replace hands every target its own copy, so these two are what keep it honest.

func TestValueClone_CopiesNestedLists(t *testing.T) {
	original := &Value{Type: TypeArray, Array: []*Value{
		{Type: TypeArray, Array: []*Value{num("1"), num("2")}},
	}}

	clone := original.Clone()
	if len(clone.Array) != 1 || len(clone.Array[0].Array) != 2 {
		t.Fatalf("the inner list was dropped: %+v", clone)
	}
	if clone.Array[0].Array[0].Worth != "1" {
		t.Errorf("inner value = %q, want 1", clone.Array[0].Array[0].Worth)
	}

	clone.Array[0].Array[0].Worth = "changed"
	if original.Array[0].Array[0].Worth != "1" {
		t.Error("the copy shares its inner list with the original")
	}

	var nilValue *Value
	if nilValue.Clone() != nil {
		t.Error("cloning nothing produced something")
	}
}

func TestNodeClone_IsIndependent(t *testing.T) {
	root := editable()
	clone := root.Clone()

	if !root.Equal(clone) {
		t.Fatal("a copy does not compare equal to what it was copied from")
	}

	clone.Replace("a", num("0"))
	if worthAt(t, root, "root.a") != "1" {
		t.Error("editing the copy changed the original")
	}
}

// TestEqual_ComparesEveryMember covers what the earlier version missed: it
// looked at an object's first member only, which made {"a":1,"b":2} equal to
// {"a":1}, and it never compared lists inside lists.
func TestEqual_ComparesEveryMember(t *testing.T) {
	twoMembers := func(second *Value) *Node {
		n := &Node{Key: "root", Value: &Value{Type: TypeObject}}
		n.Value.Node = &Node{Key: "a", Value: num("1"), Parent: n}
		if second != nil {
			n.Value.Node.Next = &Node{Key: "b", Value: second, Parent: n}
			n.Value.Node.Next.Prev = n.Value.Node
		}
		return n
	}

	if twoMembers(num("2")).Equal(twoMembers(num("3"))) {
		t.Error("two objects differing in their second member compared equal")
	}
	if twoMembers(num("2")).Equal(twoMembers(nil)) {
		t.Error("an object compared equal to one with a member missing")
	}
	if twoMembers(nil).Equal(twoMembers(num("2"))) {
		t.Error("an object compared equal to one with an extra member")
	}
	if !twoMembers(num("2")).Equal(twoMembers(num("2"))) {
		t.Error("two objects holding the same members compared unequal")
	}

	nested := func(inner ...*Value) *Node {
		return &Node{Key: "root", Value: &Value{Type: TypeArray, Array: []*Value{
			{Type: TypeArray, Array: inner},
		}}}
	}
	if nested(num("1"), num("2")).Equal(nested(num("3"), num("4"))) {
		t.Error("two lists of lists holding different values compared equal")
	}
	if nested(num("1")).Equal(nested(num("1"), num("2"))) {
		t.Error("two lists of lists of different lengths compared equal")
	}
	if !nested(num("1"), num("2")).Equal(nested(num("1"), num("2"))) {
		t.Error("two matching lists of lists compared unequal")
	}

	var nilNode *Node
	if !nilNode.Equal(nil) {
		t.Error("nothing does not compare equal to nothing")
	}
	if nilNode.Equal(twoMembers(nil)) {
		t.Error("nothing compared equal to something")
	}
}
