package alerts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var incidentTestAPIHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Incidents []AlertIncident `json:"incidents,omitempty"`
	}{}

	openIncident := AlertIncident{
		ID:                 42,
		OpenedAt:           1575502560942,
		IncidentPreference: "PER_CONDITION",
		Links: AlertIncidentLink{
			Violations: []int{123456789},
			PolicyID:   12345,
		},
	}

	closedIncident := AlertIncident{
		ID:                 24,
		OpenedAt:           1575506284796,
		ClosedAt:           1575506342161,
		IncidentPreference: "PER_POLICY",
		Links: AlertIncidentLink{
			Violations: []int{987654321},
			PolicyID:   54321,
		},
	}

	// always including the open incident
	response.Incidents = append(response.Incidents, openIncident)

	// if not "only open", add the closed incident
	params := r.URL.Query()
	oo, ok := params["only_open"]
	fmt.Printf("Only Open: %+v\n", oo)
	if !ok || (ok && len(oo) > 0 && oo[0] != "true") {
		response.Incidents = append(response.Incidents, closedIncident)
	}

	// if "exclude violations", remove the violation links
	ev, ok := params["exclude_violations"]
	fmt.Printf("Exclude Violations: %+v\n", oo)
	if ok && len(ev) > 0 && ev[0] == "true" {
		for i := range response.Incidents {
			response.Incidents[i].Links.Violations = nil
		}
	}
	fmt.Printf("Incidents: %+v\n", response.Incidents)

	// set up response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	body, err := json.Marshal(response)
	if err != nil {
		panic(fmt.Errorf("error marshalling json: %w", err))
	}

	_, err = w.Write(body)
	if err != nil {
		panic(fmt.Errorf("failed to write test response body: %w", err))
	}
})

var failingTestHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)
})

func TestListAlertIncidents(t *testing.T) {
	c := newTestAlerts(incidentTestAPIHandler)

	expected := []AlertIncident{
		{
			ID:                 42,
			OpenedAt:           1575502560942,
			IncidentPreference: "PER_CONDITION",
			Links: AlertIncidentLink{
				Violations: []int{123456789},
				PolicyID:   12345,
			},
		},
		{
			ID:                 24,
			OpenedAt:           1575506284796,
			ClosedAt:           1575506342161,
			IncidentPreference: "PER_POLICY",
			Links: AlertIncidentLink{
				Violations: []int{987654321},
				PolicyID:   54321,
			},
		},
	}

	alertIncidents, err := c.ListAlertIncidents(false, false)
	if err != nil {
		t.Log(err)
		t.Fatal("GetAlertIncident error")
	}
	if alertIncidents == nil {
		t.Log(err)
		t.Fatal("GetAlertIncident error")
	}
	if diff := cmp.Diff(alertIncidents, expected); diff != "" {
		t.Fatalf("Alert incidents not parsed correctly: %s", diff)
	}
}

func TestOpenListAlertIncidents(t *testing.T) {
	c := newTestAlerts(incidentTestAPIHandler)

	expected := []AlertIncident{
		{
			ID:                 42,
			OpenedAt:           1575502560942,
			IncidentPreference: "PER_CONDITION",
			Links: AlertIncidentLink{
				Violations: []int{123456789},
				PolicyID:   12345,
			},
		},
	}

	alertIncidents, err := c.ListAlertIncidents(true, false)
	if err != nil {
		t.Log(err)
		t.Fatal("GetAlertIncident error")
	}
	if alertIncidents == nil {
		t.Log(err)
		t.Fatal("GetAlertIncident error")
	}
	if diff := cmp.Diff(alertIncidents, expected); diff != "" {
		t.Fatalf("Alert incidents not parsed correctly: %s", diff)
	}
}

func TestListAlertIncidentsWithoutViolations(t *testing.T) {
	c := newTestAlerts(incidentTestAPIHandler)

	expected := []AlertIncident{
		{
			ID:                 42,
			OpenedAt:           1575502560942,
			IncidentPreference: "PER_CONDITION",
			Links: AlertIncidentLink{
				PolicyID: 12345,
			},
		},
		{
			ID:                 24,
			OpenedAt:           1575506284796,
			ClosedAt:           1575506342161,
			IncidentPreference: "PER_POLICY",
			Links: AlertIncidentLink{
				PolicyID: 54321,
			},
		},
	}

	alertIncidents, err := c.ListAlertIncidents(false, true)
	if err != nil {
		t.Log(err)
		t.Fatal("GetAlertIncident error")
	}
	if alertIncidents == nil {
		t.Log(err)
		t.Fatal("GetAlertIncident error")
	}
	if diff := cmp.Diff(alertIncidents, expected); diff != "" {
		t.Fatalf("Alert incidents not parsed correctly: %s", diff)
	}
}

func TestListAlertIncidentFailing(t *testing.T) {
	c := newTestAlerts(failingTestHandler)

	_, err := c.ListAlertIncidents(false, false)
	if err == nil {
		t.Fatal("ListAlertIncident expected an error")
	}
}

func TestAcknowledgeAlertIncident(t *testing.T) {
	jsonResponse := `
			{
				"incidents": [
			    {
			      "id": 42,
				  "opened_at": 1575502560942,
			      "incident_preference": "PER_CONDITION",
			      "links": {
			        "violations": [
			          123456789
			        ],
			        "policy_id": 12345
				  }
				}
				]
			}
	`
	alerts := newMockResponse(t, jsonResponse, http.StatusOK)

	err := alerts.AcknowledgeAlertIncident(42)
	if err != nil {
		t.Log(err)
		t.Fatal("AckAlertIncident error")
	}
}

func TestAcknowledgeAlertIncidentFailing(t *testing.T) {
	c := newTestAlerts(failingTestHandler)

	err := c.CloseAlertIncident(42)
	if err == nil {
		t.Fatal("AckAlertIncident expected an error")
	}
}

func TestCloseAlertIncident(t *testing.T) {
	jsonResponse := `
		{
			"incidents": [
		    	{
		    	  "id": 42,
				  "opened_at": 1575502560942,
				  "closed_at": 1575502560943,
		    	  "incident_preference": "PER_CONDITION",
		    	  "links": {
		    	    "violations": [
		    	      123456789
		    	    ],
		    	    "policy_id": 12345
				  }
				}
			]
		}
	`

	alerts := newMockResponse(t, jsonResponse, http.StatusOK)

	err := alerts.AcknowledgeAlertIncident(42)
	if err != nil {
		t.Log(err)
		t.Fatal("CloseAlertIncident error")
	}
}

func TestCloseAlertIncidentFailing(t *testing.T) {
	c := newTestAlerts(failingTestHandler)

	err := c.CloseAlertIncident(42)
	if err == nil {
		t.Fatal("CloseAlertIncident expected an error")
	}
}
