//go:build unit

package newrelic

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

// buildCtyObjectValue returns a cty object of objectType where the attributes named in set take the
// given values and every other attribute is null, mirroring what Terraform hands the SDK for a
// configuration that only mentions some of a resource's attributes.
func buildCtyObjectValue(objectType cty.Type, set map[string]cty.Value) cty.Value {
	values := make(map[string]cty.Value, len(objectType.AttributeTypes()))

	for name, attributeType := range objectType.AttributeTypes() {
		if value, ok := set[name]; ok {
			values[name] = value
			continue
		}
		values[name] = cty.NullVal(attributeType)
	}

	return cty.ObjectVal(values)
}

// newAzureIntegrationsResourceData builds a ResourceData for the azure integrations resource whose
// "monitor" block carries the given attributes, wired up the way the SDK does during an apply: the
// raw configuration is the cty value itself, while d.Get reads the flatmap shimmed from it.
//
// The block always carries metrics_polling_interval and enabled on top of monitorAttributes. A block
// whose every attribute is null shims down to a nil list element, which the expander short-circuits
// before it ever looks at the tags, so an all-null block would not exercise anything.
func newAzureIntegrationsResourceData(t *testing.T, monitorAttributes map[string]cty.Value) *schema.ResourceData {
	t.Helper()

	azureIntegrations := resourceNewRelicCloudAzureIntegrations()
	resourceType := azureIntegrations.CoreConfigSchema().ImpliedType()
	monitorType := resourceType.AttributeType("monitor").ElementType()

	monitorBlock := map[string]cty.Value{
		"metrics_polling_interval": cty.NumberIntVal(1200),
		"enabled":                  cty.True,
	}
	for name, value := range monitorAttributes {
		monitorBlock[name] = value
	}

	configValue := buildCtyObjectValue(resourceType, map[string]cty.Value{
		"id":                cty.StringVal("1234567"),
		"linked_account_id": cty.NumberIntVal(1234567),
		"monitor":           cty.ListVal([]cty.Value{buildCtyObjectValue(monitorType, monitorBlock)}),
	})

	state, err := azureIntegrations.ShimInstanceStateFromValue(configValue)
	require.NoError(t, err)
	state.RawConfig = configValue

	return azureIntegrations.Data(state)
}

// TestExpandCloudAzureIntegrationMonitorInputTags asserts the three-way distinction that
// CloudAzureMonitorIntegrationInput.IncludeTags/ExcludeTags being `*[]string` with `omitempty`
// makes expressible: a nil pointer omits the field and leaves the tags configured on the
// integration alone, while a pointer to an empty slice is serialized as `[]` and clears them.
//
// The wire forms are asserted directly, since that distinction -- not the Go value -- is what
// NerdGraph acts on. An explicit `null` is not a third option: the API rejects it with
// "ERR_INVALID_DATA include_tags must be a list".
func TestExpandCloudAzureIntegrationMonitorInputTags(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		tags     cty.Value
		expected *[]string
		// wire is what the field marshals to, or "" when it is omitted entirely.
		wire string
	}{
		{
			// Not present in the configuration: nothing should change.
			name:     "absent from configuration",
			tags:     cty.NullVal(cty.List(cty.String)),
			expected: nil,
			wire:     "",
		},
		{
			// Explicitly empty: the existing tags should be overridden with an empty list.
			name:     "explicitly set to an empty list",
			tags:     cty.ListValEmpty(cty.String),
			expected: &[]string{},
			wire:     "[]",
		},
		{
			name:     "populated",
			tags:     cty.ListVal([]cty.Value{cty.StringVal("env:production")}),
			expected: &[]string{"env:production"},
			wire:     `["env:production"]`,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			d := newAzureIntegrationsResourceData(t, map[string]cty.Value{
				"include_tags": testCase.tags,
				"exclude_tags": testCase.tags,
			})

			configureInput, _ := expandCloudAzureIntegrationsInput(d)

			require.Len(t, configureInput.Azure.AzureMonitor, 1)
			monitor := configureInput.Azure.AzureMonitor[0]

			require.Equal(t, testCase.expected, monitor.IncludeTags)
			require.Equal(t, testCase.expected, monitor.ExcludeTags)

			// A nil pointer and a pointer to an empty slice are the whole distinction the change
			// turns on, so be explicit rather than leaning on require.Equal.
			if testCase.expected == nil {
				require.Nil(t, monitor.IncludeTags)
				require.Nil(t, monitor.ExcludeTags)
			} else {
				require.NotNil(t, monitor.IncludeTags)
				require.NotNil(t, monitor.ExcludeTags)
			}

			// What actually reaches NerdGraph. `omitempty` on a pointer drops the field only when
			// the pointer is nil, so an empty slice still renders as `[]`.
			marshalled, err := json.Marshal(monitor)
			require.NoError(t, err)

			var wire map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(marshalled, &wire))

			for _, field := range []string{"includeTags", "excludeTags"} {
				raw, present := wire[field]
				if testCase.wire == "" {
					require.False(t, present, "%s should be omitted, got %s", field, raw)
					continue
				}
				require.True(t, present, "%s should be present", field)
				require.JSONEq(t, testCase.wire, string(raw))
			}
		})
	}
}

// TestExpandCloudAzureIntegrationMonitorInputTagsWithoutRawConfig covers the fallback path, where no
// raw configuration is available and GetRawConfig returns a null value. Absent tags are the safe
// reading there, since that leaves whatever is configured on the integration untouched.
func TestExpandCloudAzureIntegrationMonitorInputTagsWithoutRawConfig(t *testing.T) {
	t.Parallel()

	azureIntegrations := resourceNewRelicCloudAzureIntegrations()
	d := schema.TestResourceDataRaw(t, azureIntegrations.Schema, map[string]interface{}{
		"linked_account_id": 1234567,
		"monitor": []interface{}{
			map[string]interface{}{
				"metrics_polling_interval": 1200,
				"include_tags":             []interface{}{},
				"exclude_tags":             []interface{}{},
			},
		},
	})

	require.True(t, d.GetRawConfig().IsNull())

	configureInput, _ := expandCloudAzureIntegrationsInput(d)

	require.Len(t, configureInput.Azure.AzureMonitor, 1)
	require.Nil(t, configureInput.Azure.AzureMonitor[0].IncludeTags)
	require.Nil(t, configureInput.Azure.AzureMonitor[0].ExcludeTags)
}
