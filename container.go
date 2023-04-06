package goioc

import (
	"fmt"
	"reflect"
)

// Container holds and autowires all the registered types.
type Container struct {
	bindings   []Binding
	properties *map[string]any
}

// This is default container ready to use project-wise.
var DefaultContainer *Container = NewContainer()

// Returns fully initialized container.
func NewContainer() *Container {
	data := make(map[string]any)

	c := Container{
		bindings:   make([]Binding, 0),
		properties: &data,
	}

	targetType := reflect.TypeOf(*new(Properties))

	newB := newBinding(targetType, newProperties(&data), false)

	c.bindings = append(c.bindings, newB)

	return &c
}

func SetContainer(newC *Container) {
	DefaultContainer = newC
}

// It attaches any value to the container under some key.
// Properties can be injected in BindInject to inform
// autowiring decisions.
func SetProperty(key string, value any) {
	c := DefaultContainer

	data := c.properties
	(*data)[key] = value
}

// Creates simple binding T <==> value without any
// further injections.
func Bind[T any](value any) error {

	c := DefaultContainer

	targetType := reflect.TypeOf(*new(T))

	newB := newBinding(targetType, value, false)

	err := validate(newB, false)

	if err != nil {
		return err
	}

	idx := findBinding(c.bindings, targetType)
	if idx == -1 {
		c.bindings = append(c.bindings, newB)
	} else {
		c.bindings[idx] = newB
	}

	return nil
}

// Creates provider under type T which when injected and being resolved
// computes aproppirate value.
func BindInject[T any](value any) error {

	c := DefaultContainer

	targetType := reflect.TypeOf(*new(T))

	newB := newBinding(targetType, value, true)

	err := validate(newB, true)

	if err != nil {
		return err
	}

	idx := findBinding(c.bindings, targetType)
	if idx == -1 {
		c.bindings = append(c.bindings, newB)
	} else {
		c.bindings[idx] = newB
	}

	return nil
}

// Resolves what container holds under type T, fully autowired.
// Multiple calls result in caching of the dependencies, unless
// forcecRebind is specified true, in which case all dependencies
// are recomputed.
func Resolve[T any](forceRebind bool) (T, error) {

	c := DefaultContainer

	template := new(T)
	targetType := reflect.TypeOf(*template)

	value := reflect.ValueOf(template).Elem()

	idx := findBinding(c.bindings, targetType)

	if idx == -1 {
		return *template, fmt.Errorf("Binding for %v not found", reflect.TypeOf(value.Interface()))
	}

	_, err := c.bindings[idx].resolve(c.bindings, forceRebind)
	if err != nil {
		return *template, err
	}

	binding := reflect.ValueOf(c.bindings[idx].resolved)

	value.Set(binding)
	return *template, nil
}
