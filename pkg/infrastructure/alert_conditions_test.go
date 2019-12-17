// +build unit

package infrastructure

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/newrelic/newrelic-client-go/internal/serialization"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func NewTestInfrastructure(handler http.Handler) Infrastructure {
	ts := httptest.NewServer(handler)

	c := New(config.ReplacementConfig{
		APIKey:    "abc123",
		BaseURL:   ts.URL,
		UserAgent: "newrelic/newrelic-client-go",
	})

	return c
}

func TestListAlertConditions(t *testing.T) {
	t.Parallel()
	infra := NewTestInfrastructure(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`
		{
			"data":[
			   {
					"type":"infra_process_running",
					"name":"Java is running",
					"enabled":true,
					"where_clause":"(hostname LIKE '%cassandra%')",
					"id":13890,
					"created_at_epoch_millis":1490996713872,
					"updated_at_epoch_millis":1490996713872,
					"policy_id":111111,
					"comparison":"equal",
					"critical_threshold":{
						"value":0,
						"duration_minutes":6
					},
					"process_where_clause":"(commandName = 'java')"
			   }
			]
		 }
		`))

		if err != nil {
			t.Fatal(err)
		}
	}))

	critical := Threshold{
		Value:    0,
		Duration: 6,
	}

	expected := []AlertCondition{
		{
			Type:         "infra_process_running",
			Name:         "Java is running",
			Enabled:      true,
			Where:        "(hostname LIKE '%cassandra%')",
			ID:           13890,
			CreatedAt:    serialization.Epoch(time.Unix(1490996713872, 0)),
			UpdatedAt:    serialization.Epoch(time.Unix(1490996713872, 0)),
			PolicyID:     111111,
			Comparison:   "equal",
			Critical:     &critical,
			ProcessWhere: "(commandName = 'java')",
		},
	}

	actual, err := infra.ListAlertConditions(111111)

	if err != nil {
		t.Fatalf("ListAlertConditions error: %s", err)
	}

	if actual == nil {
		t.Fatalf("ListAlertConditions response is nil")
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("ListAlertConditions response differs from expected: %s", diff)
	}
}
