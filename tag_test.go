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
