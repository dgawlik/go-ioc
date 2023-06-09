package goioc

import (
	"fmt"
	"reflect"
)

// Container holds and autowires all the registered types.
// The autowiring is lazy and it happens on Resolve call.
type Container struct {
	bindings   []Binding
	properties *map[string]any
}

// This is default container implicitly used by the methods.
var DefaultContainer *Container = NewContainer()

// Returns ready to use new Container.
func NewContainer() *Container {
	data := make(map[string]any)

	c := Container{
		bindings:   make([]Binding, 0),
		properties: &data,
	}

	targetType := reflect.TypeOf(*new(Properties))

	newB := newBinding(targetType, newProperties(&data), false, true)

	c.bindings = append(c.bindings, newB)

	return &c
}

// Sets new implicit Container.
func SetContainer(newC *Container) {
	DefaultContainer = newC
}

// It attaches any value to the container under given key.
// It can be then injected through Properties internal type in BindInject to inform
// autowiring decisions.
func SetProperty(key string, value any) {
	c := DefaultContainer

	data := c.properties
	(*data)[key] = value
}

// Creates plain binding T <==> value without any
// further nested injections.
func Bind[T any](value any) error {

	c := DefaultContainer

	targetType := reflect.TypeOf(*new(T))

	newB := newBinding(targetType, value, false, false)

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

// Creates constructor on injections under type T, which when being resolved
// computes aproppirate value which nested bindings.
func InjectBind[T any](value any, isPrototype bool) error {

	c := DefaultContainer

	targetType := reflect.TypeOf(*new(T))

	newB := newBinding(targetType, value, true, isPrototype)

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
	return resolveInternal[T](nil, forceRebind)
}

// It differs from resolve that the injected value is not in cache
// you provide constructor for it. It returns injected value without
// putting it to cache
func InjectResolve[T any](ctor any, forceRebind bool) (T, error) {
	return resolveInternal[T](ctor, forceRebind)
}

func resolveInternal[T any](ctor any, forceRebind bool) (T, error) {
	c := DefaultContainer

	template := new(T)
	targetType := reflect.TypeOf(*template)

	value := reflect.ValueOf(template).Elem()

	var newB Binding

	if ctor != nil {
		newB = newBinding(targetType, ctor, true, false)

		err := validate(newB, true)

		if err != nil {
			return *template, err
		}

	} else {

		idx := findBinding(c.bindings, targetType)

		if idx == -1 {
			return *template, fmt.Errorf("Binding for %v not found", reflect.TypeOf(value.Interface()))
		}

		newB = c.bindings[idx]
	}

	result, err := newB.resolve(c.bindings, forceRebind)
	if err != nil {
		return *template, err
	}

	value.Set(reflect.ValueOf(result))
	return *template, nil
}
