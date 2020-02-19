package synthetics

import (
	"fmt"
	"path"
)

const (
	listMonitorsLimit = 100
)

// ListMonitors is used to retrieve New Relic Synthetics monitors.
func (s *Synthetics) ListMonitors() ([]*Monitor, error) {
	resp := listMonitorsResponse{}
	queryParams := listMonitorsParams{
		Limit: listMonitorsLimit,
	}

	_, err := s.client.Get("/v4/monitors", &queryParams, &resp)

	if err != nil {
		return nil, err
	}

	return resp.Monitors, nil
}

// GetMonitor is used to retrieve a specific New Relic Synthetics monitor.
func (s *Synthetics) GetMonitor(monitorID string) (*Monitor, error) {
	resp := Monitor{}
	url := fmt.Sprintf("/v4/monitors/%s", monitorID)
	_, err := s.client.Get(url, nil, &resp)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// CreateMonitor is used to create a New Relic Synthetics monitor.
func (s *Synthetics) CreateMonitor(monitor Monitor) (*Monitor, error) {
	resp, err := s.client.Post("/v4/monitors", nil, &monitor, nil)

	if err != nil {
		return nil, err
	}

	l := resp.Header.Get("location")
	monitorID := path.Base(l)

	monitor.ID = monitorID

	return &monitor, nil
}

// UpdateMonitor is used to update a New Relic Synthetics monitor.
func (s *Synthetics) UpdateMonitor(monitor Monitor) (*Monitor, error) {
	url := fmt.Sprintf("/v4/monitors/%s", monitor.ID)
	_, err := s.client.Put(url, nil, &monitor, nil)

	if err != nil {
		return nil, err
	}

	return &monitor, nil
}

// DeleteMonitor is used to delete a New Relic Synthetics monitor.
func (s *Synthetics) DeleteMonitor(monitorID string) error {
	url := fmt.Sprintf("/v4/monitors/%s", monitorID)
	_, err := s.client.Delete(url, nil, nil)

	if err != nil {
		return err
	}

	return nil
}

type listMonitorsResponse struct {
	Monitors []*Monitor `json:"monitors,omitempty"`
}

type listMonitorsParams struct {
	Limit int `url:"limit,omitempty"`
}
