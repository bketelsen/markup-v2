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
	<lib.FooComponent Bar="42">
</div>
	`

	b := bytes.NewBufferString(h)
	d := NewTagDecoder(b)
	root := Tag{}
	if err := d.Decode(&root); err != nil {
		t.Fatal(err)
	}

	testParseCheckRoot(t, root)
	testParseCheckH1(t, root.Children[0])
	testParseCheckBr(t, root.Children[1])
	testParseCheckInput(t, root.Children[2])
	testParseCheckFooComponent(t, root.Children[3])
}

func testParseCheckRoot(t *testing.T, tag Tag) {
	if name := tag.Name; name != "div" {
		t.Fatalf(`tag name should be "div": "%s"`, name)
	}
	if count := len(tag.Children); count != 4 {
		t.Fatal("tag should have 4 children:", count)
	}
}

func testParseCheckH1(t *testing.T, tag Tag) {
	if name := tag.Name; name != "h1" {
		t.Fatalf(`tag name should be "h1": "%s"`, name)
	}
	if count := len(tag.Children); count != 1 {
		t.Fatal("tag should have 1 children:", count)
	}
	if text := tag.Children[0]; text.Text != "hello" {
		t.Fatalf(`text.Text should be "hello": "%s"`, text.Text)
	}
}

func testParseCheckBr(t *testing.T, tag Tag) {
	if name := tag.Name; name != "br" {
		t.Fatalf(`tag name should be "br": "%s"`, name)
	}
	if count := len(tag.Children); count != 0 {
		t.Fatal("root should not have children:", count)
	}
}

func testParseCheckInput(t *testing.T, tag Tag) {
	if name := tag.Name; name != "input" {
		t.Fatalf(`tag name should be "input": "%s"`, name)
	}
	if count := len(tag.Children); count != 0 {
		t.Fatal("tag should not have children:", count)
	}
	if count := len(tag.Attrs); count != 2 {
		t.Fatal("tag should have 2 attributes:", count)
	}
	if val, _ := tag.Attrs["type"]; val != "text" {
		t.Fatalf(`tag should have an attr with value = "text": %s`, val)
	}
	if _, ok := tag.Attrs["required"]; !ok {
		t.Fatal(`tag should have an attr with key = "required"`)
	}
}

func testParseCheckFooComponent(t *testing.T, tag Tag) {
	if name := tag.Name; name != "lib.foocomponent" {
		t.Fatalf(`tag name should be "lib.foocomponent": "%s"`, name)
	}
	if count := len(tag.Children); count != 0 {
		t.Fatal("tag should not have children:", count)
	}
	if count := len(tag.Attrs); count != 1 {
		t.Fatal("tag should have 1 attribure:", count)
	}
	if val, _ := tag.Attrs["bar"]; val != "42" {
		t.Fatalf(`tag should have an attr with value = "42": %s`, val)
	}
}

func TestParseSelfClosingTagError(t *testing.T) {
	h := `
<p>
	<input />
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
