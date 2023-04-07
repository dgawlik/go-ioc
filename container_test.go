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
