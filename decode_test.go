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
	// h1 := root.Children[0]

}

func TestParseSelfClosingTagError(t *testing.T) {
	h := `
<p>
	<div/>
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
