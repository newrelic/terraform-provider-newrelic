package newrelic

import (
	"reflect"
	"testing"

	newrelic "github.com/paultyng/go-newrelic/v4/api"
)

func TestFlattenNrql(t *testing.T) {
	expanded := newrelic.AlertNrqlQuery{
		Query:      "SELECT percentile(duration, 99) FROM Transaction FACET remote_ip",
		SinceValue: "3",
	}

	flattened := []interface{}{map[string]interface{}{
		"query":       "SELECT percentile(duration, 99) FROM Transaction FACET remote_ip",
		"since_value": "3",
	}}

	result := flattenNrql(expanded)

	if result == nil {
		t.Fatal("result was nil")
	}

	if !reflect.DeepEqual(result, flattened) {
		t.Fatalf("result %s not equal to expected %s", result, flattened)
	}
}
