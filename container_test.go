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
	BindInject[TestFunc](func(props Properties) func() {
		return func() {
			_, ok := props.String("test")
			if ok {
				found = true
			}
		}
	})

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

	BindInject[Employee](func(wrk Work) Employee {
		return Employee{
			name:    "Dominik",
			surname: "Gawlik",
			age:     33,
			work:    wrk,
		}
	})

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

	err := BindInject[SomeFunc](func() func(x string) string {
		return func(x string) string {
			return x
		}
	})

	assert.EqualError(t, err, "func(string) string is not convertible to target type goioc.SomeFunc", "Should reject invalid ctor definition")
}

func TestInjectionNestedBindingNotFound(t *testing.T) {
	c := NewContainer()
	SetContainer(c)

	BindInject[SomeFunc](func(e Employee) func(x int) int {
		return func(x int) int {
			return x
		}
	})

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

	BindInject[SomeFunc](func() func(x int) int {
		return func(x int) int {
			return x * 2
		}
	})

	BindInject[SomeFunc](func() func(x int) int {
		return func(x int) int {
			return x * 4
		}
	})

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
