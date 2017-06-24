package markup

import (
	"testing"

	"github.com/google/uuid"
)

type Foo ZeroCompo

func (c *Foo) Render() string {
	return `
<div>
	<h1>Foo</h1>
	<markup.bar>
</div>
	`
}

type Bar ZeroCompo

func (c *Bar) Render() string {
	return `<h2>Bar</h2>`
}

type CompoBadTmpl ZeroCompo

func (c *CompoBadTmpl) Render() string {
	return `<h2>{{.Hello}}</h2>`
}

type CompoBadTag ZeroCompo

func (c *CompoBadTag) Render() string {
	return `<h1><div/></h1>`
}

type CompoNotRegistered ZeroCompo

func (c *CompoNotRegistered) Render() string {
	return `
<div>
	<markup.unknown>
</div>
	`
}

type CompoBadChild ZeroCompo

func (c *CompoBadChild) Render() string {
	return `
<div>
	<markup.compobadtmpl>
</div>
	`
}

func TestNewEnv(t *testing.T) {
	b := NewCompoBuilder()
	NewEnv(b)
}

func TestEnvComponent(t *testing.T) {
	compoID := uuid.New()
	foo := &Foo{}

	b := NewCompoBuilder()
	env := newEnv(b)
	env.components[compoID] = foo

	c, err := env.Component(compoID)
	if err != nil {
		t.Fatal(err)
	}
	if c != foo {
		t.Fatal("c and foo should point to the same component")
	}

	if _, err = env.Component(uuid.New()); err == nil {
		t.Fatal("err should not be nil")
	}
}

func TestEnv(t *testing.T) {
	b := NewCompoBuilder()
	b.Register(&Foo{})
	b.Register(&Bar{})
	b.Register(&CompoBadTmpl{})
	b.Register(&CompoBadTag{})
	b.Register(&CompoNotRegistered{})
	b.Register(&CompoBadChild{})

	env := newEnv(b)

	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "mount and dismount",
			test: func(t *testing.T) { testMountDismount(t, env, &Foo{}) },
		},
		{
			name: "mount mounted",
			test: func(t *testing.T) { testMountMounted(t, env, &Foo{}) },
		},
		{
			name: "mount component with bad template",
			test: func(t *testing.T) { testMountInvalid(t, env, &CompoBadTmpl{}) },
		},
		{
			name: "mount component with bad tag",
			test: func(t *testing.T) { testMountInvalid(t, env, &CompoBadTag{}) },
		},
		{
			name: "mount component with not registered child",
			test: func(t *testing.T) { testMountInvalid(t, env, &CompoNotRegistered{}) },
		},
		{
			name: "mount component with bad child",
			test: func(t *testing.T) { testMountInvalid(t, env, &CompoBadChild{}) },
		},
		{
			name: "dismount dismounted",
			test: func(t *testing.T) { testDismountDismounted(t, env, &Foo{}) },
		},
		{
			name: "dismount dismounted child",
			test: func(t *testing.T) { testDismountDismountedChild(t, env, &Foo{}) },
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.test)
	}
}

func testMountDismount(t *testing.T, env *env, c Componer) {
	// Mount.
	root, err := env.Mount(c)
	if err != nil {
		t.Fatal(err)
	}
	if count := len(env.components); count != 2 {
		t.Fatal("env shoud have 2 components:", count)
	}
	if count := len(env.compoRoots); count != 2 {
		t.Fatal("env shoud have 2 component roots:", count)
	}

	barTag := root.Children[1]
	if name := barTag.Name; name != "markup.bar" {
		t.Fatalf(`barTag.Name should be "markup.bar": "%s"`, name)
	}
	if _, err = env.Component(barTag.ID); err != nil {
		t.Fatal(err)
	}

	// Dismount
	env.Dismount(c)
	if count := len(env.components); count != 0 {
		t.Fatal("env shoud have 0 component:", count)
	}
	if count := len(env.compoRoots); count != 0 {
		t.Fatal("env shoud have 0 component root:", count)
	}
}

func testMountMounted(t *testing.T, env *env, c Componer) {
	if _, err := env.Mount(c); err != nil {
		t.Fatal(err)
	}
	defer env.Dismount(c)

	_, err := env.Mount(c)
	if err == nil {
		t.Fatal("err should not be nil")
	}
	t.Log(err)
}

func testMountInvalid(t *testing.T, env *env, c Componer) {
	_, err := env.Mount(c)
	if err == nil {
		t.Fatal("err should not be nil")
	}
	t.Log(err)
}

func testDismountDismounted(t *testing.T, env *env, c Componer) {
	if _, err := env.Mount(c); err != nil {
		t.Fatal(err)
	}
	env.Dismount(c)
	env.Dismount(c)
}

func testDismountDismountedChild(t *testing.T, env *env, c Componer) {
	root, err := env.Mount(c)
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range env.components {
		if k != root.CompoID {
			env.Dismount(v)
		}
	}
	env.Dismount(c)
}
