//go:build unit
// +build unit

package newrelic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseIDs_Basic(t *testing.T) {
	ids, err := parseIDs("1:2", 2)

	require.NoError(t, err)
	require.Equal(t, 2, len(ids))
	require.Equal(t, 1, ids[0])
	require.Equal(t, 2, ids[1])
}

func TestParseIDs_BadIDs(t *testing.T) {
	_, err := parseIDs("12", 2)
	require.Error(t, err)

	_, err = parseIDs("a:b", 2)
	require.Error(t, err)
}

func TestParseHashedIDs_Basic(t *testing.T) {
	expected := []int{1, 2, 3}
	result, err := parseHashedIDs("1:2:3")

	require.NoError(t, err)
	require.Equal(t, 3, len(result))
	require.Equal(t, expected, result)
}

func TestParseHashedIDs_Invalid(t *testing.T) {
	_, err := parseHashedIDs("123:abc")

	require.Error(t, err)
}

func TestSerializeIDs_Basic(t *testing.T) {
	id := serializeIDs([]int{1, 2})

	require.Equal(t, "1:2", id)
}

func TestStripWhitespace(t *testing.T) {
	json := " { \"key\": \"value\" } "
	e := "{\"key\":\"value\"}"
	a := stripWhitespace(json)

	require.Equal(t, e, a)
}

func TestSortIntegerSlice(t *testing.T) {
	integers := []int{2, 1, 4, 3}
	expected := []int{1, 2, 3, 4}

	sortIntegerSlice(integers)

	require.Equal(t, expected, integers)
}
