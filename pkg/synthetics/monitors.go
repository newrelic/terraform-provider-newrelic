package synthetics

import "strconv"

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
func (s *Synthetics) GetMonitor(monitorID int) ([]Monitor, error) {
	return nil, nil
}
