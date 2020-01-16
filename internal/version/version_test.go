// +build unit

package version

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var GitTag = "undefined"

func TestVersionTag(t *testing.T) {
	t.Parallel()

	if strings.HasPrefix(os.Getenv("CIRCLE_BRANCH"), "release/") {
		t.Skip("skipping version test due to release branch")
	}

	assert.Equal(t, Version, GitTag)
}
