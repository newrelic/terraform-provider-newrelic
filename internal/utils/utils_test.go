// build +unit

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntArrayToString(t *testing.T) {
	t.Parallel()

	var result string

	// empty
	result = IntArrayToString([]int{})
	assert.Equal(t, "", result)

	// single
	result = IntArrayToString([]int{1})
	assert.Equal(t, "1", result)

	// multiple
	result = IntArrayToString([]int{1, 2, 3, 4})
	assert.Equal(t, "1,2,3,4", result)
}
