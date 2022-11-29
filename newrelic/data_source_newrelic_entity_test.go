//go:build unit
// +build unit

package newrelic

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/stretchr/testify/require"
)

func TestExpandEntityTag(t *testing.T) {
	flattened := []interface{}{
		map[string]interface{}{
			"key":   "my-key",
			"value": "my-value",
		},
	}

	expected := []entities.EntitySearchQueryBuilderTag{
		{
			Key:   "my-key",
			Value: "my-value",
		},
	}

	expanded := expandEntityTag(flattened)

	require.NotNil(t, expanded)
	require.Equal(t, expected, expanded)
}
