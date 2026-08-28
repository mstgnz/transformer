package node

import (
	"reflect"
	"testing"
	"time"
)

// sample builds a tree that carries the same key at several depths, inside
// objects, inside arrays, and inside arrays of arrays.
//
//	root
//	├── name            = "top"
//	├── tags            = ["a", "b"]
//	├── matrix          = [[1, 2], [3]]
//	├── group
//	│   ├── name        = "inner"
//	│   └── members     = [ {name: "alice", scores: [9, 8]}, {name: "bob"} ]
//	└── deep            = [ [ {name: "buried", detail: {name: "deepest"}} ] ]
func sample() *Node {
	scores := &Value{Type: TypeArray, Array: []*Value{
		{Type: TypeNumber, Worth: "9"},
		{Type: TypeNumber, Worth: "8"},
	}}
	alice := &Node{Key: "item0", Value: &Value{Type: TypeObject, Node: &Node{
		Key:   "name",
		Value: &Value{Type: TypeString, Worth: "alice"},
		Next:  &Node{Key: "scores", Value: scores},
	}}}
	bob := &Node{Key: "item1", Value: &Value{Type: TypeObject, Node: &Node{
		Key:   "name",
		Value: &Value{Type: TypeString, Worth: "bob"},
	}}}

	group := &Node{Key: "group", Value: &Value{Type: TypeObject, Node: &Node{
		Key:   "name",
		Value: &Value{Type: TypeString, Worth: "inner"},
		Next: &Node{Key: "members", Value: &Value{Type: TypeArray, Array: []*Value{
			{Type: TypeObject, Node: alice},
			{Type: TypeObject, Node: bob},
		}}},
	}}}

	buried := &Node{Key: "item0", Value: &Value{Type: TypeObject, Node: &Node{
		Key:   "name",
		Value: &Value{Type: TypeString, Worth: "buried"},
		Next: &Node{Key: "detail", Value: &Value{Type: TypeObject, Node: &Node{
			Key:   "name",
			Value: &Value{Type: TypeString, Worth: "deepest"},
		}}},
	}}}

	root := &Node{Key: "root", Value: &Value{Type: TypeObject}}
	root.Value.Node = &Node{
		Key:   "name",
		Value: &Value{Type: TypeString, Worth: "top"},
		Next: &Node{
			Key: "tags",
			Value: &Value{Type: TypeArray, Array: []*Value{
				{Type: TypeString, Worth: "a"},
				{Type: TypeString, Worth: "b"},
			}},
			Next: &Node{
				Key: "matrix",
				Value: &Value{Type: TypeArray, Array: []*Value{
					{Type: TypeArray, Array: []*Value{
						{Type: TypeNumber, Worth: "1"},
						{Type: TypeNumber, Worth: "2"},
					}},
					{Type: TypeArray, Array: []*Value{
						{Type: TypeNumber, Worth: "3"},
					}},
				}},
				Next: &Node{
					Key:   "group",
					Value: group.Value,
					Next: &Node{
						Key: "deep",
						Value: &Value{Type: TypeArray, Array: []*Value{
							{Type: TypeArray, Array: []*Value{
								{Type: TypeObject, Node: buried},
							}},
						}},
					},
				},
			},
		},
	}
	return root
}

// ---- Search ----

func TestSearch_FindsEveryOccurrence(t *testing.T) {
	got := sample().Search("name")

	want := []struct {
		path  string
		depth int
		worth string
	}{
		{"root.name", 1, "top"},
		{"root.group.name", 2, "inner"},
		{"root.group.members[0].name", 4, "alice"},
		{"root.group.members[1].name", 4, "bob"},
		{"root.deep[0][0].name", 4, "buried"},
		{"root.deep[0][0].detail.name", 5, "deepest"},
	}

	if len(got) != len(want) {
		t.Fatalf("got %d matches, want %d:\n%v", len(got), len(want), got)
	}
	for i, w := range want {
		if got[i].Path != w.path {
			t.Errorf("match %d path = %q, want %q", i, got[i].Path, w.path)
		}
		if got[i].Depth != w.depth {
			t.Errorf("match %d depth = %d, want %d", i, got[i].Depth, w.depth)
		}
		if got[i].Node.Value.Worth != w.worth {
			t.Errorf("match %d value = %q, want %q", i, got[i].Node.Value.Worth, w.worth)
		}
	}
}

// TestSearch_ReachesInsideNestedArrays is the case the old traversal missed
// entirely: it followed an array element that held a node, but never one that
// held another array, so anything under a list of lists was invisible.
func TestSearch_ReachesInsideNestedArrays(t *testing.T) {
	matches := sample().Search("name")
	found := false
	for _, m := range matches {
		if m.Path == "root.deep[0][0].detail.name" {
			found = true
		}
	}
	if !found {
		t.Errorf("a key inside an array of arrays was not found:\n%v", matches)
	}
}

func TestSearch_NoMatches(t *testing.T) {
	if got := sample().Search("nope"); len(got) != 0 {
		t.Errorf("Search(nope) = %v, want none", got)
	}
}

func TestSearch_DoesNotMatchTheSearchedNode(t *testing.T) {
	if got := sample().Search("root"); len(got) != 0 {
		t.Errorf("Search looked at the node it was called on: %v", got)
	}
}

func TestSearch_NilAndEmpty(t *testing.T) {
	var nilNode *Node
	if got := nilNode.Search("x"); got != nil {
		t.Errorf("Search on nil = %v, want nil", got)
	}
	if got := sample().SearchFunc(nil); got != nil {
		t.Errorf("SearchFunc(nil) = %v, want nil", got)
	}
}

func TestSearchFunc(t *testing.T) {
	numbers := sample().SearchFunc(func(n *Node) bool { return n.Type() == TypeNumber })
	if len(numbers) != 0 {
		t.Errorf("plain values inside a list have no node of their own: %v", numbers)
	}

	lists := sample().SearchFunc(func(n *Node) bool { return n.Type() == TypeArray })
	wantPaths := []string{"root.tags", "root.matrix", "root.group.members", "root.group.members[0].scores", "root.deep"}
	var gotPaths []string
	for _, m := range lists {
		gotPaths = append(gotPaths, m.Path)
	}
	if !reflect.DeepEqual(gotPaths, wantPaths) {
		t.Errorf("lists = %v, want %v", gotPaths, wantPaths)
	}
}

// ---- Depths ----

func TestDepths(t *testing.T) {
	if got, want := sample().Depths("name"), []int{1, 2, 4, 5}; !reflect.DeepEqual(got, want) {
		t.Errorf("Depths(name) = %v, want %v", got, want)
	}
	if got := sample().Depths("nope"); len(got) != 0 {
		t.Errorf("Depths(nope) = %v, want none", got)
	}
	if got, want := sample().Depths("scores"), []int{4}; !reflect.DeepEqual(got, want) {
		t.Errorf("Depths(scores) = %v, want %v", got, want)
	}
}

// ---- Paths ----

// TestSearchPathsResolveBack is the contract between the two halves of the
// search: a path a search reports must address the node it came from.
func TestSearchPathsResolveBack(t *testing.T) {
	root := sample()
	for _, m := range root.Search("name") {
		got := root.GetNodeByPath(m.Path)
		if got != m.Node {
			t.Errorf("GetNodeByPath(%q) did not return the node the search reported", m.Path)
		}
	}
}

func TestGetNodeByPath(t *testing.T) {
	root := sample()
	tests := []struct {
		path  string
		want  string // key of the expected node, empty for nil
		worth string
	}{
		{"root", "root", ""},
		{"root.name", "name", "top"},
		{"root.group.name", "name", "inner"},
		{"root.group.members[1].name", "name", "bob"},
		{"root.deep[0][0].detail.name", "name", "deepest"},
		{"root.group.members[0]", "item0", ""},
		{"root.tags", "tags", ""},
		{"root.tags[0]", "", ""}, // a plain value in a list has no node
		{"root.matrix[0][1]", "", ""},
		{"root.nope", "", ""},
		{"root.group.nope", "", ""},
		{"root.group.members[9].name", "", ""},
		{"root.name.deeper", "", ""}, // a scalar has no children
		{"group.name", "", ""},       // the path must open with this node's key
		{"", "", ""},
		{"root.", "", ""},
		{"root.members[", "", ""},
		{"root.members[x]", "", ""},
		{".name", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := root.GetNodeByPath(tt.path)
			if tt.want == "" {
				if got != nil {
					t.Errorf("GetNodeByPath(%q) = %v, want nil", tt.path, got.Key)
				}
				return
			}
			if got == nil {
				t.Fatalf("GetNodeByPath(%q) = nil, want %q", tt.path, tt.want)
			}
			if got.Key != tt.want {
				t.Errorf("GetNodeByPath(%q) key = %q, want %q", tt.path, got.Key, tt.want)
			}
			if tt.worth != "" && got.Value.Worth != tt.worth {
				t.Errorf("GetNodeByPath(%q) value = %q, want %q", tt.path, got.Value.Worth, tt.worth)
			}
		})
	}
}

func TestGetValueByPath(t *testing.T) {
	root := sample()
	tests := []struct {
		path      string
		wantOK    bool
		wantType  ValueType
		wantWorth string
	}{
		{"root.name", true, TypeString, "top"},
		{"root.tags[0]", true, TypeString, "a"},
		{"root.tags[1]", true, TypeString, "b"},
		{"root.matrix[0][1]", true, TypeNumber, "2"},
		{"root.matrix[1][0]", true, TypeNumber, "3"},
		{"root.group.members[0].scores[1]", true, TypeNumber, "8"},
		{"root.tags", true, TypeArray, ""},
		{"root.tags[5]", false, TypeNull, ""},
		{"root.matrix[0][9]", false, TypeNull, ""},
		{"root.nope", false, TypeNull, ""},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got, ok := root.GetValueByPath(tt.path)
			if ok != tt.wantOK {
				t.Fatalf("GetValueByPath(%q) ok = %v, want %v", tt.path, ok, tt.wantOK)
			}
			if !ok {
				return
			}
			if got.Type != tt.wantType {
				t.Errorf("type = %v, want %v", got.Type, tt.wantType)
			}
			if got.Worth != tt.wantWorth {
				t.Errorf("value = %q, want %q", got.Worth, tt.wantWorth)
			}
		})
	}
}

func TestGetNodeByPath_NilNode(t *testing.T) {
	var nilNode *Node
	if got := nilNode.GetNodeByPath("root"); got != nil {
		t.Errorf("GetNodeByPath on nil = %v, want nil", got)
	}
	if _, ok := nilNode.GetValueByPath("root"); ok {
		t.Error("GetValueByPath on nil reported a hit")
	}
}

// ---- Exists ----

// TestExists_ReachesIntoArrays covers the keys the earlier version reported as
// absent: it only walked object children, so anything inside a list was
// invisible to it.
func TestExists_ReachesIntoArrays(t *testing.T) {
	root := sample()
	tests := []struct {
		key  string
		want bool
	}{
		{"root", true},
		{"name", true},
		{"group", true},
		{"members", true},
		{"scores", true}, // inside an array
		{"detail", true}, // inside an array of arrays
		{"nope", false},
		{"item0", false}, // an array element's own key belongs to the list
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if got := root.Exists(tt.key); got != tt.want {
				t.Errorf("Exists(%q) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}

	var nilNode *Node
	if nilNode.Exists("x") {
		t.Error("Exists on nil reported a hit")
	}
}

// ---- Match rendering ----

func TestMatch_String(t *testing.T) {
	matches := sample().Search("name")
	if got, want := matches[0].String(), "root.name (depth 1) = top"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}

	lists := sample().SearchFunc(func(n *Node) bool { return n.Key == "tags" })
	if got, want := lists[0].String(), "root.tags (depth 1) [array]"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}

	empty := Match{Path: "somewhere"}
	if got := empty.String(); got != "somewhere" {
		t.Errorf("String() on a match with no node = %q", got)
	}
}

// ---- Guards ----

// TestSearch_CyclicTreeTerminates covers a tree that only a hand-built one can
// be. The decoders cannot produce a cycle; a search should still fail rather
// than run forever.
func TestSearch_CyclicTreeTerminates(t *testing.T) {
	a := &Node{Key: "a", Value: &Value{Type: TypeObject}}
	b := &Node{Key: "b", Value: &Value{Type: TypeObject, Node: a}}
	a.Value.Node = b

	done := make(chan int, 1)
	go func() { done <- len(a.Search("b")) }()

	select {
	case n := <-done:
		if n == 0 {
			t.Error("expected the cycle to be walked until the depth limit, got no matches")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Search did not terminate on a cyclic tree")
	}

	if !a.Exists("b") {
		t.Error("Exists did not terminate usefully on a cyclic tree")
	}
}

// ---- The older helpers still answer ----

func TestGetNode_IncludesTheNodeItself(t *testing.T) {
	root := sample()
	if got := root.GetNode("root"); len(got) != 1 || got[0] != root {
		t.Errorf("GetNode(root) = %v, want the node itself", got)
	}
	if got := root.GetNode("name"); len(got) != 6 {
		t.Errorf("GetNode(name) found %d nodes, want 6", len(got))
	}
	var nilNode *Node
	if got := nilNode.GetNode("x"); got != nil {
		t.Errorf("GetNode on nil = %v, want nil", got)
	}
}

func TestFindNodes_IncludesTheNodeItself(t *testing.T) {
	root := sample()
	got := root.FindNodes(func(n *Node) bool { return n.Type() == TypeObject })
	if len(got) == 0 || got[0] != root {
		t.Errorf("FindNodes did not start with the node itself: %v", got)
	}
	var nilNode *Node
	if got := nilNode.FindNodes(func(*Node) bool { return true }); got != nil {
		t.Errorf("FindNodes on nil = %v, want nil", got)
	}
	if got := root.FindNodes(nil); got != nil {
		t.Errorf("FindNodes(nil) = %v, want nil", got)
	}
}
