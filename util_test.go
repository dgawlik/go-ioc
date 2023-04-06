package goioc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindExists(t *testing.T) {
	coll := [4]int{1, 2, 3, 4}

	idx := find(coll[:], func(x int) bool {
		return x == 3
	})

	assert.Equal(t, 2, idx, "Should find valid index position")
}

func TestFindDoesNotExist(t *testing.T) {
	coll := [4]int{1, 2, 3, 4}

	idx := find(coll[:], func(x int) bool {
		return x == 5
	})

	assert.Equal(t, -1, idx, "Not found should be equal to -1")
}
