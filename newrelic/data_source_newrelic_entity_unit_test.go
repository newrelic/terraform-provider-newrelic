//go:build unit
// +build unit

package newrelic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildTagsQueryFragment_SingleTag(t *testing.T) {
	t.Parallel()

	expected := "tags.`tagKey` = 'tagValue'"

	tags := []interface{}{
		map[string]interface{}{
			"key":   "tagKey",
			"value": "tagValue",
		},
	}

	result := buildTagsQueryFragment(tags)

	require.Equal(t, expected, result)
}

func TestBuildTagsQueryFragment_MultipleTags(t *testing.T) {
	t.Parallel()

	expected := "tags.`tagKey` = 'tagValue' AND tags.`tagKey2` = 'tagValue2' AND tags.`tagKey3` = 'tagValue3'"

	tags := []interface{}{
		map[string]interface{}{
			"key":   "tagKey",
			"value": "tagValue",
		},
		map[string]interface{}{
			"key":   "tagKey2",
			"value": "tagValue2",
		},
		map[string]interface{}{
			"key":   "tagKey3",
			"value": "tagValue3",
		},
	}

	result := buildTagsQueryFragment(tags)

	require.Equal(t, expected, result)
}

func TestBuildTagsQueryFragment_EmptyTags(t *testing.T) {
	t.Parallel()

	expected := ""
	tags := []interface{}{}

	result := buildTagsQueryFragment(tags)

	require.Equal(t, expected, result)
}

func TestBuildEntitySearchQuery(t *testing.T) {
	t.Parallel()

	tags := []interface{}{}

	// Name only
	expected := "name = 'Dummy App'"
	result := buildEntitySearchQuery("Dummy App", "", "", tags)
	require.Equal(t, expected, result)

	// Name & Domain
	expected = "name = 'Dummy App' AND domain = 'APM'"
	result = buildEntitySearchQuery("Dummy App", "APM", "", tags)
	require.Equal(t, expected, result)

	// Name, domain, and type
	expected = "name = 'Dummy App' AND domain = 'APM' AND type = 'APPLICATION'"
	result = buildEntitySearchQuery("Dummy App", "APM", "APPLICATION", tags)
	require.Equal(t, expected, result)

	// Name, domain, type, and tags
	expected = "name = 'Dummy App' AND domain = 'APM' AND type = 'APPLICATION' AND tags.`tagKey` = 'tagValue' AND tags.`tagKey2` = 'tagValue2'"
	tags = []interface{}{
		map[string]interface{}{
			"key":   "tagKey",
			"value": "tagValue",
		},
		map[string]interface{}{
			"key":   "tagKey2",
			"value": "tagValue2",
		},
	}
	result = buildEntitySearchQuery("Dummy App", "APM", "APPLICATION", tags)
	require.Equal(t, expected, result)
}
