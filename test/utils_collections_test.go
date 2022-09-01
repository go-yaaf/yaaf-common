// Copyright 2022. Motty Cohen
//
// Base configuration utility tests

package test

import (
	"github.com/go-yaaf/yaaf-common/utils/collections"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var str_array = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
var num_array = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

func TestCollections_Index(t *testing.T) {
	index := collections.Index(str_array, "4")
	assert.Equal(t, 4, index)
}

func TestCollections_IndexN(t *testing.T) {
	index := collections.IndexN(num_array, 4)
	assert.Equal(t, 4, index)
}

func TestCollections_Include(t *testing.T) {
	assert.True(t, collections.Include(str_array, "8"))
	assert.False(t, collections.Include(str_array, "18"))
}

func TestCollections_IncludeN(t *testing.T) {
	assert.True(t, collections.IncludeN(num_array, 8))
	assert.False(t, collections.IncludeN(num_array, 18))
}

func TestCollections_IncludeMask(t *testing.T) {

	var array = []int{8, 9, 10}

	assert.False(t, collections.IncludeMask(array, 1))
	assert.False(t, collections.IncludeMask(array, 2))
	assert.False(t, collections.IncludeMask(array, 3))
	assert.False(t, collections.IncludeMask(array, 4))
	assert.True(t, collections.IncludeMask(array, 8))
}

func TestCollections_Any(t *testing.T) {
	assert.True(t, collections.Any(str_array, func(s string) bool {
		return strings.Contains(s, "0")
	}))
	assert.False(t, collections.Any(str_array, func(s string) bool {
		return strings.Contains(s, "Y")
	}))
}

func TestCollections_AnyN(t *testing.T) {

	array := []int{5, 10, 11, 12}
	assert.True(t, collections.AnyN(array, func(n int) bool {
		return n%5 == 0
	}))
	assert.False(t, collections.AnyN(array, func(n int) bool {
		return n < 4
	}))
}
