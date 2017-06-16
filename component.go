package markup

import (
	"reflect"

	"github.com/pkg/errors"
)

// Componer is the interface that describes a component.
// Should be implemented on a non empty struct pointer.
type Componer interface {
	// Render should return a string describing the component with HTML5
	// standard.
	// It support Golang template/text API.
	// Pipeline is based on the component struct.
	// See https://golang.org/pkg/text/template for more informations.
	Render() string
}

// ZeroCompo is the type to redefine when writing an empty component.
// Every instances of an empty struct is given the same memory address, which
// causes problem for indexing components.
// ZeroCompo have a placeholder field to avoid that.
type ZeroCompo struct {
	placeholder byte
}

func ensureValidComponent(c Componer) error {
	v := reflect.ValueOf(c)
	if v.Kind() != reflect.Ptr {
		return errors.Errorf("%T must be implemented on a struct pointer", c)
	}

	if v = v.Elem(); v.Kind() != reflect.Struct {
		return errors.Errorf("%T must be implemented on a struct pointer", c)
	}

	if v.NumField() == 0 {
		return errors.Errorf("%T can't be implemented on an empty struct pointer", c)
	}
	return nil
}
