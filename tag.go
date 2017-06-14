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
	Attrs    map[string]string
	Children []Tag
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
	void      = struct{}{}
	voidElems = map[string]struct{}{
		"area":   void,
		"base":   void,
		"br":     void,
		"col":    void,
		"embed":  void,
		"hr":     void,
		"img":    void,
		"input":  void,
		"keygen": void,
		"link":   void,
		"meta":   void,
		"param":  void,
		"source": void,
		"track":  void,
		"wbr":    void,
	}
)
