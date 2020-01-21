package newrelic

import (
	newrelic "github.com/paultyng/go-newrelic/v4/api"
)

func flattenNrql(nrql newrelic.AlertNrqlQuery) []interface{} {
	m := map[string]interface{}{
		"query":       nrql.Query,
		"since_value": nrql.SinceValue,
	}

	return []interface{}{m}
}
