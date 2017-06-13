package markup

import (
	"bytes"
	"testing"
)

func TestParse(t *testing.T) {
	h := `
<div>
	<h1>hello</h1>
	<br>
	<input type="text" required>
	<FooComponent Bar="42">
</div>
	`

	b := bytes.NewBufferString(h)
	d := NewTagDecoder(b)
	root := Tag{}
	if err := d.Decode(&root); err != nil {
		t.Fatal(err)
	}

	// Check root.
	if name := root.Name; name != "div" {
		t.Fatalf(`root name should be "div": "%s"`, name)
	}
	if count := len(root.Children); count != 4 {
		t.Fatal("root should have 4 children:", count)
	}

	// Check h1.
	h1 := root.Children[0]
	if name := h1.Name; name != "h1" {
		t.Fatalf(`h1 name should be "h1": "%s"`, name)
	}
	if count := len(h1.Children); count != 1 {
		t.Fatal("h1 should have 1 children:", count)
	}
	if text := h1.Children[0]; text.Text != "hello" {
		t.Fatalf(`text.Text should be "hello": "%s"`, text.Text)
	}

	// Check br.
	br := root.Children[1]
	if name := br.Name; name != "br" {
		t.Fatalf(`br name should be "br": "%s"`, name)
	}
	if count := len(br.Children); count != 0 {
		t.Fatal("root should not have children:", count)
	}

	// Check input.
	input := root.Children[2]
	if name := input.Name; name != "input" {
		t.Fatalf(`input name should be "input": "%s"`, name)
	}
	if count := len(input.Children); count != 0 {
		t.Fatal("input should not have children:", count)
	}
	if count := len(input.Attrs); count != 2 {
		t.Fatal("input should have 2 attributes:", count)
	}
	if attr, expec := input.Attrs[0], (Attr{"type", "text"}); attr != expec {
		t.Fatalf("attr != expec: %+v != %+v", attr, expec)
	}
	if attr, expec := input.Attrs[1], (Attr{Key: "required"}); attr != expec {
		t.Fatalf("attr != expec: %+v != %+v", attr, expec)
	}

	// Check input.
	foo := root.Children[3]
	if name := foo.Name; name != "foocomponent" {
		t.Fatalf(`foo name should be "foocomponent": "%s"`, name)
	}
	if count := len(foo.Children); count != 0 {
		t.Fatal("foo should not have children:", count)
	}
	if count := len(foo.Attrs); count != 1 {
		t.Fatal("foo should have 1 attribure:", count)
	}
	if attr, expec := foo.Attrs[0], (Attr{"bar", "42"}); attr != expec {
		t.Fatalf("attr != expec: %+v != %+v", attr, expec)
	}
}

func TestParseSelfClosingTagError(t *testing.T) {
	h := `
<p>
	<input/>
</p>
`

	b := bytes.NewBufferString(h)
	d := NewTagDecoder(b)
	root := Tag{}
	if err := d.Decode(&root); err == nil {
		t.Fatal("err should not be nil")
	}
}

func TestParseEmptyHTML(t *testing.T) {
	h := ""

	b := bytes.NewBufferString(h)
	d := NewTagDecoder(b)

	root := Tag{}
	if err := d.Decode(&root); err == nil {
		t.Fatal("err should not be nil")
	}
}

func TestParseNonClosingHTML(t *testing.T) {
	h := "<body><div>"

	b := bytes.NewBufferString(h)
	d := NewTagDecoder(b)

	root := Tag{}
	if err := d.Decode(&root); err != nil {
		t.Fatal(err)
	}
}
