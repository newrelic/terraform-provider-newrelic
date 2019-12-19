package synthetics

import (
	"fmt"
	"strconv"
)

const (
	listMonitorsLimit = 100
)

type listMonitorsResponse struct {
	Monitors []Monitor `json:"monitors,omitempty"`
}

// ListMonitors is used to retrieve New Relic Synthetics monitors.
func (s *Synthetics) ListMonitors() ([]Monitor, error) {
	res := listMonitorsResponse{}
	paramsMap := map[string]string{
		"limit": strconv.Itoa(listMonitorsLimit),
	}

	_, err := s.client.Get("/monitors", &paramsMap, &res)

	if err != nil {
		return nil, err
	}

	return res.Monitors, nil
}

// GetMonitor is used to retrieve a specific New Relic Synthetics monitor.
func (s *Synthetics) GetMonitor(monitorID string) (*Monitor, error) {
	res := Monitor{}
	url := fmt.Sprintf("/monitors/%s", monitorID)
	_, err := s.client.Get(url, nil, &res)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// CreateMonitor is used to create a New Relic Synthetics monitor.
func (s *Synthetics) CreateMonitor(monitor Monitor) (*Monitor, error) {
	res := Monitor{}
	_, err := s.client.Post("/monitors", nil, &monitor, &res)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// UpdateMonitor is used to create a New Relic Synthetics monitor.
func (s *Synthetics) UpdateMonitor(monitor Monitor) (*Monitor, error) {
	res := Monitor{}
	url := fmt.Sprintf("/monitors/%s", monitor.ID)
	_, err := s.client.Put(url, nil, &monitor, &res)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// DeleteMonitor is used to create a New Relic Synthetics monitor.
func (s *Synthetics) DeleteMonitor(monitorID string) (*Monitor, error) {
	return nil, nil
}
