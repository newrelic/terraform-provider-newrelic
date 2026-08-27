//go:build unit

package newrelic

import (
	"encoding/base64"
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/pkg/pathpoint"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── accountIDFromGUID ─────────────────────────────────────────────────────────

func TestAccountIDFromGUID_Valid(t *testing.T) {
	raw := "12345|NR1_ENTITY|NA|some-id"
	guid := base64.StdEncoding.EncodeToString([]byte(raw))
	assert.Equal(t, 12345, accountIDFromGUID(guid))
}

func TestAccountIDFromGUID_InvalidBase64(t *testing.T) {
	assert.Equal(t, 0, accountIDFromGUID("not-valid-base64!!!"))
}

func TestAccountIDFromGUID_NoSeparator(t *testing.T) {
	guid := base64.StdEncoding.EncodeToString([]byte("nopipe"))
	// SplitN with n=2 always returns at least 1 element; Atoi will fail on "nopipe"
	assert.Equal(t, 0, accountIDFromGUID(guid))
}

func TestAccountIDFromGUID_NonNumericAccount(t *testing.T) {
	guid := base64.StdEncoding.EncodeToString([]byte("abc|NR1_ENTITY|NA|id"))
	assert.Equal(t, 0, accountIDFromGUID(guid))
}

// ── expandPathpointKpiNRQLSelectInput ─────────────────────────────────────────

func TestExpandPathpointKpiNRQLSelectInput_Basic(t *testing.T) {
	m := map[string]interface{}{
		"aggregation_type": "COUNT",
		"alias":            "",
		"attribute":        "",
		"threshold":        float64(0),
	}
	got := expandPathpointKpiNRQLSelectInput(m)
	assert.Equal(t, pathpoint.PathPointKpiNRQLAggregations("COUNT"), got.AggregationType)
	assert.Empty(t, got.Alias)
	assert.Empty(t, got.Attribute)
	assert.Equal(t, float64(0), got.Threshold)
}

func TestExpandPathpointKpiNRQLSelectInput_AllFields(t *testing.T) {
	m := map[string]interface{}{
		"aggregation_type": "PERCENTILE",
		"alias":            "p99",
		"attribute":        "duration",
		"threshold":        float64(99.0),
	}
	got := expandPathpointKpiNRQLSelectInput(m)
	assert.Equal(t, pathpoint.PathPointKpiNRQLAggregations("PERCENTILE"), got.AggregationType)
	assert.Equal(t, "p99", got.Alias)
	assert.Equal(t, "duration", got.Attribute)
	assert.Equal(t, float64(99.0), got.Threshold)
}

// ── expandPathpointKpiTimeWindowInput ─────────────────────────────────────────

func TestExpandPathpointKpiTimeWindowInput_Empty(t *testing.T) {
	m := map[string]interface{}{
		"custom_range":   "",
		"relative_range": []interface{}{},
	}
	assert.Nil(t, expandPathpointKpiTimeWindowInput(m))
}

func TestExpandPathpointKpiTimeWindowInput_CustomRange(t *testing.T) {
	m := map[string]interface{}{
		"custom_range":   "SINCE 3 days ago",
		"relative_range": []interface{}{},
	}
	got := expandPathpointKpiTimeWindowInput(m)
	require.NotNil(t, got)
	assert.Equal(t, pathpoint.NRQL("SINCE 3 days ago"), got.CustomRange)
	assert.Nil(t, got.RelativeRange)
}

func TestExpandPathpointKpiTimeWindowInput_RelativeRange(t *testing.T) {
	m := map[string]interface{}{
		"custom_range": "",
		"relative_range": []interface{}{
			map[string]interface{}{
				"since":           "SEVEN_DAYS",
				"compare_against": "THIRTY_DAYS",
			},
		},
	}
	got := expandPathpointKpiTimeWindowInput(m)
	require.NotNil(t, got)
	require.NotNil(t, got.RelativeRange)
	assert.Equal(t, pathpoint.PathPointKpiTimeDuration("SEVEN_DAYS"), got.RelativeRange.Since)
	assert.Equal(t, pathpoint.PathPointKpiTimeDuration("THIRTY_DAYS"), got.RelativeRange.CompareAgainst)
}

func TestExpandPathpointKpiTimeWindowInput_RelativeRangeNoCompare(t *testing.T) {
	m := map[string]interface{}{
		"custom_range": "",
		"relative_range": []interface{}{
			map[string]interface{}{
				"since":           "SIXTY_MINUTES",
				"compare_against": "",
			},
		},
	}
	got := expandPathpointKpiTimeWindowInput(m)
	require.NotNil(t, got)
	require.NotNil(t, got.RelativeRange)
	assert.Equal(t, pathpoint.PathPointKpiTimeDuration("SIXTY_MINUTES"), got.RelativeRange.Since)
	assert.Empty(t, got.RelativeRange.CompareAgainst)
}

// ── expandPathpointKpiNRQLInput ───────────────────────────────────────────────

func TestExpandPathpointKpiNRQLInput_Basic(t *testing.T) {
	m := map[string]interface{}{
		"from":  "Transaction",
		"where": "",
		"select": []interface{}{
			map[string]interface{}{
				"aggregation_type": "COUNT",
				"alias":            "",
				"attribute":        "",
				"threshold":        float64(0),
			},
		},
		"time_window": []interface{}{},
	}
	got := expandPathpointKpiNRQLInput(m)
	assert.Equal(t, pathpoint.NRQL("Transaction"), got.From)
	assert.Equal(t, pathpoint.PathPointKpiNRQLAggregations("COUNT"), got.Select.AggregationType)
	assert.Nil(t, got.TimeWindow)
}

func TestExpandPathpointKpiNRQLInput_WithWhere(t *testing.T) {
	m := map[string]interface{}{
		"from":  "Metric",
		"where": "appName = 'myApp'",
		"select": []interface{}{
			map[string]interface{}{
				"aggregation_type": "SUM",
				"alias":            "",
				"attribute":        "value",
				"threshold":        float64(0),
			},
		},
		"time_window": []interface{}{},
	}
	got := expandPathpointKpiNRQLInput(m)
	assert.Equal(t, "appName = 'myApp'", got.Where)
}

// ── expandPathpointKpiInputs ──────────────────────────────────────────────────

func TestExpandPathpointKpiInputs_Empty(t *testing.T) {
	got := expandPathpointKpiInputs([]interface{}{}, 1111)
	assert.Empty(t, got)
}

func TestExpandPathpointKpiInputs_UsesFlowAccountIDWhenZero(t *testing.T) {
	raw := []interface{}{
		map[string]interface{}{
			"name":        "Request Rate",
			"description": "",
			"category":    "",
			"account_id":  0,
			"query": []interface{}{
				map[string]interface{}{
					"from":        "Transaction",
					"where":       "",
					"time_window": []interface{}{},
					"select": []interface{}{
						map[string]interface{}{
							"aggregation_type": "COUNT",
							"alias":            "",
							"attribute":        "",
							"threshold":        float64(0),
						},
					},
				},
			},
		},
	}
	got := expandPathpointKpiInputs(raw, 9999)
	require.Len(t, got, 1)
	assert.Equal(t, "Request Rate", got[0].Name)
	assert.Equal(t, 9999, got[0].AccountID)
}

func TestExpandPathpointKpiInputs_UsesKpiAccountIDWhenProvided(t *testing.T) {
	raw := []interface{}{
		map[string]interface{}{
			"name":        "Error Rate",
			"description": "desc",
			"category":    "SRE",
			"account_id":  4321,
			"query": []interface{}{
				map[string]interface{}{
					"from":        "TransactionError",
					"where":       "",
					"time_window": []interface{}{},
					"select": []interface{}{
						map[string]interface{}{
							"aggregation_type": "COUNT",
							"alias":            "",
							"attribute":        "",
							"threshold":        float64(0),
						},
					},
				},
			},
		},
	}
	got := expandPathpointKpiInputs(raw, 1111)
	require.Len(t, got, 1)
	assert.Equal(t, 4321, got[0].AccountID)
	assert.Equal(t, "desc", got[0].Description)
	assert.Equal(t, "SRE", got[0].Category)
}

// ── expandPathpointKpiUpdateInputsResolved ────────────────────────────────────

func TestExpandPathpointKpiUpdateInputsResolved_MatchByName(t *testing.T) {
	oldRaw := []interface{}{
		map[string]interface{}{"id": "kpi-1", "name": "Request Rate", "account_id": 0, "description": "", "category": "", "query": []interface{}{}},
	}
	newRaw := []interface{}{
		map[string]interface{}{
			"name": "Request Rate", "account_id": 0, "description": "", "category": "",
			"query": []interface{}{
				map[string]interface{}{
					"from": "Transaction", "where": "", "time_window": []interface{}{},
					"select": []interface{}{map[string]interface{}{"aggregation_type": "COUNT", "alias": "", "attribute": "", "threshold": float64(0)}},
				},
			},
		},
	}
	got := expandPathpointKpiUpdateInputsResolved(newRaw, oldRaw, 1111)
	require.Len(t, got, 1)
	assert.Equal(t, "kpi-1", got[0].ID)
	assert.Equal(t, "Request Rate", got[0].Name)
}

func TestExpandPathpointKpiUpdateInputsResolved_MatchByPosition(t *testing.T) {
	oldRaw := []interface{}{
		map[string]interface{}{"id": "kpi-old", "name": "Old Name", "account_id": 0, "description": "", "category": "", "query": []interface{}{}},
	}
	newRaw := []interface{}{
		map[string]interface{}{
			"name": "New Name", "account_id": 0, "description": "", "category": "",
			"query": []interface{}{
				map[string]interface{}{
					"from": "Transaction", "where": "", "time_window": []interface{}{},
					"select": []interface{}{map[string]interface{}{"aggregation_type": "COUNT", "alias": "", "attribute": "", "threshold": float64(0)}},
				},
			},
		},
	}
	got := expandPathpointKpiUpdateInputsResolved(newRaw, oldRaw, 1111)
	require.Len(t, got, 1)
	assert.Equal(t, "kpi-old", got[0].ID, "should fall back to positional match")
	assert.Equal(t, "New Name", got[0].Name)
}

func TestExpandPathpointKpiUpdateInputsResolved_NewKpiNoID(t *testing.T) {
	got := expandPathpointKpiUpdateInputsResolved(
		[]interface{}{
			map[string]interface{}{
				"name": "Brand New", "account_id": 0, "description": "", "category": "",
				"query": []interface{}{
					map[string]interface{}{
						"from": "Transaction", "where": "", "time_window": []interface{}{},
						"select": []interface{}{map[string]interface{}{"aggregation_type": "COUNT", "alias": "", "attribute": "", "threshold": float64(0)}},
					},
				},
			},
		},
		[]interface{}{},
		1111,
	)
	require.Len(t, got, 1)
	assert.Empty(t, got[0].ID)
}

// ── expandPathpointRelatedInput ───────────────────────────────────────────────

func TestExpandPathpointRelatedInput(t *testing.T) {
	got := expandPathpointRelatedInput(map[string]interface{}{"source": true, "target": false})
	assert.True(t, got.Source)
	assert.False(t, got.Target)
}

// ── expandPathpointSignalQueryInput ──────────────────────────────────────────

func TestExpandPathpointSignalQueryInput(t *testing.T) {
	got := expandPathpointSignalQueryInput(map[string]interface{}{
		"query":       "domain = 'APM'",
		"is_excluded": true,
	})
	require.NotNil(t, got)
	assert.Equal(t, "domain = 'APM'", got.Query)
	assert.True(t, got.IsExcluded)
}

// ── expandPathpointStepStatusThresholdInput ───────────────────────────────────

func TestExpandPathpointStepStatusThresholdInput_Empty(t *testing.T) {
	got := expandPathpointStepStatusThresholdInput(map[string]interface{}{
		"health_rollup":   "",
		"threshold_type":  "",
		"threshold_value": 0,
	})
	assert.Empty(t, string(got.HealthRollup))
	assert.Empty(t, string(got.ThresholdType))
	assert.Equal(t, 0, got.ThresholdValue)
}

func TestExpandPathpointStepStatusThresholdInput_Full(t *testing.T) {
	got := expandPathpointStepStatusThresholdInput(map[string]interface{}{
		"health_rollup":   "WORST_STATUS_WINS",
		"threshold_type":  "PERCENTAGE",
		"threshold_value": 75,
	})
	assert.Equal(t, pathpoint.PathPointStepHealthRollup("WORST_STATUS_WINS"), got.HealthRollup)
	assert.Equal(t, pathpoint.PathPointThresholdType("PERCENTAGE"), got.ThresholdType)
	assert.Equal(t, 75, got.ThresholdValue)
}

// ── expandPathpointSignalInputs ───────────────────────────────────────────────

func TestExpandPathpointSignalInputs_Empty(t *testing.T) {
	got := expandPathpointSignalInputs([]interface{}{})
	assert.Empty(t, got)
}

func TestExpandPathpointSignalInputs_Full(t *testing.T) {
	raw := []interface{}{
		map[string]interface{}{
			"guid":        "MjM0NTZ8TlIxX0FQUHxBUFBMSUNBVElPTnwxMjM0NTY",
			"name":        "My Service",
			"type":        "ENTITY",
			"is_excluded": false,
		},
		map[string]interface{}{
			"guid":        "alert-guid",
			"name":        "",
			"type":        "ALERT",
			"is_excluded": true,
		},
	}
	got := expandPathpointSignalInputs(raw)
	require.Len(t, got, 2)
	assert.Equal(t, pathpoint.EntityGUID("MjM0NTZ8TlIxX0FQUHxBUFBMSUNBVElPTnwxMjM0NTY"), got[0].GUID)
	assert.Equal(t, "My Service", got[0].Name)
	assert.Equal(t, pathpoint.PathPointSignalType("ENTITY"), got[0].Type)
	assert.False(t, got[0].IsExcluded)
	assert.True(t, got[1].IsExcluded)
}

// ── expandPathpointStepInputs ─────────────────────────────────────────────────

func TestExpandPathpointStepInputs_Basic(t *testing.T) {
	raw := []interface{}{
		map[string]interface{}{
			"name":                "My Step",
			"is_excluded":         false,
			"link":                "",
			"scoped_accounts":     []interface{}{},
			"entity_search_query": []interface{}{},
			"config":              []interface{}{},
			"signals":             []interface{}{},
		},
	}
	got := expandPathpointStepInputs(raw)
	require.Len(t, got, 1)
	assert.Equal(t, "My Step", got[0].Name)
	assert.False(t, got[0].IsExcluded)
	assert.Nil(t, got[0].EntitySearchQuery)
}

func TestExpandPathpointStepInputs_WithOptionals(t *testing.T) {
	raw := []interface{}{
		map[string]interface{}{
			"name":            "Step With Extras",
			"is_excluded":     true,
			"link":            "https://example.com",
			"scoped_accounts": []interface{}{111, 222},
			"entity_search_query": []interface{}{
				map[string]interface{}{
					"query":       "type = 'APM'",
					"is_excluded": false,
				},
			},
			"config": []interface{}{
				map[string]interface{}{
					"health_rollup":   "BEST_STATUS_WINS",
					"threshold_type":  "FIXED",
					"threshold_value": 5,
				},
			},
			"signals": []interface{}{},
		},
	}
	got := expandPathpointStepInputs(raw)
	require.Len(t, got, 1)
	s := got[0]
	assert.True(t, s.IsExcluded)
	assert.Equal(t, "https://example.com", s.Link)
	assert.Equal(t, []int{111, 222}, s.ScopedAccounts)
	require.NotNil(t, s.EntitySearchQuery)
	assert.Equal(t, "type = 'APM'", s.EntitySearchQuery.Query)
	assert.Equal(t, pathpoint.PathPointStepHealthRollup("BEST_STATUS_WINS"), s.Config.HealthRollup)
}

// ── expandPathpointStepUpdateInputsResolved ───────────────────────────────────

func TestExpandPathpointStepUpdateInputsResolved_MatchByName(t *testing.T) {
	oldRaw := []interface{}{
		map[string]interface{}{
			"id": "step-1", "name": "Step A",
			"is_excluded": false, "link": "", "scoped_accounts": []interface{}{},
			"entity_search_query": []interface{}{}, "config": []interface{}{}, "signals": []interface{}{},
		},
	}
	newRaw := []interface{}{
		map[string]interface{}{
			"name":        "Step A",
			"is_excluded": false, "link": "", "scoped_accounts": []interface{}{},
			"entity_search_query": []interface{}{}, "config": []interface{}{}, "signals": []interface{}{},
		},
	}
	got := expandPathpointStepUpdateInputsResolved(newRaw, oldRaw)
	require.Len(t, got, 1)
	assert.Equal(t, "step-1", got[0].ID)
}

func TestExpandPathpointStepUpdateInputsResolved_PositionalFallback(t *testing.T) {
	oldRaw := []interface{}{
		map[string]interface{}{
			"id": "step-pos", "name": "Old Step",
			"is_excluded": false, "link": "", "scoped_accounts": []interface{}{},
			"entity_search_query": []interface{}{}, "config": []interface{}{}, "signals": []interface{}{},
		},
	}
	newRaw := []interface{}{
		map[string]interface{}{
			"name":        "Renamed Step",
			"is_excluded": false, "link": "", "scoped_accounts": []interface{}{},
			"entity_search_query": []interface{}{}, "config": []interface{}{}, "signals": []interface{}{},
		},
	}
	got := expandPathpointStepUpdateInputsResolved(newRaw, oldRaw)
	require.Len(t, got, 1)
	assert.Equal(t, "step-pos", got[0].ID)
}

// ── expandPathpointLevelInputs ────────────────────────────────────────────────

func TestExpandPathpointLevelInputs_Empty(t *testing.T) {
	got := expandPathpointLevelInputs([]interface{}{})
	assert.Empty(t, got)
}

func TestExpandPathpointLevelInputs_WithSteps(t *testing.T) {
	raw := []interface{}{
		map[string]interface{}{
			"steps": []interface{}{
				map[string]interface{}{
					"name":        "Step One",
					"is_excluded": false, "link": "", "scoped_accounts": []interface{}{},
					"entity_search_query": []interface{}{}, "config": []interface{}{}, "signals": []interface{}{},
				},
			},
		},
	}
	got := expandPathpointLevelInputs(raw)
	require.Len(t, got, 1)
	require.Len(t, got[0].Steps, 1)
	assert.Equal(t, "Step One", got[0].Steps[0].Name)
}

// ── expandPathpointLevelUpdateInputsResolved ──────────────────────────────────

func TestExpandPathpointLevelUpdateInputsResolved_MatchByStepNames(t *testing.T) {
	oldRaw := []interface{}{
		map[string]interface{}{
			"id": "level-1",
			"steps": []interface{}{
				map[string]interface{}{
					"id": "s1", "name": "Step A",
					"is_excluded": false, "link": "", "scoped_accounts": []interface{}{},
					"entity_search_query": []interface{}{}, "config": []interface{}{}, "signals": []interface{}{},
				},
			},
		},
	}
	newRaw := []interface{}{
		map[string]interface{}{
			"steps": []interface{}{
				map[string]interface{}{
					"name":        "Step A",
					"is_excluded": false, "link": "", "scoped_accounts": []interface{}{},
					"entity_search_query": []interface{}{}, "config": []interface{}{}, "signals": []interface{}{},
				},
			},
		},
	}
	got := expandPathpointLevelUpdateInputsResolved(newRaw, oldRaw)
	require.Len(t, got, 1)
	assert.Equal(t, "level-1", got[0].ID)
	require.Len(t, got[0].Steps, 1)
	assert.Equal(t, "s1", got[0].Steps[0].ID)
}

// ── expandPathpointStageInputs ────────────────────────────────────────────────

func TestExpandPathpointStageInputs_Basic(t *testing.T) {
	raw := []interface{}{
		map[string]interface{}{
			"name":          "Stage One",
			"health_rollup": "",
			"is_excluded":   false,
			"link":          "",
			"related":       []interface{}{},
			"stage_kpis":    []interface{}{},
			"levels":        []interface{}{},
		},
	}
	got := expandPathpointStageInputs(raw, 1111)
	require.Len(t, got, 1)
	assert.Equal(t, "Stage One", got[0].Name)
	assert.Empty(t, got[0].Levels)
}

func TestExpandPathpointStageInputs_WithRelatedAndLink(t *testing.T) {
	raw := []interface{}{
		map[string]interface{}{
			"name":          "Source Stage",
			"health_rollup": "AUTOMATIC_ROLL_UP",
			"is_excluded":   false,
			"link":          "https://wiki.example.com",
			"related": []interface{}{
				map[string]interface{}{
					"source": true,
					"target": false,
				},
			},
			"stage_kpis": []interface{}{},
			"levels":     []interface{}{},
		},
	}
	got := expandPathpointStageInputs(raw, 1111)
	require.Len(t, got, 1)
	assert.Equal(t, pathpoint.PathPointStageHealthRollup("AUTOMATIC_ROLL_UP"), got[0].HealthRollup)
	assert.Equal(t, "https://wiki.example.com", got[0].Link)
	assert.True(t, got[0].Related.Source)
	assert.False(t, got[0].Related.Target)
}

// ── expandPathpointStageUpdateInputsResolved ──────────────────────────────────

func TestExpandPathpointStageUpdateInputsResolved_MatchByName(t *testing.T) {
	oldRaw := []interface{}{
		map[string]interface{}{
			"id": "stage-1", "name": "Checkout",
			"health_rollup": "", "is_excluded": false, "link": "",
			"related": []interface{}{}, "stage_kpis": []interface{}{}, "levels": []interface{}{},
		},
	}
	newRaw := []interface{}{
		map[string]interface{}{
			"name":          "Checkout",
			"health_rollup": "", "is_excluded": false, "link": "",
			"related": []interface{}{}, "stage_kpis": []interface{}{}, "levels": []interface{}{},
		},
	}
	got := expandPathpointStageUpdateInputsResolved(newRaw, oldRaw, 1111)
	require.Len(t, got, 1)
	assert.Equal(t, "stage-1", got[0].ID)
}

// ── flattenPathpointSignals ───────────────────────────────────────────────────

func TestFlattenPathpointSignals_Empty(t *testing.T) {
	got := flattenPathpointSignals([]pathpoint.PathPointSignal{})
	assert.Empty(t, got)
}

func TestFlattenPathpointSignals(t *testing.T) {
	signals := []pathpoint.PathPointSignal{
		{
			GUID:       "guid-1",
			Name:       "My Signal",
			Type:       "ENTITY",
			IsExcluded: false,
		},
	}
	got := flattenPathpointSignals(signals)
	require.Len(t, got, 1)
	assert.Equal(t, "guid-1", got[0]["guid"])
	assert.Equal(t, "My Signal", got[0]["name"])
	assert.Equal(t, "ENTITY", got[0]["type"])
	assert.Equal(t, false, got[0]["is_excluded"])
}

// ── flattenPathpointSteps ─────────────────────────────────────────────────────

func TestFlattenPathpointSteps_Empty(t *testing.T) {
	got := flattenPathpointSteps([]pathpoint.PathPointStep{})
	assert.Empty(t, got)
}

func TestFlattenPathpointSteps_Basic(t *testing.T) {
	steps := []pathpoint.PathPointStep{
		{
			ID:         "step-1",
			Name:       "My Step",
			IsExcluded: false,
			Link:       "",
		},
	}
	got := flattenPathpointSteps(steps)
	require.Len(t, got, 1)
	assert.Equal(t, "step-1", got[0]["id"])
	assert.Equal(t, "My Step", got[0]["name"])
}

func TestFlattenPathpointSteps_WithEntitySearchQuery(t *testing.T) {
	steps := []pathpoint.PathPointStep{
		{
			ID:   "step-2",
			Name: "Filtered Step",
			EntitySearchQuery: pathpoint.PathPointSignalQuery{
				Query:      "domain='APM'",
				IsExcluded: true,
			},
		},
	}
	got := flattenPathpointSteps(steps)
	require.Len(t, got, 1)
	esq, ok := got[0]["entity_search_query"].([]map[string]interface{})
	require.True(t, ok)
	require.Len(t, esq, 1)
	assert.Equal(t, "domain='APM'", esq[0]["query"])
	assert.Equal(t, true, esq[0]["is_excluded"])
}

func TestFlattenPathpointSteps_WithConfig(t *testing.T) {
	steps := []pathpoint.PathPointStep{
		{
			ID:   "step-3",
			Name: "Configured Step",
			Config: pathpoint.PathPointStepStatusThreshold{
				HealthRollup:   "BEST_STATUS_WINS",
				ThresholdType:  "FIXED",
				ThresholdValue: 3,
			},
		},
	}
	got := flattenPathpointSteps(steps)
	require.Len(t, got, 1)
	cfg, ok := got[0]["config"].([]map[string]interface{})
	require.True(t, ok)
	require.Len(t, cfg, 1)
	assert.Equal(t, "BEST_STATUS_WINS", cfg[0]["health_rollup"])
	assert.Equal(t, "FIXED", cfg[0]["threshold_type"])
	assert.Equal(t, 3, cfg[0]["threshold_value"])
}

// ── flattenPathpointLevels ────────────────────────────────────────────────────

func TestFlattenPathpointLevels_Empty(t *testing.T) {
	got := flattenPathpointLevels([]pathpoint.PathPointLevel{})
	assert.Empty(t, got)
}

func TestFlattenPathpointLevels(t *testing.T) {
	levels := []pathpoint.PathPointLevel{
		{
			ID: "level-1",
			Steps: pathpoint.PathPointStepsItems{
				Items: []pathpoint.PathPointStep{
					{ID: "step-1", Name: "Step One"},
				},
			},
		},
	}
	got := flattenPathpointLevels(levels)
	require.Len(t, got, 1)
	assert.Equal(t, "level-1", got[0]["id"])
	steps, ok := got[0]["steps"].([]map[string]interface{})
	require.True(t, ok)
	require.Len(t, steps, 1)
	assert.Equal(t, "Step One", steps[0]["name"])
}

// ── flattenPathpointKpiTimeWindow ─────────────────────────────────────────────

func TestFlattenPathpointKpiTimeWindow_Empty(t *testing.T) {
	tw := pathpoint.PathPointKpiTimeWindow{}
	got := flattenPathpointKpiTimeWindow(tw)
	assert.Nil(t, got)
}

func TestFlattenPathpointKpiTimeWindow_CustomRange(t *testing.T) {
	tw := pathpoint.PathPointKpiTimeWindow{
		CustomRange: "SINCE 3 days ago",
	}
	got := flattenPathpointKpiTimeWindow(tw)
	require.Len(t, got, 1)
	assert.Equal(t, "SINCE 3 days ago", got[0]["custom_range"])
	_, hasRR := got[0]["relative_range"]
	assert.False(t, hasRR)
}

func TestFlattenPathpointKpiTimeWindow_RelativeRange(t *testing.T) {
	tw := pathpoint.PathPointKpiTimeWindow{
		RelativeRange: pathpoint.PathPointKpiTimeWindowRelativeRange{
			Since:          "SEVEN_DAYS",
			CompareAgainst: "THIRTY_DAYS",
		},
	}
	got := flattenPathpointKpiTimeWindow(tw)
	require.Len(t, got, 1)
	rr, ok := got[0]["relative_range"].([]map[string]interface{})
	require.True(t, ok)
	require.Len(t, rr, 1)
	assert.Equal(t, "SEVEN_DAYS", rr[0]["since"])
	assert.Equal(t, "THIRTY_DAYS", rr[0]["compare_against"])
}

// ── flattenPathpointKpiNRQL ───────────────────────────────────────────────────

func TestFlattenPathpointKpiNRQL(t *testing.T) {
	q := pathpoint.PathPointKpiNRQL{
		From:  "Transaction",
		Where: "appName = 'test'",
		Select: pathpoint.PathPointKpiNRQLSelect{
			AggregationType: "COUNT",
			Alias:           "reqs",
			Attribute:       "",
			Threshold:       0,
		},
	}
	got := flattenPathpointKpiNRQL(q)
	require.Len(t, got, 1)
	assert.Equal(t, "Transaction", got[0]["from"].(string))
	assert.Equal(t, "appName = 'test'", got[0]["where"])
	sel := got[0]["select"].([]map[string]interface{})
	require.Len(t, sel, 1)
	assert.Equal(t, "COUNT", sel[0]["aggregation_type"])
	assert.Equal(t, "reqs", sel[0]["alias"])
}

// ── flattenPathpointKpis ──────────────────────────────────────────────────────

func TestFlattenPathpointKpis_Empty(t *testing.T) {
	got := flattenPathpointKpis([]pathpoint.PathPointKpi{})
	assert.Empty(t, got)
}

func TestFlattenPathpointKpis(t *testing.T) {
	kpis := []pathpoint.PathPointKpi{
		{
			ID:          "kpi-1",
			Name:        "Error Rate",
			Description: "Rate of errors",
			Category:    "SRE",
			AccountID:   1234,
			MetricQuery: "SELECT count(*) FROM Metric",
			Query: pathpoint.PathPointKpiNRQL{
				From: "TransactionError",
				Select: pathpoint.PathPointKpiNRQLSelect{
					AggregationType: "COUNT",
				},
			},
		},
	}
	got := flattenPathpointKpis(kpis)
	require.Len(t, got, 1)
	assert.Equal(t, "kpi-1", got[0]["id"])
	assert.Equal(t, "Error Rate", got[0]["name"])
	assert.Equal(t, "Rate of errors", got[0]["description"])
	assert.Equal(t, "SRE", got[0]["category"])
	assert.Equal(t, 1234, got[0]["account_id"])
	assert.Equal(t, "SELECT count(*) FROM Metric", got[0]["metric_query"])
}

// ── flattenPathpointStages ────────────────────────────────────────────────────

func TestFlattenPathpointStages_Empty(t *testing.T) {
	got := flattenPathpointStages([]pathpoint.PathPointStage{})
	assert.Empty(t, got)
}

func TestFlattenPathpointStages_Basic(t *testing.T) {
	stages := []pathpoint.PathPointStage{
		{
			ID:           "stage-1",
			Name:         "Checkout",
			HealthRollup: "AUTOMATIC_ROLL_UP",
			IsExcluded:   false,
			Link:         "https://wiki.example.com",
			Related:      pathpoint.PathPointRelated{Source: false, Target: false},
		},
	}
	got := flattenPathpointStages(stages)
	require.Len(t, got, 1)
	assert.Equal(t, "stage-1", got[0]["id"])
	assert.Equal(t, "Checkout", got[0]["name"])
	assert.Equal(t, "AUTOMATIC_ROLL_UP", got[0]["health_rollup"])
	assert.Equal(t, "https://wiki.example.com", got[0]["link"])
	_, hasRelated := got[0]["related"]
	assert.False(t, hasRelated, "related should not be set when both source and target are false")
}

func TestFlattenPathpointStages_WithRelated(t *testing.T) {
	stages := []pathpoint.PathPointStage{
		{
			ID:      "stage-2",
			Name:    "Source Stage",
			Related: pathpoint.PathPointRelated{Source: true, Target: false},
		},
	}
	got := flattenPathpointStages(stages)
	require.Len(t, got, 1)
	related, ok := got[0]["related"].([]map[string]interface{})
	require.True(t, ok)
	require.Len(t, related, 1)
	assert.True(t, related[0]["source"].(bool))
	assert.False(t, related[0]["target"].(bool))
}
