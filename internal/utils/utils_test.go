// build +unit

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntArrayToString(t *testing.T) {
	t.Parallel()

	arr := []int{1, 2, 3, 4}
	result := IntArrayToString(arr)

	assert.Equal(t, "1,2,3,4", result)
}
