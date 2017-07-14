package markup

import (
	"bytes"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
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
		"area":   {},
		"base":   {},
		"br":     {},
		"col":    {},
		"embed":  {},
		"hr":     {},
		"img":    {},
		"input":  {},
		"keygen": {},
		"link":   {},
		"meta":   {},
		"param":  {},
		"source": {},
		"track":  {},
		"wbr":    {},
	}
)

// HTML returns the HTML5 representation of t.
func (t *Tag) HTML(env Env) (h string, err error) {
	var b bytes.Buffer
	if err = t.print(&b, 0, env); err != nil {
		return
	}

	h = b.String()
	return
}

func (t *Tag) print(b *bytes.Buffer, indent int, env Env) error {
	if env == nil {
		return errors.New("env is not set")
	}

	if t.IsText() {
		t.printIndent(b, indent)
		b.WriteString(t.Text)
		return nil
	}

	if t.IsComponent() {
		return t.printComponent(b, indent, env)
	}

	t.printIndent(b, indent)
	b.WriteString("<")
	b.WriteString(t.Name)
	t.printAttributes(b)
	b.WriteRune('>')

	if t.IsVoidElem() {
		return nil
	}

	if len(t.Children) == 0 {
		b.WriteString("</")
		b.WriteString(t.Name)
		b.WriteRune('>')
		return nil
	}

	for _, child := range t.Children {
		b.WriteRune('\n')
		child.print(b, indent+1, env)
	}

	b.WriteRune('\n')
	t.printIndent(b, indent)
	b.WriteString("</")
	b.WriteString(t.Name)
	b.WriteRune('>')
	return nil
}

func (t *Tag) printComponent(b *bytes.Buffer, indent int, env Env) error {
	c, err := env.Component(t.ID)
	if err != nil {
		return errors.Wrap(err, "can't print component")
	}

	root, _ := env.Root(c)
	return root.print(b, indent, env)
}

func (t *Tag) printAttributes(b *bytes.Buffer) {
	for k, v := range t.Attrs {
		if len(v) == 0 {
			b.WriteRune(' ')
			b.WriteString(k)
			continue
		}

		if strings.HasPrefix(k, "on") {
			b.WriteRune(' ')
			b.WriteString(k)
			b.WriteString(`="CallGoHandler('`)
			b.WriteString(t.CompoID.String())
			b.WriteString(`', '`)
			b.WriteString(v)
			b.WriteString(`', this, event)"`)
			continue
		}

		b.WriteRune(' ')
		b.WriteString(k)
		b.WriteString(`="`)
		b.WriteString(v)
		b.WriteString(`"`)
	}

	b.WriteString(` data-go-id="`)
	b.WriteString(t.ID.String())
	b.WriteString(`"`)
}

func (t *Tag) printIndent(b *bytes.Buffer, indent int) {
	for i := 0; i < indent; i++ {
		b.WriteString("  ")
	}
}

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
