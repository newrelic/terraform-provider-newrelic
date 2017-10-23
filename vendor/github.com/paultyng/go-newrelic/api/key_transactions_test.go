package api

import (
	"net/http"
	"testing"
)

func TestKeyTransactions_Basic(t *testing.T) {
	c := newTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
      {
				"key_transactions": [{
					"id": 1,
					"name": "foo",
					"transaction_name": "/bar",
					"health_status": "unknown",
					"reporting": true,
					"last_reported_at": "2017-10-19T18:16:08+00:00",
					"application_summary": {
						"response_time": 0.0,
						"throughput": 0.0,
						"error_rate": 0,
						"apdex_target": 0.5,
						"apdex_score": 0.0
					},
					"links": {
						"application": 2
					}
				}],
				"links": {
					"key_transaction.application": "/v2/applications/{application_id}"
				}
			}
    `))
	}))

	apps, err := c.queryKeyTransactions()
	if err != nil {
		t.Log(err)
		t.Fatal("queryKeyTransactions error")
	}

	if len(apps) == 0 {
		t.Fatal("No applications found")
	}
}
