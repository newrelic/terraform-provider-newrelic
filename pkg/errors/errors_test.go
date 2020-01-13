// +build unit

package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorNotFound(t *testing.T) {
	t.Parallel()

	var e NotFound

	assert.Equal(t, "resource not found", e.Error())
}

func TestErrorUnexpectedStatusCode(t *testing.T) {
	t.Parallel()

	e := NewUnexpectedStatusCode(99, "wat")

	assert.Equal(t, "99 response returned: wat", e.Error())
}
