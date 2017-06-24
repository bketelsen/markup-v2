package markup

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

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

func mapComponentFields(c Componer, attrs map[string]string) error {
	if len(attrs) == 0 {
		return nil
	}

	v := reflect.ValueOf(c).Elem()
	t := v.Type()

	for i, numField := 0, t.NumField(); i < numField; i++ {
		f := v.Field(i)
		finfo := t.Field(i)

		if finfo.Anonymous {
			continue
		}

		if len(finfo.PkgPath) != 0 {
			continue
		}

		key := strings.ToLower(finfo.Name)
		val, ok := attrs[key]
		if !ok {
			if f.Kind() == reflect.Bool {
				f.SetBool(false)
			}
			continue
		}

		if err := mapComponentField(f, val); err != nil {
			return errors.Wrapf(err, `fail to map %s="%s" to %T.%s`, key, val, c, finfo.Name)
		}
	}
	return nil
}

func mapComponentField(f reflect.Value, v string) error {
	switch f.Kind() {
	case reflect.String:
		f.SetString(v)

	case reflect.Bool:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return err
		}
		f.SetBool(b)

	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		n, err := strconv.ParseInt(v, 0, 64)
		if err != nil {
			return err
		}
		f.SetInt(n)

	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uintptr:
		n, err := strconv.ParseUint(v, 0, 64)
		if err != nil {
			return err
		}
		f.SetUint(n)

	case reflect.Float64, reflect.Float32:
		n, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		f.SetFloat(n)

	default:
		addr := f.Addr()
		i := addr.Interface()
		if err := json.Unmarshal([]byte(v), i); err != nil {
			return err
		}
	}
	return nil
}
