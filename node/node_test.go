package node

import (
	"reflect"
	"testing"
)

var node *Node

func TestNode_Exists(t *testing.T) {
	if got := node.Exists(); !reflect.DeepEqual(got, false) {
		t.Errorf("exists expected=%v, got=%v", false, got)
	}
}

func TestNode_AddToStart(t *testing.T) {
	node = node.AddToStart(nil)
	if got := node; got != nil {
		t.Errorf("key expected=%v, got=%v", nil, got)
	}
	node = &Node{Key: "first"}
	node = node.AddToStart(&Node{Key: "new first"})
	if got := node.Key; !reflect.DeepEqual(got, "new first") {
		t.Errorf("value expected=%v, got=%v", "new first", got)
	}
	if got := node.Next.Key; !reflect.DeepEqual(got, "first") {
		t.Errorf("value expected=%v, got=%v", "first", got)
	}
}

func TestNode_AddToNext(t *testing.T) {
	node = node.AddToNext(node, nil, "first")
	node = node.AddToValue(node, Value{Worth: "first value"})
	node = node.AddToNext(node, nil, "second")
	node = node.AddToValue(node, Value{Worth: "second value"})
	if got := node.Key; !reflect.DeepEqual(got, "second") {
		t.Errorf("key expected=%v, got=%v", "second", got)
	}
	if got := node.Value.Worth; !reflect.DeepEqual(got, "second value") {
		t.Errorf("value expected=%v, got=%v", "second value", got)
	}
	if got := node.Next; got != nil {
		t.Errorf("next expected=%v, got=%v", nil, got)
	}
	if got := node.Prev.Key; !reflect.DeepEqual(got, "first") {
		t.Errorf("prev expected=%v, got=%v", "first", got)
	}
	if got := node.Prev.Value.Worth; !reflect.DeepEqual(got, "first value") {
		t.Errorf("prev expected=%v, got=%v", "first value", got)
	}
	if got := node.Parent; got != nil {
		t.Errorf("parent expected=%v, got=%v", nil, got)
	}
}

func TestNode_AddToValue(t *testing.T) {
	// test value worth
	node = node.AddToNext(node, nil, "third")
	node = node.AddToValue(node, Value{Worth: "third value"})
	if got := node.Key; !reflect.DeepEqual(got, "third") {
		t.Errorf("key expected=%v, got=%v", "third", got)
	}
	if got := node.Value.Worth; !reflect.DeepEqual(got, "third value") {
		t.Errorf("key expected=%v, got=%v", "third value", got)
	}
	// test value node
	node = node.AddToValue(node, Value{Node: &Node{Key: "value node"}})
	if got := node.Key; !reflect.DeepEqual(got, "value node") {
		t.Errorf("key expected=%v, got=%v", "value node", got)
	}
	// test value map
	attr := map[string]string{"xml": "value"}
	node = node.AddToValue(node, Value{Attr: attr})
	if got := node.Value.Attr; !reflect.DeepEqual(got, attr) {
		t.Errorf("key expected=%v, got=%v", attr, got)
	}
	// test slice
	attr = map[string]string{"slc": "value"}
	node = node.AddToValue(node, Value{
		Slice: []Value{
			{Worth: "slc value"},
			{Attr: attr},
			{Slice: []Value{{}, {Worth: "middle"}, {}}},
		},
	})
	if got := node.Value.Slice[0].Worth; !reflect.DeepEqual(got, "slc value") {
		t.Errorf("key expected=%v, got=%v", "slc value", got)
	}
	if got := node.Value.Slice[1].Attr; !reflect.DeepEqual(got, attr) {
		t.Errorf("key expected=%v, got=%v", attr, got)
	}
	if got := node.Value.Slice[2].Slice[1].Worth; !reflect.DeepEqual(got, "middle") {
		t.Errorf("key expected=%v, got=%v", "middle", got)
	}
	// test slice with node
	node = node.AddToValue(node, Value{
		Slice: []Value{
			{Node: &Node{Key: "slc node key"}},
		},
	})
	if got := node.Key; !reflect.DeepEqual(got, "slc node key") {
		t.Errorf("key expected=%v, got=%v", "slc node key", got)
	}
}

func TestNode_AddToEnd(t *testing.T) {
	node = node.AddToEnd(nil)
	if got := node; got != nil {
		t.Errorf("key expected=%v, got=%v", nil, got)
	}
	node = node.AddToNext(node, nil, "end first")
	second := node.AddToNext(node, nil, "end second")
	node = node.AddToNext(second, nil, "end third")
	node = second.AddToEnd(&Node{Key: "end to end"})
	if got := node.Key; !reflect.DeepEqual(got, "end to end") {
		t.Errorf("key expected=%v, got=%v", "end to end", got)
	}
}

func TestNode_Delete(t *testing.T) {
	node = node.AddToNext(node, nil, "del first")
	second := node.AddToNext(node, nil, "del second")
	node = node.AddToNext(node, nil, "del third")
	node = node.Delete(second)
	if got := node.Key; !reflect.DeepEqual(got, "del third") {
		t.Errorf("key expected=%v, got=%v", "del third", got)
	}
	if got := node.Prev.Key; !reflect.DeepEqual(got, "del first") {
		t.Errorf("key expected=%v, got=%v", "del first", got)
	}

}

func TestNode_GetNode(t *testing.T) {
	node = node.AddToNext(node, nil, "first")
	node = node.AddToValue(node, Value{Worth: "first value"})
	node = node.AddToNext(node, nil, "second")
	node = node.AddToValue(node, Value{Worth: "second value"})
	node = node.AddToNext(node, nil, "third")
	node = node.AddToValue(node, Value{Worth: "third value"})
	nodes := node.GetNode("second")
	if got := nodes[0].Key; !reflect.DeepEqual(got, "second") {
		t.Errorf("key expected=%v, got=%v", "second", got)
	}
}

func TestNode_Reset(t *testing.T) {

}
