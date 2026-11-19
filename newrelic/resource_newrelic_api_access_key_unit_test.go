//go:build unit

package newrelic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceNewRelicAPIAccessKey_KeyIsSensitive(t *testing.T) {
	t.Parallel()
	r := resourceNewRelicAPIAccessKey()
	s, ok := r.Schema["key"]
	require.True(t, ok, "key schema attribute must exist")
	// Without Sensitive, the key is printed in plain text in plans (e.g. on destroy).
	assert.True(t, s.Sensitive, "key schema attribute must be Sensitive")
}
