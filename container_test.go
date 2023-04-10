package goioc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContainerEmpty(t *testing.T) {
	c := NewContainer()

	assert.Equal(t, 1, len(c.bindings))
	assert.Equal(t, 0, len(*c.properties))
}

type TestFunc func()

func TestPropertyAvailableForInjection(t *testing.T) {
	c := NewContainer()
	SetContainer(c)

	found := false
	InjectBind[TestFunc](func(props Properties) func() {
		return func() {
			_, ok := props.String("test")
			if ok {
				found = true
			}
		}
	}, true)

	SetProperty("test", "test")

	fn, _ := Resolve[TestFunc](false)

	fn()

	assert.True(t, found, "Should properly inject Properties with all values")
}

type Work struct {
	description string
}

type Employee struct {
	name    string
	surname string
	age     int
	work    Work
}

func TestInjection(t *testing.T) {
	c := NewContainer()
	SetContainer(c)

	Bind[Work](Work{"typing on the keyboard"})

	InjectBind[Employee](func(wrk Work) Employee {
		return Employee{
			name:    "Dominik",
			surname: "Gawlik",
			age:     33,
			work:    wrk,
		}
	}, true)

	res, _ := Resolve[Employee](false)

	assert.Equal(t, Employee{
		name:    "Dominik",
		surname: "Gawlik",
		age:     33,
		work: Work{
			"typing on the keyboard",
		},
	}, res, "Injections should work")
}

type SomeFunc func(int) int
type SomeFunc2 func(int) bool

func TestInjectionWrongType(t *testing.T) {
	c := NewContainer()
	SetContainer(c)

	err := InjectBind[SomeFunc](func() func(x string) string {
		return func(x string) string {
			return x
		}
	}, true)

	assert.EqualError(t, err, "func(string) string is not convertible to target type goioc.SomeFunc", "Should reject invalid ctor definition")
}

func TestInjectionNestedBindingNotFound(t *testing.T) {
	c := NewContainer()
	SetContainer(c)

	InjectBind[SomeFunc](func(e Employee) func(x int) int {
		return func(x int) int {
			return x
		}
	}, true)

	_, err := Resolve[SomeFunc](false)

	assert.EqualError(t, err, "Binding for goioc.Employee not found", "Should reject invalid ctor definition")
}

func TestInjectionBindingNotFound(t *testing.T) {
	c := NewContainer()
	SetContainer(c)

	_, err := Resolve[SomeFunc](false)

	assert.EqualError(t, err, "Binding for goioc.SomeFunc not found", "Should reject invalid ctor definition")
}

func TestOverwriteBinding(t *testing.T) {
	c := NewContainer()
	SetContainer(c)

	Bind[SomeFunc](func(x int) int {
		return x * 2
	})

	Bind[SomeFunc](func(x int) int {
		return x * 4
	})

	v, _ := Resolve[SomeFunc](true)

	assert.Equal(t, 8, v(2), "Should overwrite dependency for same type")
}

func TestOverwriteBindInject(t *testing.T) {
	c := NewContainer()
	SetContainer(c)

	InjectBind[SomeFunc](func() func(x int) int {
		return func(x int) int {
			return x * 2
		}
	}, true)

	InjectBind[SomeFunc](func() func(x int) int {
		return func(x int) int {
			return x * 4
		}
	}, true)

	v, _ := Resolve[SomeFunc](true)

	assert.Equal(t, 8, v(2), "Should overwrite dependency for same type")
}

func TestValidationNotPassSimpleBinding(t *testing.T) {
	c := NewContainer()
	SetContainer(c)

	err := Bind[SomeFunc](func(x string) string {
		return x
	})

	assert.EqualError(t, err, "func(string) string is not convertible to target type goioc.SomeFunc", "Should reject invalid ctor definition")
}

func TestInjectResolveOk(t *testing.T) {
	c := NewContainer()
	SetContainer(c)

	Bind[SomeFunc](func(a int) int {
		return a * 2
	})

	res, _ := InjectResolve[SomeFunc2](func(fn SomeFunc) func(b int) bool {
		return func(b int) bool {
			return fn(1) == b
		}
	}, false)

	assert.Equal(t, true, res(2), "Should inject successfully")
}

func TestInjectResolveErr(t *testing.T) {
	c := NewContainer()
	SetContainer(c)

	Bind[SomeFunc](func(a int) int {
		return a * 2
	})

	_, err := InjectResolve[SomeFunc2](func(fn SomeFunc) func(b int) string {
		return func(b int) string {
			return "hello"
		}
	}, false)

	assert.EqualError(t, err, "func(int) string is not convertible to target type goioc.SomeFunc2", "Should reject validation")
}
