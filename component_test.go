package markup

import "testing"

type ValidCompo ZeroCompo

func (c *ValidCompo) Render() string {
	return `<p>Hello World</p>`
}

type EmptyCompo struct{}

func (c *EmptyCompo) Render() string {
	return `<p>Goodbye World</p>`
}

type NonPtrCompo ZeroCompo

func (c NonPtrCompo) Render() string {
	return `<p>Bye World</p>`
}

type IntCompo int

func (i *IntCompo) Render() string {
	return `<p>Aurevoir World</p>`
}

func TestEnsureValidCompo(t *testing.T) {
	valc := &ValidCompo{}
	if err := ensureValidComponent(valc); err != nil {
		t.Error(err)
	}

	noptrc := NonPtrCompo{}
	if err := ensureValidComponent(noptrc); err == nil {
		t.Error("err should not be nil")
	}

	empc := &EmptyCompo{}
	if err := ensureValidComponent(empc); err == nil {
		t.Error("err should not be nil")
	}

	intc := IntCompo(42)
	if err := ensureValidComponent(&intc); err == nil {
		t.Error("err should not be nil")
	}
}
