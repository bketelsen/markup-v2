package markup

import (
	"reflect"
	"strings"
	"testing"
)

type RegisterTest struct {
}

func (c *RegisterTest) Render() string {
	return `<p>Hello World</p>`
}

func TestCompoBuilderRegister(t *testing.T) {
	c := &RegisterTest{}
	ct := reflect.TypeOf(*c)
	cname := strings.ToLower(ct.String())

	b := make(compoBuilder)
	if b.Register(c) {
		t.Fatalf("%s should not be overriden", cname)
	}

	if _, ok := b[cname]; !ok {
		t.Fatalf("%s should have been registered", cname)
	}

	if !b.Register(c) {
		t.Fatalf("%s should have been overriden", cname)
	}
}

func TestCompoBuilderNew(t *testing.T) {
	c := &RegisterTest{}
	cname := "markup.registertest"
	b := make(compoBuilder)
	b.Register(c)

	n, err := b.New(cname)
	if err != nil {
		t.Fatal(err)
	}
	if n == nil {
		t.Fatalf("%s should have been created: %v", cname, n)
	}

	if n, err = b.New("unknown"); err == nil {
		t.Fatal("unknown should not have been created")
	}
}

func TestNormalizeCompoName(t *testing.T) {
	if name := "lib.FooBar"; normalizeCompoName(name) != "lib.foobar" {
		t.Errorf(`name should be "lib.foobar": "%s"`, name)
	}

	if name := "main.FooBar"; normalizeCompoName(name) != "foobar" {
		t.Errorf(`name should be "foobar": "%s"`, name)
	}
}
