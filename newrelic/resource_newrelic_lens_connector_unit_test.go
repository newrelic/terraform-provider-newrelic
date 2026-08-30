//go:build unit

package newrelic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// --- lensConnectorID / parseLensConnectorID ---

func TestLensConnectorID_Basic(t *testing.T) {
	t.Parallel()

	id := lensConnectorID("my-connector", lensScope{Type: "ORGANIZATION", ID: "org-123"})
	require.Equal(t, "my-connector:ORGANIZATION:org-123", id)
}

func TestParseLensConnectorID_Basic(t *testing.T) {
	t.Parallel()

	name, scopeType, scopeID := parseLensConnectorID("my-connector:ORGANIZATION:org-123")
	require.Equal(t, "my-connector", name)
	require.Equal(t, "ORGANIZATION", scopeType)
	require.Equal(t, "org-123", scopeID)
}

func TestLensConnectorID_RoundTrip(t *testing.T) {
	t.Parallel()

	scope := lensScope{Type: "ORGANIZATION", ID: "org-abc-456"}
	id := lensConnectorID("glue-prod", scope)
	name, scopeType, scopeID := parseLensConnectorID(id)

	require.Equal(t, "glue-prod", name)
	require.Equal(t, scope.Type, scopeType)
	require.Equal(t, scope.ID, scopeID)
}

func TestParseLensConnectorID_ScopeIDWithColons(t *testing.T) {
	t.Parallel()

	// scopeID itself contains colons — SplitN(3) must preserve them
	id := "my-connector:ORGANIZATION:org:uuid:with:colons"
	name, scopeType, scopeID := parseLensConnectorID(id)

	require.Equal(t, "my-connector", name)
	require.Equal(t, "ORGANIZATION", scopeType)
	require.Equal(t, "org:uuid:with:colons", scopeID)
}

func TestParseLensConnectorID_Malformed(t *testing.T) {
	t.Parallel()

	// Falls back gracefully when there are fewer than 3 segments
	name, scopeType, scopeID := parseLensConnectorID("no-separators")
	require.Equal(t, "no-separators", name)
	require.Equal(t, "", scopeType)
	require.Equal(t, "", scopeID)
}

// --- expandLensProperties ---

func TestExpandLensProperties_Empty(t *testing.T) {
	t.Parallel()

	result := expandLensProperties([]interface{}{})
	require.Empty(t, result)
}

func TestExpandLensProperties_Single(t *testing.T) {
	t.Parallel()

	raw := []interface{}{
		map[string]interface{}{"key": "region", "value": "us-east-1"},
	}
	result := expandLensProperties(raw)

	require.Len(t, result, 1)
	require.Equal(t, "region", result[0]["key"])
	require.Equal(t, "us-east-1", result[0]["value"])
}

func TestExpandLensProperties_Multiple(t *testing.T) {
	t.Parallel()

	raw := []interface{}{
		map[string]interface{}{"key": "region", "value": "us-east-1"},
		map[string]interface{}{"key": "database", "value": "mydb"},
	}
	result := expandLensProperties(raw)

	require.Len(t, result, 2)
	require.Equal(t, "region", result[0]["key"])
	require.Equal(t, "database", result[1]["key"])
}

// --- flattenLensProperties ---

func TestFlattenLensProperties_Empty(t *testing.T) {
	t.Parallel()

	result := flattenLensProperties([]lensProperty{})
	require.Empty(t, result)
}

func TestFlattenLensProperties_Single(t *testing.T) {
	t.Parallel()

	props := []lensProperty{{Key: "region", Value: "us-east-1"}}
	result := flattenLensProperties(props)

	require.Len(t, result, 1)
	m := result[0].(map[string]interface{})
	require.Equal(t, "region", m["key"])
	require.Equal(t, "us-east-1", m["value"])
}

func TestFlattenLensProperties_Multiple(t *testing.T) {
	t.Parallel()

	props := []lensProperty{
		{Key: "region", Value: "us-east-1"},
		{Key: "database", Value: "mydb"},
	}
	result := flattenLensProperties(props)

	require.Len(t, result, 2)
	require.Equal(t, "region", result[0].(map[string]interface{})["key"])
	require.Equal(t, "database", result[1].(map[string]interface{})["key"])
}

// --- flattenLensScope ---

func TestFlattenLensScope_Basic(t *testing.T) {
	t.Parallel()

	result := flattenLensScope(lensScope{ID: "org-123", Type: "ORGANIZATION"})

	require.Len(t, result, 1)
	m := result[0].(map[string]interface{})
	require.Equal(t, "org-123", m["id"])
	require.Equal(t, "ORGANIZATION", m["type"])
}

// --- flattenLensConnectors ---

func TestFlattenLensConnectors_Empty(t *testing.T) {
	t.Parallel()

	result := flattenLensConnectors([]lensConnectorCatalogItem{})
	require.Empty(t, result)
}

func TestFlattenLensConnectors_Single(t *testing.T) {
	t.Parallel()

	items := []lensConnectorCatalogItem{
		{
			Name:      "glue-prod",
			Connector: "AWSGLUE",
			Type:      "CATALOG",
			Scope:     lensScope{ID: "org-123", Type: "ORGANIZATION"},
			Properties: []lensProperty{
				{Key: "region", Value: "us-east-1"},
			},
		},
	}
	result := flattenLensConnectors(items)

	require.Len(t, result, 1)
	m := result[0].(map[string]interface{})
	require.Equal(t, "glue-prod", m["name"])
	require.Equal(t, "AWSGLUE", m["connector"])
	require.Equal(t, "CATALOG", m["type"])
	require.Len(t, m["properties"].([]interface{}), 1)
	require.Len(t, m["scope"].([]interface{}), 1)
}

func TestFlattenLensConnectors_Multiple(t *testing.T) {
	t.Parallel()

	items := []lensConnectorCatalogItem{
		{Name: "connector-a", Connector: "AWSGLUE", Scope: lensScope{ID: "org-1", Type: "ORGANIZATION"}},
		{Name: "connector-b", Connector: "AWSGLUE", Scope: lensScope{ID: "org-2", Type: "ORGANIZATION"}},
	}
	result := flattenLensConnectors(items)

	require.Len(t, result, 2)
	require.Equal(t, "connector-a", result[0].(map[string]interface{})["name"])
	require.Equal(t, "connector-b", result[1].(map[string]interface{})["name"])
}

// --- expand / flatten round-trip ---

func TestLensProperties_ExpandFlattenRoundTrip(t *testing.T) {
	t.Parallel()

	original := []interface{}{
		map[string]interface{}{"key": "region", "value": "us-east-1"},
		map[string]interface{}{"key": "database", "value": "mydb"},
	}

	expanded := expandLensProperties(original)

	props := make([]lensProperty, len(expanded))
	for i, p := range expanded {
		props[i] = lensProperty{Key: p["key"].(string), Value: p["value"].(string)}
	}

	flattened := flattenLensProperties(props)

	require.Len(t, flattened, 2)
	require.Equal(t, "region", flattened[0].(map[string]interface{})["key"])
	require.Equal(t, "database", flattened[1].(map[string]interface{})["key"])
}
