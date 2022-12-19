//go:build unit
// +build unit

package newrelic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestGenerateNameForIntegrationTestResource(t *testing.T) {
	result := generateNameForIntegrationTestResource()

	require.Contains(t, result, "tf_test_")
}
