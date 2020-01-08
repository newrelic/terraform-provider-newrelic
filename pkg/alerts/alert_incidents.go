package alerts

import (
	"fmt"
)

// ListAlertIncidents returns all alert incidents
func (alerts *Alerts) ListAlertIncidents(onlyOpen bool, excludeViolations bool) ([]AlertIncident, error) {
	incidentsResponse := alertIncidentsResponse{}
	incidents := []AlertIncident{}
	queryParams := listAlertIncidentsParams{
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

// AcknowledgeAlertIncident acknowledges an existing incident
func (alerts *Alerts) AcknowledgeAlertIncident(id int) error {
	return alerts.updateAlertIncident(id, "acknowledge")
}

// CloseAlertIncident closes an existing open incident
func (alerts *Alerts) CloseAlertIncident(id int) error {
	return alerts.updateAlertIncident(id, "close")
}

func (alerts *Alerts) updateAlertIncident(id int, verb string) error {
	path := fmt.Sprintf("/alerts_incidents/%v/%v.json", id, verb)
	_, err := alerts.client.Put(path, nil, nil, nil)
	return err
}

type listAlertIncidentsParams struct {
	OnlyOpen          bool `url:"only_open,omitempty"`
	ExcludeViolations bool `url:"exclude_violations,omitempty"`
}

type alertIncidentsResponse struct {
	Incidents []AlertIncident `json:"incidents,omitempty"`
}
