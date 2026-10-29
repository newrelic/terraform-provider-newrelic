package newrelic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/servicelevel"
	"github.com/stretchr/testify/assert"
)

func TestExpandServiceLevelFacets(t *testing.T) {
	t.Parallel()

	assert.Equal(t, []string{"appName", "transactionType"}, expandServiceLevelFacets(map[string]interface{}{
		"facets": []interface{}{"appName", "transactionType"},
	}))

	assert.Nil(t, expandServiceLevelFacets(map[string]interface{}{
		"facets": []interface{}{},
	}))

	assert.Nil(t, expandServiceLevelFacets(map[string]interface{}{}))
}

func TestExpandServiceLevelEventsCreateInput_Facets(t *testing.T) {
	t.Parallel()

	cfg := map[string]interface{}{
		"account_id": 12345678,
		"facets":     []interface{}{"appName"},
		"valid_events": []interface{}{
			map[string]interface{}{
				"from": "Transaction",
			},
		},
	}

	got := expandServiceLevelEventsCreateInput(cfg)

	assert.Equal(t, []string{"appName"}, got.Facets)
}

func TestExpandServiceLevelEventsUpdateInput_Facets(t *testing.T) {
	t.Parallel()

	cfg := map[string]interface{}{
		"facets": []interface{}{"appName"},
		"valid_events": []interface{}{
			map[string]interface{}{
				"from": "Transaction",
			},
		},
	}

	got := expandServiceLevelEventsUpdateInput(cfg)

	assert.Equal(t, []string{"appName"}, got.Facets)
}

func TestFlattenServiceLevelIndicator_Facets(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, resourceNewRelicServiceLevel().Schema, map[string]interface{}{})

	indicator := servicelevel.ServiceLevelIndicator{
		Events: servicelevel.ServiceLevelEvents{
			ValidEvents: &servicelevel.ServiceLevelEventsQuery{
				From:   "Transaction",
				Facets: []string{"appName"},
				Select: servicelevel.ServiceLevelEventsQuerySelect{Function: "COUNT"},
			},
		},
	}

	identifier := &serviceLevelIdentifier{AccountID: 12345678, ID: "1", EntityGUID: "guid"}

	err := flattenServiceLevelIndicator(indicator, identifier, d, "sli-guid")
	assert.NoError(t, err)

	events := d.Get("events").([]interface{})[0].(map[string]interface{})
	assert.Equal(t, []interface{}{"appName"}, events["facets"])
}

func TestFlattenServiceLevelIndicator_NoFacets(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, resourceNewRelicServiceLevel().Schema, map[string]interface{}{})

	indicator := servicelevel.ServiceLevelIndicator{
		Events: servicelevel.ServiceLevelEvents{
			ValidEvents: &servicelevel.ServiceLevelEventsQuery{
				From:   "Transaction",
				Select: servicelevel.ServiceLevelEventsQuerySelect{Function: "COUNT"},
			},
		},
	}

	identifier := &serviceLevelIdentifier{AccountID: 12345678, ID: "1", EntityGUID: "guid"}

	err := flattenServiceLevelIndicator(indicator, identifier, d, "sli-guid")
	assert.NoError(t, err)

	events := d.Get("events").([]interface{})[0].(map[string]interface{})
	assert.Empty(t, events["facets"])
}
