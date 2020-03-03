package newrelic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExpandChannelIDs(t *testing.T) {
	flattened := []interface{}{123, 456}
	expected := []int{123, 456}

	expanded := expandChannelIDs(flattened)

	require.NotNil(t, expanded)
	require.Equal(t, expected, expanded)
}
