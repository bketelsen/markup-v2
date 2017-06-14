package markup

import (
	"io"
	"strings"

	"github.com/pkg/errors"

	"golang.org/x/net/html"
)

// TagDecoder is the interface that describes a decoder that can read HTML5 code
// and translate it to a Tag tree.
// Additionally, HTML5 can embed custom component tags.
type TagDecoder interface {
	Decode(t *Tag) error
}

// NewTagDecoder creates a new tad decoder.
func NewTagDecoder(r io.Reader) TagDecoder {
	return &tagDecoder{
		tokenizer: html.NewTokenizer(r),
	}
}

type tagDecoder struct {
	tokenizer *html.Tokenizer
	err       error
}

func (d *tagDecoder) Decode(t *Tag) error {
	d.decode(t)

	if t.IsEmpty() {
		return errors.New("can't decode an empty html")
	}

	return d.err
}

func (d *tagDecoder) decode(t *Tag) bool {
	z := d.tokenizer
	switch tok := z.Next(); tok {
	case html.StartTagToken:
		return d.decodeTag(t)

	case html.TextToken:
		return d.decodeText(t)

	case html.SelfClosingTagToken:
		return d.decodeSelfClosingTag(t)

	case html.ErrorToken:
		return false
	}
	return true
}

func (d *tagDecoder) decodeTag(t *Tag) bool {
	z := d.tokenizer

	bname, hasAttr := z.TagName()
	name := string(bname)
	t.Name = name

	if hasAttr {
		d.decodeAttrs(t)
	}

	if t.IsComponent() || t.IsVoidElem() {
		return true
	}

	for {
		c := Tag{}
		if !d.decode(&c) {
			return false
		}
		if c.IsEmpty() {
			return true
		}
		t.Children = append(t.Children, c)
	}
}

func (d *tagDecoder) decodeAttrs(t *Tag) {
	z := d.tokenizer

	attrs := make(map[string]string)
	for {
		key, val, more := z.TagAttr()
		attrs[string(key)] = string(val)
		if !more {
			break
		}
	}
	t.Attrs = attrs
}

func (d *tagDecoder) decodeText(t *Tag) bool {
	z := d.tokenizer

	text := string(z.Text())
	text = strings.TrimSpace(text)
	t.Text = text

	if t.IsEmpty() {
		return d.decode(t)
	}
	return true
}

func (d *tagDecoder) decodeSelfClosingTag(t *Tag) bool {
	z := d.tokenizer

	bname, _ := z.TagName()
	name := string(bname)
	d.err = errors.Errorf("%s should not be a closing tag", name)
	return false
}
