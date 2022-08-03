//go:build unit
// +build unit

package newrelic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func TestMergeSchemas(t *testing.T) {
	schema1 := map[string]*schema.Schema{
		"string_attribute": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "test",
			Description: "This attribute has a TypeString",
		},
		"boolean_attribute": {
			Type:        schema.TypeBool,
			Description: "This attribute is a TypeBool",
		},
	}

	schema2 := map[string]*schema.Schema{
		"typeset_attribute": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "This attribute is a TypeSet",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "This nested attribute is a TypeString",
					},
				},
			},
		},
	}

	result := mergeSchemas(schema1, schema2)
	require.Equal(t, 3, len(result))

	defaultStringAttrValue, err := result["string_attribute"].DefaultValue()
	require.NoError(t, err)
	require.Equal(t, "test", defaultStringAttrValue)
}
