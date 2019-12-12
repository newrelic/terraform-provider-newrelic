package synthetics

type listMonitorsResponse struct {
	Monitors []Monitor `json:"monitors,omitempty"`
}

// ListMonitors is used to retrieve New Relic Synthetics monitors.
func (s *Synthetics) ListMonitors() ([]Monitor, error) {
	res := listMonitorsResponse{}
	err := s.client.Get("monitors", nil, &res)

	if err != nil {
		return nil, err
	}

	return res.Monitors, nil
}
