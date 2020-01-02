// +build unit

package apm

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testListKeyTransactionsResponseJSON = `{
		"key_transactions": [
			{
				"id": 1,
				"name": "get /",
				"transaction_name": "get /",
				"health_status": "unknown",
				"reporting": true,
				"last_reported_at": "2020-01-02T21:09:07+00:00",
				"application_summary": {
					"response_time": 0.381,
					"throughput": 24,
					"error_rate": 0,
					"apdex_target": 0.01,
					"apdex_score": 1
				},
				"links": {
					"application": 12345
				}
			}
		]
	}`
)

func TestListKeyTransactions(t *testing.T) {
	t.Parallel()
	apm := newMockResponse(t, testListKeyTransactionsResponseJSON, http.StatusOK)

	expected := []*KeyTransaction{
		{
			ID:              1,
			Name:            "get /",
			TransactionName: "get /",
			HealthStatus:    "unknown",
			LastReportedAt:  "2020-01-02T21:09:07+00:00",
			Reporting:       true,
			Summary: ApplicationSummary{
				ResponseTime: 0.381,
				Throughput:   24,
				ErrorRate:    0,
				ApdexTarget:  0.01,
				ApdexScore:   1,
			},
			Links: KeyTransactionLinks{
				Application: 12345,
			},
		},
	}

	actual, err := apm.ListKeyTransactions()

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
