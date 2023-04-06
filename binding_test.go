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
