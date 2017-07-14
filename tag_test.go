package markup

import (
	"testing"
)

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

func TestTagHTML(t *testing.T) {
	b := NewCompoBuilder()
	b.Register(&Hello{})
	b.Register(&World{})

	env := newEnv(b)

	hello := &Hello{
		Name: "JonhyMaxoo",
	}

	root, err := env.Mount(hello)
	if err != nil {
		t.Fatal(err)
	}

	h, err := root.HTML(env)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(h)

	if _, err = root.HTML(nil); err == nil {
		t.Fatal("err should not be nil")
	}
	t.Log(err)

	errRoot := Tag{
		Name: "markup.hello",
	}
	if _, err = errRoot.HTML(env); err == nil {
		t.Fatal("err should not be nil")
	}
	t.Log(err)
}

func BenchmarkTagHTML(b *testing.B) {
	bui := NewCompoBuilder()
	bui.Register(&Hello{})
	bui.Register(&World{})

	env := newEnv(bui)

	hello := &Hello{
		Name: "JonhyMaxoo",
	}

	root, _ := env.Mount(hello)
	for i := 0; i < b.N; i++ {
		root.HTML(env)
	}
}

type Hello2 struct {
	Greeting      string
	Name          string
	Placeholder   string
	TextBye       bool
	TmplErr       bool
	ChildErr      bool
	CompoFieldErr bool
}

func (h *Hello2) Render() string {
	return `
<div>
	<h1>{{html .Greeting}}</h1>
	<input type="text" placeholder="{{.Placeholder}}" onchange="Name" />
	<p>
		{{if .Name}}
			<World2 name="{{html .Name}}" err="{{.ChildErr}}" {{if .CompoFieldErr}}fielderr="-42"{{end}} />
		{{else}}
			<span>World</span>
		{{end}}
	</p>

	{{if .TmplErr}}
		<div>{{.UnknownField}}</div>
	{{end}}

	{{if .TextBye}}
		Goodbye
	{{else}}
		<span>Goodbye</span>
		<p>world</p>
	{{end}}
</div>
	`
}

type World2 struct {
	Name     string
	Err      bool
	FieldErr uint
}

func (w *World2) Render() string {
	return `
<div>
	{{html .Name}}

	{{if .Err}}
		<markup.componotregistered>
	{{end}}
</div>
	`
}
