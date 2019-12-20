package synthetics

// GetMonitorScript is used to retrieve the script that belongs
// to a New Relic Synthetics scripted monitor.
func (s *Synthetics) GetMonitorScript(monitorID string) (*MonitorScript, error) {
	return nil, nil
}

// UpdateMonitorScript is used to add a script to an existing New Relic Synthetics monitor.
func (s *Synthetics) UpdateMonitorScript(monitorID string, script MonitorScript) error {
	return nil
}
