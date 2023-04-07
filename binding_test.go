package goioc

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBindingIsCtor(t *testing.T) {
	binding := newBinding(reflect.TypeOf(1), 1, true)

	assert.Equal(t, Binding{
		targetType: reflect.TypeOf(1),
		resolved:   nil,
		ctor:       1,
	}, binding, "Ctor binding should have resolved nil")
}

func TestNewBindingIsNotCtor(t *testing.T) {
	binding := newBinding(reflect.TypeOf(1), 1, false)

	assert.Equal(t, Binding{
		targetType: reflect.TypeOf(1),
		resolved:   1,
		ctor:       nil,
	}, binding, "Plain binding should have ctor nil")
}

func TestValidatorNotCtor(t *testing.T) {
	fn := func(x int) (int, int) {
		return x, x
	}

	binding := newBinding(reflect.TypeOf(fn).Out(0), fn, true)

	err := validate(binding, true)

	assert.EqualError(t, err, "func(int) (int, int) is not valid constructor prototype", "Should reject if not ctor")
}

func TestValidatorNotCtor2(t *testing.T) {
	fn := func(x int) (func(int) int, bool) {
		return func(x int) int {
			return x
		}, false
	}

	binding := newBinding(reflect.TypeOf(fn).Out(0), fn, true)

	err := validate(binding, true)

	assert.EqualError(t, err, "func(int) (func(int) int, bool) is not valid constructor prototype", "Should reject if not ctor")
}

func TestValidatorCtor(t *testing.T) {
	fn := func(x int) func(int) int {
		return func(x int) int {
			return x
		}
	}

	binding := newBinding(reflect.TypeOf(fn).Out(0), fn, true)

	err := validate(binding, true)

	assert.Nil(t, err, "On valid ctor should not return error")
}

func TestValidatorPlainNotConvertible(t *testing.T) {
	fn1 := func(x int) int {
		return x
	}

	fn2 := func(x string) string {
		return x
	}

	binding := newBinding(reflect.TypeOf(fn1), fn2, false)

	err := validate(binding, false)

	assert.EqualError(t, err, "func(string) string is not convertible to target type func(int) int", "Should not validate non convertible type")
}

func TestFindBindingMatches(t *testing.T) {
	bindings := [2]Binding{
		{
			targetType: reflect.TypeOf(1),
			ctor:       nil,
			resolved:   nil,
		},
		{
			targetType: reflect.TypeOf("Hello"),
			ctor:       nil,
			resolved:   nil,
		},
	}

	idx := findBinding(bindings[:], reflect.TypeOf(1))

	assert.Equal(t, 0, idx, "Should find binding under type")
}

func TestFindBindingNotMatches(t *testing.T) {
	bindings := [2]Binding{
		{
			targetType: reflect.TypeOf(1),
			ctor:       nil,
			resolved:   nil,
		},
		{
			targetType: reflect.TypeOf("Hello"),
			ctor:       nil,
			resolved:   nil,
		},
	}

	idx := findBinding(bindings[:], reflect.TypeOf(false))

	assert.Equal(t, -1, idx, "Should find binding under type")
}

func TestResolveCached(t *testing.T) {
	b := Binding{
		targetType: reflect.TypeOf(1),
		resolved:   1,
		ctor:       nil,
	}

	result, _ := b.resolve(make([]Binding, 0), false)

	assert.Equal(t, 1, result, "Should return value from cache")
}

func TestResolveNullConstructor(t *testing.T) {
	b := Binding{
		targetType: reflect.TypeOf(1),
		resolved:   1,
		ctor:       nil,
	}

	result, _ := b.resolve(make([]Binding, 0), true)

	assert.Equal(t, 1, result, "Should return value from cache")
}

type Type1 int
type Type2 struct {
	value1 Type1
}

type Type3 struct {
	value2 Type2
}

func TestResolveEverytingOk(t *testing.T) {
	b := Binding{
		targetType: reflect.TypeOf(*new(Type2)),
		resolved:   nil,
		ctor: func(value1 Type1) Type2 {
			return Type2{
				value1: value1,
			}
		},
	}

	bindings := [1]Binding{
		{
			targetType: reflect.TypeOf(*new(Type1)),
			resolved:   10,
			ctor:       nil,
		},
	}

	res, _ := b.resolve(bindings[:], true)

	assert.Equal(t, Type2{
		value1: 10,
	}, res, "Should resolve dependency correctly")
}

func TestResolveDependencyNotFound(t *testing.T) {
	b := Binding{
		targetType: reflect.TypeOf(*new(Type2)),
		resolved:   nil,
		ctor: func(value1 Type1) Type2 {
			return Type2{
				value1: value1,
			}
		},
	}

	bindings := [1]Binding{
		{
			targetType: reflect.TypeOf(*new(Type3)),
			resolved:   10,
			ctor:       nil,
		},
	}

	_, err := b.resolve(bindings[:], true)

	assert.EqualError(t, err, "Binding for goioc.Type1 not found", "Should throw on dependency not found")
}

func TestResolveNestedDependencyNotFound(t *testing.T) {
	b := Binding{
		targetType: reflect.TypeOf(*new(Type3)),
		resolved:   nil,
		ctor: func(value1 Type2) Type3 {
			return Type3{
				value2: value1,
			}
		},
	}

	bindings := [1]Binding{
		{
			targetType: reflect.TypeOf(*new(Type2)),
			resolved:   0,
			ctor: func(value1 Type1) Type2 {
				return Type2{
					value1: value1,
				}
			},
		},
	}

	_, err := b.resolve(bindings[:], true)

	assert.EqualError(t, err, "Binding for goioc.Type1 not found", "Should throw on dependency not found")
}
