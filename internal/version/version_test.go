// +build unit

package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var GitTag = "undefined"

func TestVersionTag(t *testing.T) {
	t.Parallel()

	assert.Equal(t, Version, GitTag)
}
