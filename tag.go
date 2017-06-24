package markup

import (
	"github.com/google/uuid"
	"golang.org/x/net/html/atom"
)

// Tag represents an HTML tag.
type Tag struct {
	ID       uuid.UUID
	CompoID  uuid.UUID
	Name     string
	Text     string
	Attrs    AttrMap
	Children []Tag
}

// IsEmpty reports whether its argument t is nil.
// Empty tags have empty name and empty text.
func (t *Tag) IsEmpty() bool {
	return len(t.Name) == 0 && len(t.Text) == 0
}

// IsText reports whether its argument t represents a text.
// Text tags have empty name and non empty text.
func (t *Tag) IsText() bool {
	if t.IsEmpty() {
		return false
	}
	return len(t.Name) == 0 && len(t.Text) != 0
}

// IsComponent reports whether its argument t represents a component.
// Component tags have non standard HTML5 tag name.
func (t *Tag) IsComponent() bool {
	if len(t.Name) == 0 {
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
	voidElems = map[string]struct{}{
		"area":   struct{}{},
		"base":   struct{}{},
		"br":     struct{}{},
		"col":    struct{}{},
		"embed":  struct{}{},
		"hr":     struct{}{},
		"img":    struct{}{},
		"input":  struct{}{},
		"keygen": struct{}{},
		"link":   struct{}{},
		"meta":   struct{}{},
		"param":  struct{}{},
		"source": struct{}{},
		"track":  struct{}{},
		"wbr":    struct{}{},
	}
)

// AttrMap represents a map of attributes.
type AttrMap map[string]string

// AttrEquals reports wheter its arguments l and r are equals.
func AttrEquals(l, r AttrMap) bool {
	if len(l) != len(r) {
		return false
	}

	for k, v := range l {
		otherVal, ok := r[k]
		if !ok {
			return false
		}
		if v != otherVal {
			return false
		}
	}
	return true
}
