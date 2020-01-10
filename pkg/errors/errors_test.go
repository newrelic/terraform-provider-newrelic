// +build unit

package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorNotFound(t *testing.T) {
	t.Parallel()

	var e ErrorNotFound

	assert.Equal(t, "404 not found", e.Error())
}

func TestErrorUnexpectedStatusCode(t *testing.T) {
	t.Parallel()

	e := ErrorUnexpectedStatusCode{
		StatusCode: 99,
		Err:        "wat",
	}

	assert.Equal(t, "99 response returned: wat", e.Error())
}
