package markup

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

var (
	// Components allows to register and create new instances of components.
	Components CompoBuilder = make(compoBuilder)
)

// CompoBuilder is the interface that describes a component factory.
type CompoBuilder interface {
	// Register registers component of type c into the builder.
	// Components must be registered to be used.
	// During a rendering, it allows to create components of same type as c when
	// a tag named like c is found.
	Register(c Componer) (override bool)

	// New creates a component named n.
	New(n string) (c Componer, err error)
}

type compoBuilder map[string]reflect.Type

func (b compoBuilder) Register(c Componer) (override bool) {
	v := reflect.ValueOf(c)
	v = reflect.Indirect(v)
	t := v.Type()

	name := strings.ToLower(t.String())
	name = normalizeCompoName(name)

	_, override = b[name]
	b[name] = t
	return
}

func (b compoBuilder) New(name string) (c Componer, err error) {
	t, ok := b[name]
	if !ok {
		err = errors.Errorf("component %s is not registered", name)
		return
	}
	v := reflect.New(t)
	c = v.Interface().(Componer)
	return
}

func normalizeCompoName(name string) string {
	name = strings.ToLower(name)
	if pkgsep := strings.IndexByte(name, '.'); pkgsep != -1 {
		pkgname := name[:pkgsep]
		if pkgname == "main" {
			name = name[pkgsep+1:]
		}
	}
	return name
}
