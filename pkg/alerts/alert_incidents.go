package alerts

import (
	"fmt"
)

// ListIncidents returns all alert incidents.
func (alerts *Alerts) ListIncidents(onlyOpen bool, excludeViolations bool) ([]*Incident, error) {
	incidentsResponse := alertIncidentsResponse{}
	incidents := []*Incident{}
	queryParams := listIncidentsParams{
		OnlyOpen:          onlyOpen,
		ExcludeViolations: excludeViolations,
	}

	nextURL := "/alerts_incidents.json"

	for nextURL != "" {
		resp, err := alerts.client.Get(nextURL, queryParams, &incidentsResponse)

		if err != nil {
			return nil, err
		}

		incidents = append(incidents, incidentsResponse.Incidents...)

		paging := alerts.pager.Parse(resp)
		nextURL = paging.Next
	}

	return incidents, nil
}

// AcknowledgeIncident acknowledges an existing incident.
func (alerts *Alerts) AcknowledgeIncident(id int) (*Incident, error) {
	return alerts.updateIncident(id, "acknowledge")
}

// CloseIncident closes an existing open incident.
func (alerts *Alerts) CloseIncident(id int) (*Incident, error) {
	return alerts.updateIncident(id, "close")
}

func (alerts *Alerts) updateIncident(id int, verb string) (*Incident, error) {
	response := alertIncidentResponse{}
	path := fmt.Sprintf("/alerts_incidents/%v/%v.json", id, verb)
	_, err := alerts.client.Put(path, nil, nil, &response)

	if err != nil {
		return nil, err
	}

	return &response.Incident, nil
}

type listIncidentsParams struct {
	OnlyOpen          bool `url:"only_open,omitempty"`
	ExcludeViolations bool `url:"exclude_violations,omitempty"`
}

type alertIncidentsResponse struct {
	Incidents []*Incident `json:"incidents,omitempty"`
}

type alertIncidentResponse struct {
	Incident Incident `json:"incident,omitempty"`
}
