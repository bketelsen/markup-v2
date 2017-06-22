package markup

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Env is the interface that describes an environment that handle component life
// cycle.
type Env interface {
	// Component returns the component mounted under the identifier id.
	// err should be set if there is no mounted component under id.
	Component(id uuid.UUID) (c Componer, err error)

	// Mount indexes the component c into the env.
	// The component will live until it is dismounted.
	//
	// Mount should call the Render method from the component and create a tree
	// of Tag.
	// An id should be assigned to each tags.
	// Tags describing other components should trigger their creation and their
	// mount.
	Mount(c Componer) (root Tag, err error)

	// Dismount removes references to a component and its children.
	Dismount(c Componer)
}

// NewEnv creates a new environment.
func NewEnv(b CompoBuilder) Env {
	return newEnv(b)
}

type env struct {
	components   map[uuid.UUID]Componer
	compoRoots   map[Componer]Tag
	compoBuilder CompoBuilder
}

func newEnv(b CompoBuilder) *env {
	return &env{
		components:   make(map[uuid.UUID]Componer),
		compoRoots:   make(map[Componer]Tag),
		compoBuilder: b,
	}
}

func (e *env) Component(id uuid.UUID) (c Componer, err error) {
	ok := false
	if c, ok = e.components[id]; !ok {
		err = errors.Errorf("no component with id %v is mounted", id)
	}
	return
}

func (e *env) Mount(c Componer) (root Tag, err error) {
	if _, ok := e.compoRoots[c]; ok {
		err = errors.Errorf("%T is already mounted", c)
		return
	}

	r := c.Render()
	tmpl := template.Must(template.New(fmt.Sprintf("%T", c)).Parse(r))

	b := bytes.Buffer{}
	if err = tmpl.Execute(&b, e); err != nil {
		err = errors.Wrapf(err, "fail to execute render from %T", c)
		return
	}

	dec := NewTagDecoder(&b)
	if err = dec.Decode(&root); err != nil {
		err = errors.Wrapf(err, "fail to decode render from %T", c)
		return
	}

	compoID := uuid.New()
	if err = e.mountTag(&root, compoID); err != nil {
		err = errors.Wrapf(err, "fail to mount %T", c)
		return
	}

	e.components[compoID] = c
	e.compoRoots[c] = root
	return
}

func (e *env) mountTag(t *Tag, compoID uuid.UUID) error {
	t.ID = uuid.New()
	t.CompoID = compoID

	if t.IsText() {
		return nil
	}

	if t.IsComponent() {
		c, err := e.compoBuilder.New(t.Name)
		if err != nil {
			return errors.Wrapf(err, "fail to create %s", t.Name)
		}

		root, err := e.Mount(c)
		if err != nil {
			return errors.Wrapf(err, "fail to mount %s", t.Name)
		}
		t.ID = root.CompoID
		return nil
	}

	for i := range t.Children {
		if err := e.mountTag(&t.Children[i], compoID); err != nil {
			return errors.Wrapf(err, "fail to mount %s child", t.Name)
		}
	}
	return nil
}

func (e *env) Dismount(c Componer) {
	root, ok := e.compoRoots[c]
	if !ok {
		return
	}

	e.dismountTag(&root)
	delete(e.components, root.CompoID)
	delete(e.compoRoots, c)
	return
}

func (e *env) dismountTag(t *Tag) {
	if t.IsComponent() {
		c, err := e.Component(t.ID)
		if err != nil {
			return
		}

		e.Dismount(c)
		return
	}

	for i := range t.Children {
		e.dismountTag(&t.Children[i])
	}
	return
}
