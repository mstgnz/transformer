package node

import (
	"reflect"
	"testing"
)

var node *Node

// Each method is tested on a single node. Batch-based testing
func TestNodeSingle_Exists(t *testing.T) {
	if got := node.Exists(); !reflect.DeepEqual(got, false) {
		t.Errorf("exists expected=%v, got=%v", false, got)
	}
}

func TestNodeSingle_AddToNext(t *testing.T) {
	node = node.AddToNext(node, nil, "first", "first value")
	node = node.AddToNext(node, nil, "second", "second value")
	if got := node.Key; !reflect.DeepEqual(got, "second") {
		t.Errorf("key expected=%v, got=%v", "second", got)
	}
	if got := node.Value; !reflect.DeepEqual(got, "second value") {
		t.Errorf("value expected=%v, got=%v", "second value", got)
	}
	if got := node.Next; got != nil {
		t.Errorf("next expected=%v, got=%v", nil, got)
	}
	if got := node.Prev.Key; !reflect.DeepEqual(got, "first") {
		t.Errorf("prev expected=%v, got=%v", "first", got)
	}
	if got := node.Prev.Value; !reflect.DeepEqual(got, "first value") {
		t.Errorf("prev expected=%v, got=%v", "first value", got)
	}
	if got := node.Parent; got != nil {
		t.Errorf("parent expected=%v, got=%v", nil, got)
	}
}

func TestNodeSingle_AddToNextWithAttr(t *testing.T) {
	attr := map[string]string{"attr": "attr value"}
	node = node.AddToNextWithAttr(node, nil, "attr", "attr value", attr)
	if got := node.Key; !reflect.DeepEqual(got, "attr") {
		t.Errorf("key expected=%v, got=%v", "attr", got)
	}
	if got := node.Value; !reflect.DeepEqual(got, "attr value") {
		t.Errorf("value expected=%v, got=%v", "attr value", got)
	}
	if got := node.Next; got != nil {
		t.Errorf("next expected=%v, got=%v", nil, got)
	}
	if got := node.Prev.Key; !reflect.DeepEqual(got, "second") {
		t.Errorf("prev expected=%v, got=%v", "second", got)
	}
	if got := node.Parent; got != nil {
		t.Errorf("parent expected=%v, got=%v", nil, got)
	}
	if got := node.Attr; !reflect.DeepEqual(got, attr) {
		t.Errorf("attr expected=%v, got=%v", attr, got)
	}
}

func TestNodeSingle_AddToValue(t *testing.T) {
	node = node.AddToValue(node, nil, "add", "to value")
	node = node.AddToValue(node, nil, "add 2", "to value 2")
	if got := node.Key; !reflect.DeepEqual(got, "add 2") {
		t.Errorf("prev key expected=%v, got=%v", "add 2", got)
	}
	if got := node.Value; !reflect.DeepEqual(got, "to value 2") {
		t.Errorf("value expected=%v, got=%v", "to value 2", got)
	}
	if got := node.Next; got != nil {
		t.Errorf("next ext expected=%v, got=%v", nil, got)
	}
	if got := node.Prev.Key; !reflect.DeepEqual(got, "add") {
		t.Errorf("prev key expected=%v, got=%v", "add", got)
	}
	if got := node.Parent; got != nil {
		t.Errorf("parent expected=%v, got=%v", nil, got)
	}
}

func TestNodeSingle_AddToValueWithAttr(t *testing.T) {
	attr := map[string]string{"attr": "attr value"}
	node = node.AddToValueWithAttr(node, nil, "attr key", "attr value", attr)
	if got := node.Key; !reflect.DeepEqual(got, "attr key") {
		t.Errorf("prev key expected=%v, got=%v", "attr key", got)
	}
	if got := node.Value; !reflect.DeepEqual(got, "attr value") {
		t.Errorf("value expected=%v, got=%v", "attr value", got)
	}
	if got := node.Next; got != nil {
		t.Errorf("next ext expected=%v, got=%v", nil, got)
	}
	got := node.Prev.Value
	if obj, ok := got.(*Node); ok {
		if !reflect.DeepEqual(obj.Value, "attr value") {
			t.Errorf("prev key expected=%v, got=%v", "attr value", obj.Value)
		}
	} else {
		t.Errorf("prev value not a Node")
	}
	if got := node.Parent; got != nil {
		t.Errorf("parent expected=%v, got=%v", nil, got)
	}
}

func TestNodeSingle_SetToValue(t *testing.T) {
	node = node.SetToValue(node, "set key", "set value")
	if got := node.Key; !reflect.DeepEqual(got, "set key") {
		t.Errorf("prev key expected=%v, got=%v", "set key", got)
	}
	if got := node.Value; !reflect.DeepEqual(got, "set value") {
		t.Errorf("value expected=%v, got=%v", "set value", got)
	}
	if got := node.Next; got != nil {
		t.Errorf("next ext expected=%v, got=%v", nil, got)
	}
	if got := node.Prev.Key; !reflect.DeepEqual(got, "add 2") {
		t.Errorf("prev key expected=%v, got=%v", "add 2", got)
	}
	if got := node.Parent; got != nil {
		t.Errorf("parent expected=%v, got=%v", nil, got)
	}
}

func TestNodeSingle_AddToArr(t *testing.T) {
	node = node.AddToNext(node, node, "slice", []any{})
	node = node.AddToArr(node, "arr element")
	if got := node.Key; !reflect.DeepEqual(got, "slice") {
		t.Errorf("prev key expected=%v, got=%v", "slice", got)
	}
	if got := node.Value; !reflect.DeepEqual(got, []any{"arr element"}) {
		t.Errorf("value expected=%v, got=%v", []any{"arr element"}, got)
	}
	if got := node.Next; got != nil {
		t.Errorf("next ext expected=%v, got=%v", nil, got)
	}
	if got := node.Prev.Key; !reflect.DeepEqual(got, "set key") {
		t.Errorf("prev key expected=%v, got=%v", "set key", got)
	}
	if got := node.Parent.Value; !reflect.DeepEqual(got, "set value") {
		t.Errorf("parent expected=%v, got=%v", "set value", got)
	}
}

func TestNodeSingle_AddObjToArr(t *testing.T) {
	node = node.AddObjToArr(node)
	if got := node.Key; !reflect.DeepEqual(got, "") {
		t.Errorf("prev key expected=%v, got=%v", "", got)
	}
	if got := node.Value; !reflect.DeepEqual(got, nil) {
		t.Errorf("value expected=%v, got=%v", nil, got)
	}
	if got := node.Next; got != nil {
		t.Errorf("next ext expected=%v, got=%v", nil, got)
	}
	if got := node.Prev.Key; !reflect.DeepEqual(got, "slice") {
		t.Errorf("prev key expected=%v, got=%v", "slice", got)
	}
	if slc, ok := node.Parent.Value.([]any); ok {
		if !reflect.DeepEqual(slc[0], "arr element") {
			t.Errorf("parent expected=%v, got=%v", "arr element", slc)
		}
	} else {
		t.Errorf("prev value not a slice")
	}
}

func TestNodeSingle_GetNode(t *testing.T) {

}

// Each method is tested on a new node.
func TestNode_Exists(t *testing.T) {
	var knot *Node
	if got := knot.Exists(); !reflect.DeepEqual(got, false) {
		t.Errorf("exists expected=%v, got=%v", false, got)
	}
}

func TestNode_AddToNext(t *testing.T) {
	var knot *Node
	knot = knot.AddToNext(knot, nil, "first", "first value")
	expected := "first"
	if got := knot.Key; !reflect.DeepEqual(got, expected) {
		t.Errorf("key expected=%v, got=%v", expected, got)
	}
	expected = "first value"
	if got := knot.Value; !reflect.DeepEqual(got, expected) {
		t.Errorf("value expected=%v, got=%v", expected, got)
	}
	if got := knot.Next; got != nil {
		t.Errorf("next expected=%v, got=%v", nil, got)
	}
	if got := knot.Prev; got != nil {
		t.Errorf("prev expected=%v, got=%v", nil, got)
	}
	if got := knot.Parent; got != nil {
		t.Errorf("parent expected=%v, got=%v", nil, got)
	}
}

func TestNode_AddToNextWithAttr(t *testing.T) {
	var knot *Node
	attr := map[string]string{"attr": "attr value"}
	knot = knot.AddToNextWithAttr(knot, nil, "second", "second value", attr)
	if got := knot.Key; !reflect.DeepEqual(got, "second") {
		t.Errorf("key expected=%v, got=%v", "second", got)
	}
	if got := knot.Value; !reflect.DeepEqual(got, "second value") {
		t.Errorf("value expected=%v, got=%v", "second value", got)
	}
	if got := knot.Next; got != nil {
		t.Errorf("next expected=%v, got=%v", nil, got)
	}
	if got := knot.Prev; got != nil {
		t.Errorf("prev expected=%v, got=%v", nil, got)
	}
	if got := knot.Parent; got != nil {
		t.Errorf("parent expected=%v, got=%v", nil, got)
	}
	if got := knot.Attr; !reflect.DeepEqual(got, attr) {
		t.Errorf("attr expected=%v, got=%v", attr, got)
	}
}

func TestNode_AddToValue(t *testing.T) {
	var knot *Node
	knot = knot.AddToValue(knot, nil, "add", "to value")
	knot = knot.AddToValue(knot, nil, "add 2", "to value 2")
	if got := knot.Key; !reflect.DeepEqual(got, "add 2") {
		t.Errorf("prev key expected=%v, got=%v", "add 2", got)
	}
	if got := knot.Value; !reflect.DeepEqual(got, "to value 2") {
		t.Errorf("value expected=%v, got=%v", "to value 2", got)
	}
	if got := knot.Next; got != nil {
		t.Errorf("next ext expected=%v, got=%v", nil, got)
	}
	if got := knot.Prev.Key; !reflect.DeepEqual(got, "add") {
		t.Errorf("prev key expected=%v, got=%v", "add", got)
	}
	if got := knot.Parent; got != nil {
		t.Errorf("parent expected=%v, got=%v", nil, got)
	}
}

func TestNode_AddToValueWithAttr(t *testing.T) {
	var knot *Node
	attr := map[string]string{"attr": "attr value"}
	knot = knot.AddToValueWithAttr(knot, nil, "add", "to value", attr)
	if got := knot.Key; !reflect.DeepEqual(got, "add") {
		t.Errorf("prev key expected=%v, got=%v", "add", got)
	}
	if got := knot.Value; !reflect.DeepEqual(got, "to value") {
		t.Errorf("value expected=%v, got=%v", "to value", got)
	}
	if got := knot.Next; got != nil {
		t.Errorf("next ext expected=%v, got=%v", nil, got)
	}
	if got := knot.Prev; got != nil {
		t.Errorf("prev key expected=%v, got=%v", nil, got)
	}
	if got := knot.Parent; got != nil {
		t.Errorf("parent expected=%v, got=%v", nil, got)
	}
}

func TestNode_SetToValue(t *testing.T) {
	var knot *Node
	knot = knot.SetToValue(knot, "set key", "set value")
	if got := knot.Key; !reflect.DeepEqual(got, "set key") {
		t.Errorf("prev key expected=%v, got=%v", "set key", got)
	}
	if got := knot.Value; !reflect.DeepEqual(got, "set value") {
		t.Errorf("value expected=%v, got=%v", "set value", got)
	}
	if got := knot.Next; got != nil {
		t.Errorf("next ext expected=%v, got=%v", nil, got)
	}
	if got := knot.Prev; got != nil {
		t.Errorf("prev key expected=%v, got=%v", nil, got)
	}
	if got := knot.Parent; got != nil {
		t.Errorf("parent expected=%v, got=%v", nil, got)
	}
}

func TestNode_AddToArr(t *testing.T) {
	var knot *Node
	knot = knot.AddToArr(knot, "arr element")
	if got := knot.Key; !reflect.DeepEqual(got, "") {
		t.Errorf("prev key expected=%v, got=%v", "", got)
	}
	if got := knot.Value; !reflect.DeepEqual(got, []any{"arr element"}) {
		t.Errorf("value expected=%v, got=%v", []any{"arr element"}, got)
	}
	if got := knot.Next; got != nil {
		t.Errorf("next ext expected=%v, got=%v", nil, got)
	}
	if got := knot.Prev; got != nil {
		t.Errorf("prev key expected=%v, got=%v", nil, got)
	}
	if got := knot.Parent; got != nil {
		t.Errorf("parent expected=%v, got=%v", nil, got)
	}
}

func TestNode_AddObjToArr(t *testing.T) {
	var knot *Node
	knot = knot.AddObjToArr(knot)
	if got := knot.Key; !reflect.DeepEqual(got, "") {
		t.Errorf("prev key expected=%v, got=%v", "", got)
	}
	if got := knot.Value; !reflect.DeepEqual(got, nil) {
		t.Errorf("value expected=%v, got=%v", nil, got)
	}
	if got := knot.Next; got != nil {
		t.Errorf("next ext expected=%v, got=%v", nil, got)
	}
	if got := knot.Prev; got != nil {
		t.Errorf("prev key expected=%v, got=%v", nil, got)
	}
	if got := knot.Parent; got != nil {
		t.Errorf("parent expected=%v, got=%v", nil, got)
	}
}
