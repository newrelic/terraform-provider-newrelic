package synthetics

import (
	"fmt"
	"path"
)

const (
	listMonitorsLimit = 100
)

// ListMonitors is used to retrieve New Relic Synthetics monitors.
func (s *Synthetics) ListMonitors() ([]Monitor, error) {
	resp := listMonitorsResponse{}

	queryParams := listMonitorsParams{
		Limit: listMonitorsLimit,
	}

	_, err := s.client.Get("/monitors", &queryParams, &resp)

	if err != nil {
		return nil, err
	}

	return resp.Monitors, nil
}

// GetMonitor is used to retrieve a specific New Relic Synthetics monitor.
func (s *Synthetics) GetMonitor(monitorID string) (*Monitor, error) {
	resp := Monitor{}
	url := fmt.Sprintf("/monitors/%s", monitorID)
	_, err := s.client.Get(url, nil, &resp)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// CreateMonitor is used to create a New Relic Synthetics monitor.
// If successful it returns the ID of the created resource.
func (s *Synthetics) CreateMonitor(monitor Monitor) (string, error) {
	resp, err := s.client.Post("/monitors", nil, &monitor, nil)

	if err != nil {
		return "", err
	}

	l := resp.Header.Get("location")
	monitorID := path.Base(l)

	return monitorID, nil
}

// UpdateMonitor is used to create a New Relic Synthetics monitor.
func (s *Synthetics) UpdateMonitor(monitor Monitor) error {
	url := fmt.Sprintf("/monitors/%s", monitor.ID)
	_, err := s.client.Put(url, nil, &monitor, nil)

	if err != nil {
		return err
	}

	return nil
}

// DeleteMonitor is used to create a New Relic Synthetics monitor.
func (s *Synthetics) DeleteMonitor(monitorID string) error {
	url := fmt.Sprintf("/monitors/%s", monitorID)
	_, err := s.client.Delete(url, nil, nil)

	if err != nil {
		return err
	}

	return nil
}

type listMonitorsResponse struct {
	Monitors []Monitor `json:"monitors,omitempty"`
}

type listMonitorsParams struct {
	Limit int `url:"limit,omitempty"`
}
