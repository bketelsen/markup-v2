package markup

import (
	"github.com/google/uuid"
	"golang.org/x/net/html/atom"
)

// Tag represents an HTML tag.
type Tag struct {
	ID       uuid.UUID
	Name     string
	Text     string
	Attrs    []Attr
	Children []Tag
}

// Attr represents a tag attribute.
type Attr struct {
	Key string
	Val string
}

// IsEmpty reports whether its argument t is nil.
// Empty tags have empty name and empty text
func (t *Tag) IsEmpty() bool {
	return len(t.Name) == 0 && len(t.Text) == 0
}

// IsText reports whether its argument t represents a text.
// Texts doesn't have a name and are not empty.
func (t *Tag) IsText() bool {
	if t.IsEmpty() {
		return false
	}
	return len(t.Name) == 0 && len(t.Text) != 0
}

// IsComponent reports whether its argument t represents a component.
// Components are not void elements and not empty.
func (t *Tag) IsComponent() bool {
	if t.IsEmpty() {
		return false
	}

	a := atom.Lookup([]byte(t.Name))
	return a == 0
}

// IsVoidElem reports whether its argument t represents a void element.
// Void elements are tags listed at
// https://www.w3.org/TR/html5/syntax.html#void-elements.
func (t *Tag) IsVoidElem() bool {
	_, ok := voidElems[t.Name]
	return ok
}

var (
	voidElems = map[string]bool{
		"area":   true,
		"base":   true,
		"br":     true,
		"col":    true,
		"embed":  true,
		"hr":     true,
		"img":    true,
		"input":  true,
		"keygen": true,
		"link":   true,
		"meta":   true,
		"param":  true,
		"source": true,
		"track":  true,
		"wbr":    true,
	}
)
