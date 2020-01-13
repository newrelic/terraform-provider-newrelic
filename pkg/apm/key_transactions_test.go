// +build unit

package apm

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	testListKeyTransactionsWithParamsResponseJSON = `{
		"key_transactions": [
			{
				"id": 456,
				"name": "test-key-transaction",
				"transaction_name": "test-key-transaction-name",
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
					"application": 54321
				}
			}
		]
	}`

	keyTransactionJSON = `{
		"key_transaction": {
			"id": 1,
			"name": "get /",
			"transaction_name": "get /",
			"health_status": "unknown",
			"reporting": true,
			"last_reported_at": "2020-01-02T21:56:07+00:00",
			"application_summary": {
				"response_time": 0.382,
				"throughput": 24,
				"error_rate": 0,
				"apdex_target": 0.01,
				"apdex_score": 1
			},
			"links": {
				"application": 12345
			}
		}
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

	actual, err := apm.ListKeyTransactions(nil)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestListKeyTransactionsWithParams(t *testing.T) {
	t.Parallel()
	keyTransactionName := "test-key-transaction"
	idsFilter := "123,456,789"

	apm := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		name := values.Get("filter[name]")

		require.Equal(t, keyTransactionName, name)

		ids := values.Get("filter[ids]")

		require.Equal(t, idsFilter, ids)

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(testListKeyTransactionsWithParamsResponseJSON))

		require.NoError(t, err)
	}))

	expected := []*KeyTransaction{
		{
			ID:              456,
			Name:            keyTransactionName,
			TransactionName: "test-key-transaction-name",
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
				Application: 54321,
			},
		},
	}

	params := ListKeyTransactionsParams{
		Name: keyTransactionName,
		IDs:  []int{123, 456, 789},
	}

	actual, err := apm.ListKeyTransactions(&params)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestGetKeyTransaction(t *testing.T) {
	t.Parallel()
	apm := newMockResponse(t, keyTransactionJSON, http.StatusOK)

	expected := &KeyTransaction{
		ID:              1,
		Name:            "get /",
		TransactionName: "get /",
		HealthStatus:    "unknown",
		LastReportedAt:  "2020-01-02T21:56:07+00:00",
		Reporting:       true,
		Summary: ApplicationSummary{
			ResponseTime: 0.382,
			Throughput:   24,
			ErrorRate:    0,
			ApdexTarget:  0.01,
			ApdexScore:   1,
		},
		Links: KeyTransactionLinks{
			Application: 12345,
		},
	}

	actual, err := apm.GetKeyTransaction(1)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
