package markup

import "testing"

func TestTagIsEmpty(t *testing.T) {
	tag := Tag{}
	if !tag.IsEmpty() {
		t.Error("tag should be empty")
	}
}

func TestTagIsText(t *testing.T) {
	tag := Tag{Text: "foo"}
	if !tag.IsText() {
		t.Error("tag should be a text")
	}

	tag = Tag{Name: "div", Text: "foo"}
	if tag.IsText() {
		t.Error("tag should not be a text")
	}

	tag = Tag{}
	if tag.IsText() {
		t.Error("tag should not be a text")
	}
}

func TestTagIsComponent(t *testing.T) {
	tag := Tag{Name: "foo"}
	if !tag.IsComponent() {
		t.Error("tag should be a component")
	}

	tag = Tag{Name: "div"}
	if tag.IsComponent() {
		t.Error("tag should not be a component")
	}

	tag = Tag{}
	if tag.IsComponent() {
		t.Error("tag should not be a component")
	}
}

func TestTagIsVoidElement(t *testing.T) {
	tag := Tag{Name: "input"}
	if !tag.IsVoidElem() {
		t.Error("tag should be a void element")
	}

	tag = Tag{Name: "div"}
	if tag.IsVoidElem() {
		t.Error("tag should not be a void element")
	}
}

func TestAttrEquals(t *testing.T) {
	attr := AttrMap{
		"hello": "world",
		"foo":   "bar",
	}

	attr2 := AttrMap{
		"foo":   "bar",
		"hello": "world",
	}

	if !AttrEquals(attr, attr2) {
		t.Error("attr and attr2 should be equals")
	}

	if AttrEquals(attr, nil) {
		t.Error("attr and nil should not be equals")
	}

	attr3 := AttrMap{
		"foo":   "bar",
		"hello": "maxoo",
	}

	if AttrEquals(attr, attr3) {
		t.Error("attr and attr3 should not be equals")
	}

	attr4 := AttrMap{
		"foo": "bar",
		"bye": "world",
	}

	if AttrEquals(attr, attr4) {
		t.Error("attr and attr4 should not be equals")
	}
}
