package goioc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTypePropertyFound(t *testing.T) {
	data := make(map[string]any)
	data["test1"] = "Hello"

	v, _ := getType[string](&data, "test1")

	assert.Equal(t, "Hello", v, "Should return correct value")
}

func TestGetTypePropertyNotFound(t *testing.T) {
	data := make(map[string]any)
	data["test1"] = "Hello"

	_, ok := getType[string](&data, "test2")

	assert.False(t, ok, "Should return  false on key not present")
}

func TestGetTypePropertyWrongType(t *testing.T) {
	data := make(map[string]any)
	data["test1"] = "Hello"

	_, ok := getType[int](&data, "test1")

	assert.False(t, ok, "Should return  false on wrong type")
}
