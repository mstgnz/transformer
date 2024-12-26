package node

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestValueType_String(t *testing.T) {
	tests := []struct {
		name string
		t    ValueType
		want string
	}{
		{
			name: "TypeNull",
			t:    TypeNull,
			want: "null",
		},
		{
			name: "TypeObject",
			t:    TypeObject,
			want: "object",
		},
		{
			name: "TypeArray",
			t:    TypeArray,
			want: "array",
		},
		{
			name: "TypeString",
			t:    TypeString,
			want: "string",
		},
		{
			name: "TypeNumber",
			t:    TypeNumber,
			want: "number",
		},
		{
			name: "TypeBoolean",
			t:    TypeBoolean,
			want: "boolean",
		},
		{
			name: "Unknown type",
			t:    ValueType(99),
			want: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.String(); got != tt.want {
				t.Errorf("ValueType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNode(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want *Node
	}{
		{
			name: "Create node with key",
			key:  "test",
			want: &Node{
				Key: "test",
				Value: &Value{
					Type: TypeNull,
				},
			},
		},
		{
			name: "Create node with empty key",
			key:  "",
			want: &Node{
				Key: "",
				Value: &Value{
					Type: TypeNull,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewNode(tt.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_AddToStart(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() (*Node, *Node)
		wantErr bool
	}{
		{
			name: "Add to empty node",
			setup: func() (*Node, *Node) {
				parent := NewNode("parent")
				child := NewNode("child")
				return parent, child
			},
			wantErr: false,
		},
		{
			name: "Add to node with existing children",
			setup: func() (*Node, *Node) {
				parent := NewNode("parent")
				existing := NewNode("existing")
				parent.AddToStart(existing)
				child := NewNode("child")
				return parent, child
			},
			wantErr: false,
		},
		{
			name: "Add to nil node",
			setup: func() (*Node, *Node) {
				return nil, NewNode("child")
			},
			wantErr: true,
		},
		{
			name: "Add nil node",
			setup: func() (*Node, *Node) {
				return NewNode("parent"), nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parent, child := tt.setup()
			err := parent.AddToStart(child)
			if (err != nil) != tt.wantErr {
				t.Errorf("Node.AddToStart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if parent.Value.Node != child {
					t.Error("Node.AddToStart() did not add child to start")
				}
				if child.Parent != parent {
					t.Error("Node.AddToStart() did not set parent reference")
				}
			}
		})
	}
}

func TestNode_AddToEnd(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() (*Node, *Node)
		wantErr bool
	}{
		{
			name: "Add to empty node",
			setup: func() (*Node, *Node) {
				parent := NewNode("parent")
				child := NewNode("child")
				return parent, child
			},
			wantErr: false,
		},
		{
			name: "Add to node with existing children",
			setup: func() (*Node, *Node) {
				parent := NewNode("parent")
				existing := NewNode("existing")
				parent.AddToEnd(existing)
				child := NewNode("child")
				return parent, child
			},
			wantErr: false,
		},
		{
			name: "Add to nil node",
			setup: func() (*Node, *Node) {
				return nil, NewNode("child")
			},
			wantErr: true,
		},
		{
			name: "Add nil node",
			setup: func() (*Node, *Node) {
				return NewNode("parent"), nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parent, child := tt.setup()
			err := parent.AddToEnd(child)
			if (err != nil) != tt.wantErr {
				t.Errorf("Node.AddToEnd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if child.Parent != parent {
					t.Error("Node.AddToEnd() did not set parent reference")
				}
				// Find last child
				last := parent.Value.Node
				for last.Next != nil {
					last = last.Next
				}
				if last != child {
					t.Error("Node.AddToEnd() did not add child to end")
				}
			}
		})
	}
}

func TestNode_AddToValue(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() (*Node, *Value)
		wantErr bool
	}{
		{
			name: "Add string value",
			setup: func() (*Node, *Value) {
				node := NewNode("test")
				value := &Value{
					Type:  TypeString,
					Worth: "test value",
				}
				return node, value
			},
			wantErr: false,
		},
		{
			name: "Add object value",
			setup: func() (*Node, *Value) {
				node := NewNode("test")
				child := NewNode("child")
				value := &Value{
					Type: TypeObject,
					Node: child,
				}
				return node, value
			},
			wantErr: false,
		},
		{
			name: "Add array value",
			setup: func() (*Node, *Value) {
				node := NewNode("test")
				value := &Value{
					Type:  TypeArray,
					Array: []*Value{{Type: TypeString, Worth: "item"}},
				}
				return node, value
			},
			wantErr: false,
		},
		{
			name: "Add to nil node",
			setup: func() (*Node, *Value) {
				return nil, &Value{Type: TypeString}
			},
			wantErr: true,
		},
		{
			name: "Add nil value",
			setup: func() (*Node, *Value) {
				return NewNode("test"), nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, value := tt.setup()
			err := node.AddToValue(value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Node.AddToValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if node.Value != value {
					t.Error("Node.AddToValue() did not set value")
				}
				if value.Node != nil && value.Node.Parent != node {
					t.Error("Node.AddToValue() did not set parent reference for nested node")
				}
			}
		})
	}
}

func TestNode_Delete(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *Node
		wantErr bool
	}{
		{
			name: "Delete middle node",
			setup: func() *Node {
				parent := NewNode("parent")
				prev := NewNode("prev")
				target := NewNode("target")
				next := NewNode("next")
				parent.AddToStart(prev)
				parent.AddToEnd(target)
				parent.AddToEnd(next)
				return target
			},
			wantErr: false,
		},
		{
			name: "Delete first node",
			setup: func() *Node {
				parent := NewNode("parent")
				target := NewNode("target")
				next := NewNode("next")
				parent.AddToStart(target)
				parent.AddToEnd(next)
				return target
			},
			wantErr: false,
		},
		{
			name: "Delete last node",
			setup: func() *Node {
				parent := NewNode("parent")
				prev := NewNode("prev")
				target := NewNode("target")
				parent.AddToStart(prev)
				parent.AddToEnd(target)
				return target
			},
			wantErr: false,
		},
		{
			name: "Delete nil node",
			setup: func() *Node {
				return nil
			},
			wantErr: true,
		},
		{
			name: "Delete root node",
			setup: func() *Node {
				return NewNode("root")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := tt.setup()
			err := node.Delete()
			if (err != nil) != tt.wantErr {
				t.Errorf("Node.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if node.Parent != nil {
					t.Error("Node.Delete() did not clear parent reference")
				}
				if node.Next != nil {
					t.Error("Node.Delete() did not clear next reference")
				}
				if node.Prev != nil {
					t.Error("Node.Delete() did not clear prev reference")
				}
			}
		})
	}
}

func TestNode_GetNode(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *Node
		key      string
		wantLen  int
		wantKeys []string
	}{
		{
			name: "Find single node",
			setup: func() *Node {
				root := NewNode("root")
				target := NewNode("target")
				root.AddToStart(target)
				return root
			},
			key:      "target",
			wantLen:  1,
			wantKeys: []string{"target"},
		},
		{
			name: "Find multiple nodes",
			setup: func() *Node {
				root := NewNode("root")
				target1 := NewNode("target")
				target2 := NewNode("target")
				root.AddToStart(target1)
				root.AddToEnd(target2)
				return root
			},
			key:      "target",
			wantLen:  2,
			wantKeys: []string{"target", "target"},
		},
		{
			name: "Find nested nodes",
			setup: func() *Node {
				root := NewNode("root")
				parent := NewNode("parent")
				target := NewNode("target")
				parent.AddToStart(target)
				root.AddToStart(parent)
				return root
			},
			key:      "target",
			wantLen:  1,
			wantKeys: []string{"target"},
		},
		{
			name: "Find in array",
			setup: func() *Node {
				root := NewNode("root")
				root.Value = &Value{
					Type: TypeArray,
					Array: []*Value{
						{
							Node: NewNode("target"),
						},
					},
				}
				return root
			},
			key:      "target",
			wantLen:  1,
			wantKeys: []string{"target"},
		},
		{
			name: "No matches",
			setup: func() *Node {
				root := NewNode("root")
				other := NewNode("other")
				root.AddToStart(other)
				return root
			},
			key:      "target",
			wantLen:  0,
			wantKeys: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := tt.setup()
			got := root.GetNode(tt.key)
			if len(got) != tt.wantLen {
				t.Errorf("Node.GetNode() returned %d nodes, want %d", len(got), tt.wantLen)
				return
			}
			for i, node := range got {
				if node.Key != tt.wantKeys[i] {
					t.Errorf("Node.GetNode() returned node with key %s, want %s", node.Key, tt.wantKeys[i])
				}
			}
		})
	}
}

func TestNode_GetNodeByPath(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *Node
		path    string
		wantKey string
		wantNil bool
	}{
		{
			name: "Find direct child",
			setup: func() *Node {
				root := NewNode("root")
				target := NewNode("child")
				root.AddToStart(target)
				return root
			},
			path:    "root.child",
			wantKey: "child",
			wantNil: false,
		},
		{
			name: "Find nested child",
			setup: func() *Node {
				root := NewNode("root")
				parent := NewNode("parent")
				child := NewNode("child")
				parent.AddToStart(child)
				root.AddToStart(parent)
				return root
			},
			path:    "root.parent.child",
			wantKey: "child",
			wantNil: false,
		},
		{
			name: "Path not found",
			setup: func() *Node {
				root := NewNode("root")
				other := NewNode("other")
				root.AddToStart(other)
				return root
			},
			path:    "root.nonexistent",
			wantNil: true,
		},
		{
			name: "Empty path",
			setup: func() *Node {
				return NewNode("root")
			},
			path:    "",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := tt.setup()
			got := root.GetNodeByPath(tt.path)
			if tt.wantNil {
				if got != nil {
					t.Error("Node.GetNodeByPath() returned non-nil node, want nil")
				}
				return
			}
			if got == nil {
				t.Error("Node.GetNodeByPath() returned nil node, want non-nil")
				return
			}
			if got.Key != tt.wantKey {
				t.Errorf("Node.GetNodeByPath() returned node with key %s, want %s", got.Key, tt.wantKey)
			}
		})
	}
}

func TestNode_FindNodes(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *Node
		predicate func(*Node) bool
		wantLen   int
		wantKeys  []string
	}{
		{
			name: "Find by value type",
			setup: func() *Node {
				root := NewNode("root")
				node1 := NewNode("node1")
				node1.Value.Type = TypeString
				node2 := NewNode("node2")
				node2.Value.Type = TypeNumber
				root.AddToStart(node1)
				root.AddToEnd(node2)
				return root
			},
			predicate: func(n *Node) bool {
				return n.Value != nil && n.Value.Type == TypeString
			},
			wantLen:  1,
			wantKeys: []string{"node1"},
		},
		{
			name: "Find by key prefix",
			setup: func() *Node {
				root := NewNode("root")
				test1 := NewNode("test_1")
				test2 := NewNode("test_2")
				other := NewNode("other")
				root.AddToStart(test1)
				root.AddToEnd(test2)
				root.AddToEnd(other)
				return root
			},
			predicate: func(n *Node) bool {
				return strings.HasPrefix(n.Key, "test_")
			},
			wantLen:  2,
			wantKeys: []string{"test_1", "test_2"},
		},
		{
			name: "No matches",
			setup: func() *Node {
				root := NewNode("root")
				node := NewNode("node")
				root.AddToStart(node)
				return root
			},
			predicate: func(n *Node) bool {
				return false
			},
			wantLen:  0,
			wantKeys: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := tt.setup()
			got := root.FindNodes(tt.predicate)
			if len(got) != tt.wantLen {
				t.Errorf("Node.FindNodes() returned %d nodes, want %d", len(got), tt.wantLen)
				return
			}
			for i, node := range got {
				if node.Key != tt.wantKeys[i] {
					t.Errorf("Node.FindNodes() returned node with key %s, want %s", node.Key, tt.wantKeys[i])
				}
			}
		})
	}
}

func TestNode_Validate(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *Node
		wantErr bool
	}{
		{
			name: "Valid simple node",
			setup: func() *Node {
				return NewNode("test")
			},
			wantErr: false,
		},
		{
			name: "Valid object node",
			setup: func() *Node {
				root := NewNode("root")
				child := NewNode("child")
				root.AddToStart(child)
				return root
			},
			wantErr: false,
		},
		{
			name: "Valid array node",
			setup: func() *Node {
				node := NewNode("array")
				node.Value = &Value{
					Type:  TypeArray,
					Array: []*Value{{Type: TypeString, Worth: "item"}},
				}
				return node
			},
			wantErr: false,
		},
		{
			name: "Nil node",
			setup: func() *Node {
				return nil
			},
			wantErr: true,
		},
		{
			name: "Node with nil value",
			setup: func() *Node {
				node := NewNode("test")
				node.Value = nil
				return node
			},
			wantErr: true,
		},
		{
			name: "Invalid parent-child relationship",
			setup: func() *Node {
				root := NewNode("root")
				child := NewNode("child")
				child.Parent = root // Set parent without proper AddToStart/AddToEnd
				return child
			},
			wantErr: true,
		},
		{
			name: "Invalid next-prev relationship",
			setup: func() *Node {
				root := NewNode("root")
				node1 := NewNode("node1")
				node2 := NewNode("node2")
				root.AddToStart(node1)
				root.AddToEnd(node2)
				node1.Next = node2
				node2.Prev = nil // Break the prev reference
				return root
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := tt.setup()
			err := node.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Node.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNode_Print(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *Node
		want  []string
	}{
		{
			name: "Print simple node",
			setup: func() *Node {
				node := NewNode("test")
				node.Value.Type = TypeString
				node.Value.Worth = "value"
				return node
			},
			want: []string{
				"test: \"value\"",
			},
		},
		{
			name: "Print object node",
			setup: func() *Node {
				root := NewNode("root")
				child1 := NewNode("child1")
				child1.Value.Type = TypeString
				child1.Value.Worth = "value1"
				child2 := NewNode("child2")
				child2.Value.Type = TypeNumber
				child2.Value.Worth = "42"
				root.AddToStart(child1)
				root.AddToEnd(child2)
				return root
			},
			want: []string{
				"root: {",
				"  child1: \"value1\"",
				"  child2: 42",
				"}",
			},
		},
		{
			name: "Print array node",
			setup: func() *Node {
				root := NewNode("array")
				root.Value = &Value{
					Type: TypeArray,
					Array: []*Value{
						{Type: TypeString, Worth: "item1"},
						{Type: TypeNumber, Worth: "42"},
						{Type: TypeBoolean, Worth: "true"},
					},
				}
				return root
			},
			want: []string{
				"array: [",
				"  \"item1\"",
				"  ,",
				"  42",
				"  ,",
				"  true",
				"]",
			},
		},
		{
			name: "Print complex node",
			setup: func() *Node {
				root := NewNode("root")
				obj := NewNode("object")
				arr := NewNode("array")
				arr.Value = &Value{
					Type: TypeArray,
					Array: []*Value{
						{Type: TypeString, Worth: "item"},
						{Type: TypeNull},
					},
				}
				str := NewNode("string")
				str.Value.Type = TypeString
				str.Value.Worth = "value"
				obj.AddToStart(arr)
				obj.AddToEnd(str)
				root.AddToStart(obj)
				return root
			},
			want: []string{
				"root: {",
				"  object: {",
				"    array: [",
				"      \"item\"",
				"      ,",
				"      null",
				"    ]",
				"    string: \"value\"",
				"  }",
				"}",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Print node
			node := tt.setup()
			node.Print()

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Read captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			got := strings.Split(strings.TrimSpace(buf.String()), "\n")

			// Compare output
			if len(got) != len(tt.want) {
				t.Errorf("Print() output has %d lines, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("Print() line %d = %q, want %q", i+1, got[i], tt.want[i])
				}
			}
		})
	}
}
