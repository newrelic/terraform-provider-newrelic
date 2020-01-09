// +build unit

package alerts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var incidentTestAPIHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Incidents []Incident `json:"incidents,omitempty"`
	}{}

	openIncident := Incident{
		ID:                 42,
		OpenedAt:           1575502560942,
		IncidentPreference: "PER_CONDITION",
		Links: IncidentLink{
			Violations: []int{123456789},
			PolicyID:   12345,
		},
	}

	closedIncident := Incident{
		ID:                 24,
		OpenedAt:           1575506284796,
		ClosedAt:           1575506342161,
		IncidentPreference: "PER_POLICY",
		Links: IncidentLink{
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

func TestListIncidents(t *testing.T) {
	t.Parallel()

	c := newTestAlerts(incidentTestAPIHandler)

	expected := []*Incident{
		{
			ID:                 42,
			OpenedAt:           1575502560942,
			IncidentPreference: "PER_CONDITION",
			Links: IncidentLink{
				Violations: []int{123456789},
				PolicyID:   12345,
			},
		},
		{
			ID:                 24,
			OpenedAt:           1575506284796,
			ClosedAt:           1575506342161,
			IncidentPreference: "PER_POLICY",
			Links: IncidentLink{
				Violations: []int{987654321},
				PolicyID:   54321,
			},
		},
	}

	alertIncidents, err := c.ListIncidents(false, false)

	assert.NoError(t, err)
	assert.NotNil(t, alertIncidents)
	assert.Equal(t, expected, alertIncidents)
}

func TestOpenListIncidents(t *testing.T) {
	t.Parallel()

	c := newTestAlerts(incidentTestAPIHandler)

	expected := []*Incident{
		{
			ID:                 42,
			OpenedAt:           1575502560942,
			IncidentPreference: "PER_CONDITION",
			Links: IncidentLink{
				Violations: []int{123456789},
				PolicyID:   12345,
			},
		},
	}

	alertIncidents, err := c.ListIncidents(true, false)

	assert.NoError(t, err)
	assert.NotNil(t, alertIncidents)
	assert.Equal(t, expected, alertIncidents)
}

func TestListIncidentsWithoutViolations(t *testing.T) {
	t.Parallel()

	c := newTestAlerts(incidentTestAPIHandler)

	expected := []*Incident{
		{
			ID:                 42,
			OpenedAt:           1575502560942,
			IncidentPreference: "PER_CONDITION",
			Links: IncidentLink{
				PolicyID: 12345,
			},
		},
		{
			ID:                 24,
			OpenedAt:           1575506284796,
			ClosedAt:           1575506342161,
			IncidentPreference: "PER_POLICY",
			Links: IncidentLink{
				PolicyID: 54321,
			},
		},
	}

	alertIncidents, err := c.ListIncidents(false, true)

	assert.NoError(t, err)
	assert.NotNil(t, alertIncidents)
	assert.Equal(t, expected, alertIncidents)
}

func TestListIncidentFailing(t *testing.T) {
	t.Parallel()

	c := newTestAlerts(failingTestHandler)

	_, err := c.ListIncidents(false, false)

	assert.Error(t, err)
}

func TestAcknowledgeIncident(t *testing.T) {
	t.Parallel()

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

	_, err := alerts.AcknowledgeIncident(42)

	assert.NoError(t, err)
}

func TestAcknowledgeIncidentFailing(t *testing.T) {
	t.Parallel()

	c := newTestAlerts(failingTestHandler)

	_, err := c.CloseIncident(42)

	assert.Error(t, err)
}

func TestCloseIncident(t *testing.T) {
	t.Parallel()

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

	_, err := alerts.AcknowledgeIncident(42)
	if err != nil {
		t.Log(err)
		t.Fatal("CloseIncident error")
	}
}

func TestCloseIncidentFailing(t *testing.T) {
	t.Parallel()

	c := newTestAlerts(failingTestHandler)

	_, err := c.CloseIncident(42)
	if err == nil {
		t.Fatal("CloseIncident expected an error")
	}
}
