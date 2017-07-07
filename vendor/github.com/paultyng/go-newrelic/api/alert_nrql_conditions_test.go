package api

import (
	"fmt"
	"net/http"
	"testing"
)

func TestQueryAlertNrqlConditions(t *testing.T) {
	c := newTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
			  "nrql_conditions": [
			    {
			      "id": 12345,
			      "name": "NRQL Condition",
			      "runbook_url": "https://example.com/runbook.md",
			      "enabled": true,
			      "terms": [
			        {
			          "duration": "10",
			          "operator": "below",
			          "priority": "critical",
			          "threshold": "2",
			          "time_function": "all"
			         }
			      ],
			      "value_function": "single_value",
			      "nrql": {
			        "query": "SELECT uniqueCount(fieldname) FROM indexname WHERE fieldname2 = 'somevaluetofilterby'",
			        "since_value": "5"
			      }
			    }
			  ]
			}
			`))
	}))

	policyID := 123

	nrqlAlertConditions, err := c.queryAlertNrqlConditions(policyID)
	if err != nil {
		t.Log(err)
		t.Fatal("queryAlertNrqlConditions error")
	}

	if len(nrqlAlertConditions) == 0 {
		t.Fatal("No NRQL Alert Conditions found")
	}
}

func TestGetAlertNrqlCondition(t *testing.T) {
	c := newTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
			  "nrql_conditions": [
			    {
			      "id": 12345,
			      "name": "NRQL Condition",
			      "runbook_url": "https://example.com/runbook.md",
			      "enabled": true,
			      "terms": [
			        {
			          "duration": "10",
			          "operator": "below",
			          "priority": "critical",
			          "threshold": "2",
			          "time_function": "all"
			         }
			      ],
			      "value_function": "single_value",
			      "nrql": {
			        "query": "SELECT uniqueCount(fieldname) FROM indexname WHERE fieldname2 = 'somevaluetofilterby'",
			        "since_value": "5"
			      }
			    }
			  ]
			}
			`))
	}))

	policyID := 123
	conditionID := 12345

	nrqlAlertCondition, err := c.GetAlertNrqlCondition(policyID, conditionID)
	if err != nil {
		t.Log(err)
		t.Fatal("GetAlertNrqlCondition error")
	}
	if nrqlAlertCondition == nil {
		t.Log(err)
		t.Fatal("GetAlertNrqlCondition error")
	}
}

func TestListAlertNrqlConditions(t *testing.T) {
	c := newTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
			  "nrql_conditions": [
			    {
			      "id": 12345,
			      "name": "NRQL Condition",
			      "runbook_url": "https://example.com/runbook.md",
			      "enabled": true,
			      "terms": [
			        {
			          "duration": "10",
			          "operator": "below",
			          "priority": "critical",
			          "threshold": "2",
			          "time_function": "all"
			         }
			      ],
			      "value_function": "single_value",
			      "nrql": {
			        "query": "SELECT uniqueCount(fieldname) FROM indexname WHERE fieldname2 = 'somevaluetofilterby'",
			        "since_value": "5"
			      }
			    }
			  ]
			}
			`))
	}))

	policyID := 123

	nrqlAlertConditions, err := c.ListAlertNrqlConditions(policyID)
	if err != nil {
		t.Log(err)
		t.Fatal("ListAlertNrqlConditions error")
	}
	if len(nrqlAlertConditions) == 0 {
		t.Log(err)
		t.Fatal("ListAlertNrqlConditions error")
	}
}

func TestCreateAlertNrqlCondition(t *testing.T) {
	c := newTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
			  "nrql_condition":
			    {
			      "id": 12345,
			      "name": "NRQL Condition",
			      "runbook_url": "https://example.com/runbook.md",
			      "enabled": true,
			      "terms": [
			        {
			          "duration": "10",
			          "operator": "below",
			          "priority": "critical",
			          "threshold": "2",
			          "time_function": "all"
			         }
			      ],
			      "value_function": "single_value",
			      "nrql": {
			        "query": "SELECT uniqueCount(fieldname) FROM indexname WHERE fieldname2 = 'somevaluetofilterby'",
			        "since_value": "5"
			      }
			    }
			}
			`))
	}))

	nrqlAlertConditionTerms := []AlertConditionTerm{
		{
			Duration:     10,
			Operator:     "below",
			Priority:     "critical",
			Threshold:    2.0,
			TimeFunction: "all",
		},
	}

	nrqlAlertQuery := AlertNrqlQuery{
		Query:      "SELECT uniqueCount(fieldname) FROM indexname WHERE fieldname2 = 'somevaluetofilterby'",
		SinceValue: "5",
	}

	nrqlAlertCondition := AlertNrqlCondition{
		PolicyID:      123,
		Name:          "Test Condition",
		Enabled:       true,
		RunbookURL:    "https://example.com/runbook.md",
		Terms:         nrqlAlertConditionTerms,
		ValueFunction: "all",
		Nrql:          nrqlAlertQuery,
	}

	nrqlAlertConditionResp, err := c.CreateAlertNrqlCondition(nrqlAlertCondition)
	if err != nil {
		t.Log(err)
		t.Fatal("CreateAlertNrqlCondition error")
	}
	if nrqlAlertConditionResp == nil {
		t.Log(err)
		t.Fatal("CreateAlertNrqlCondition error")
	}
	if nrqlAlertConditionResp.ID != 12345 {
		t.Fatal("Condition ID was not parsed correctly")
	}
}

func TestUpdateAlertNrqlCondition(t *testing.T) {
	c := newTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
			  "nrql_condition":
			    {
			      "id": 12345,
			      "name": "NRQL Condition",
			      "runbook_url": "https://example.com/runbook.md",
			      "enabled": true,
			      "terms": [
			        {
			          "duration": "10",
			          "operator": "below",
			          "priority": "critical",
			          "threshold": "2",
			          "time_function": "all"
			         }
			      ],
			      "value_function": "single_value",
			      "nrql": {
			        "query": "SELECT uniqueCount(fieldname) FROM indexname WHERE fieldname2 = 'somevaluetofilterby'",
			        "since_value": "5"
			      }
			    }
			}
			`))
	}))

	nrqlAlertConditionTerms := []AlertConditionTerm{
		{
			Duration:     10,
			Operator:     "below",
			Priority:     "critical",
			Threshold:    2.0,
			TimeFunction: "all",
		},
	}

	nrqlAlertQuery := AlertNrqlQuery{
		Query:      "SELECT uniqueCount(fieldname) FROM indexname WHERE fieldname2 = 'somevaluetofilterby'",
		SinceValue: "5",
	}

	nrqlAlertCondition := AlertNrqlCondition{
		PolicyID:      123,
		Name:          "Test Condition",
		Enabled:       true,
		RunbookURL:    "https://example.com/runbook.md",
		Terms:         nrqlAlertConditionTerms,
		ValueFunction: "all",
		Nrql:          nrqlAlertQuery,
	}

	nrqlAlertConditionResp, err := c.UpdateAlertNrqlCondition(nrqlAlertCondition)
	if err != nil {
		t.Log(err)
		t.Fatal("UpdateAlertNrqlCondition error")
	}
	if nrqlAlertConditionResp == nil {
		t.Log(err)
		t.Fatal("UpdateAlertNrqlCondition error")
	}
	if nrqlAlertConditionResp.ID != 12345 {
		t.Fatal("Condition ID was not parsed correctly")
	}
}

func TestDeleteAlertNrqlCondition(t *testing.T) {
	policyID := 123
	conditionID := 12345
	c := newTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if r.Method != "DELETE" {
			t.Fatal("DeleteAlertNrqlCondition did not use DELETE method")
		}
		if r.URL.Path != fmt.Sprintf("/alerts_nrql_conditions/%v.json", conditionID) {
			t.Fatal("DeleteAlertNrqlCondtion did not use the correct URL")
		}
	}))
	err := c.DeleteAlertNrqlCondition(policyID, conditionID)
	if err != nil {
		t.Log(err)
		t.Fatal("DeleteAlertNrqlCondition error")
	}
}
