package api

import (
	"fmt"
	"net/http"
	"testing"
)

func TestQueryAlertSyntheticsConditions(t *testing.T) {
	c := newTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
			  "synthetics_conditions": [
			    {
			      "id": 12345,
			      "name": "Synthetics Condition",
			      "runbook_url": "https://example.com/runbook.md",
						"monitor_id": "12345678-1234-1234-1234-1234567890ab",
			      "enabled": true
			    }
			  ]
			}
			`))
	}))

	policyID := 123

	SyntheticsAlertConditions, err := c.queryAlertSyntheticsConditions(policyID)
	if err != nil {
		t.Log(err)
		t.Fatal("queryAlertSyntheticsConditions error")
	}

	if len(SyntheticsAlertConditions) == 0 {
		t.Fatal("No Synthetics Alert Conditions found")
	}
	if SyntheticsAlertConditions[0].MonitorID != "12345678-1234-1234-1234-1234567890ab" {
		t.Fatal("MonitorID was not parsed properly")
	}
}

func TestGetAlertSyntheticsCondition(t *testing.T) {
	c := newTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
			  "synthetics_conditions": [
			    {
			      "id": 12345,
			      "name": "Synthetics Condition",
			      "runbook_url": "https://example.com/runbook.md",
						"monitor_id": "12345678-1234-1234-1234-1234567890ab",
			      "enabled": true
			    }
			  ]
			}
			`))
	}))

	policyID := 123
	conditionID := 12345

	SyntheticsAlertCondition, err := c.GetAlertSyntheticsCondition(policyID, conditionID)
	if err != nil {
		t.Log(err)
		t.Fatal("GetAlertSyntheticsCondition error")
	}
	if SyntheticsAlertCondition == nil {
		t.Log(err)
		t.Fatal("GetAlertSyntheticsCondition error")
	}
	if SyntheticsAlertCondition.MonitorID != "12345678-1234-1234-1234-1234567890ab" {
		t.Fatal("MonitorID was not parsed properly")
	}
}

func TestListAlertSyntheticsConditions(t *testing.T) {
	c := newTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
			  "synthetics_conditions": [
			    {
			      "id": 12345,
			      "name": "Synthetics Condition",
			      "runbook_url": "https://example.com/runbook.md",
						"monitor_id": "12345678-1234-1234-1234-1234567890ab",
			      "enabled": true
			    }
			  ]
			}
			`))
	}))

	policyID := 123

	SyntheticsAlertConditions, err := c.ListAlertSyntheticsConditions(policyID)
	if err != nil {
		t.Log(err)
		t.Fatal("ListAlertSyntheticsConditions error")
	}
	if len(SyntheticsAlertConditions) == 0 {
		t.Log(err)
		t.Fatal("ListAlertSyntheticsConditions error")
	}
	if SyntheticsAlertConditions[0].MonitorID != "12345678-1234-1234-1234-1234567890ab" {
		t.Fatal("MonitorID was not parsed properly")
	}
}

func TestCreateAlertSyntheticsCondition(t *testing.T) {
	c := newTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
		  {
		    "synthetics_condition": {
		      "id": 12345,
		      "name": "Synthetics Condition",
		      "runbook_url": "https://example.com/runbook.md",
		      "monitor_id": "12345678-1234-1234-1234-1234567890ab",
		      "enabled": true
		    }
		  }
		  `))
	}))

	SyntheticsAlertCondition := AlertSyntheticsCondition{
		PolicyID:   123,
		Name:       "Synthetics Condition",
		Enabled:    true,
		RunbookURL: "https://example.com/runbook.md",
		MonitorID:  "12345678-1234-1234-1234-1234567890ab",
	}

	SyntheticsAlertConditionResp, err := c.CreateAlertSyntheticsCondition(SyntheticsAlertCondition)
	if err != nil {
		t.Log(err)
		t.Fatal("CreateAlertSyntheticsCondition error")
	}
	if SyntheticsAlertConditionResp == nil {
		t.Log(err)
		t.Fatal("CreateAlertSyntheticsCondition error")
	}
	if SyntheticsAlertConditionResp.MonitorID != "12345678-1234-1234-1234-1234567890ab" {
		t.Log(SyntheticsAlertConditionResp.MonitorID)
		t.Fatal("MonitorID was not parsed correctly")
	}
}

func TestUpdateAlertSyntheticsCondition(t *testing.T) {
	c := newTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
		  {
		    "synthetics_condition": {
		      "id": 12345,
		      "name": "Synthetics Condition",
		      "runbook_url": "https://example.com/runbook.md",
		      "monitor_id": "12345678-1234-1234-1234-1234567890ab",
		      "enabled": true
		    }
		  }
		  `))
	}))

	SyntheticsAlertCondition := AlertSyntheticsCondition{
		PolicyID:   123,
		Name:       "Synthetics Condition",
		Enabled:    true,
		RunbookURL: "https://example.com/runbook.md",
		MonitorID:  "12345678-1234-1234-1234-1234567890ab",
	}

	SyntheticsAlertConditionResp, err := c.UpdateAlertSyntheticsCondition(SyntheticsAlertCondition)
	if err != nil {
		t.Log(err)
		t.Fatal("UpdateAlertSyntheticsCondition error")
	}
	if SyntheticsAlertConditionResp == nil {
		t.Log(err)
		t.Fatal("UpdateAlertSyntheticsCondition error")
	}
	if SyntheticsAlertConditionResp.MonitorID != "12345678-1234-1234-1234-1234567890ab" {
		t.Log(SyntheticsAlertConditionResp.MonitorID)
		t.Fatal("MonitorID was not parsed correctly")
	}
}

func TestDeleteAlertSyntheticsCondition(t *testing.T) {
	policyID := 123
	conditionID := 12345
	c := newTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if r.Method != "DELETE" {
			t.Fatal("DeleteAlertSyntheticsCondition did not use DELETE method")
		}
		if r.URL.Path != fmt.Sprintf("/alerts_synthetics_conditions/%v.json", conditionID) {
			t.Fatal("DeleteAlertSyntheticsCondtion did not use the correct URL")
		}
	}))
	err := c.DeleteAlertSyntheticsCondition(policyID, conditionID)
	if err != nil {
		t.Log(err)
		t.Fatal("DeleteAlertNrqlCondition error")
	}
}
