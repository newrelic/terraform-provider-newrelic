package apm

// ListComponentsParams represents a set of filters to be
// used when querying New Relic applications.
type ListComponentsParams struct {
	Name         string
	IDs          []int
	PluginID     int
	HealthStatus string
}

// ListComponents is used to retrieve the components associated with
// a New Relic account.
func (a *APM) ListComponents(params *ListComponentsParams) (*[]Component, error) {
	return nil, nil
}

// ShowComponents is used to retrieve a specific New Relic component.
func (a *APM) GetComponent(componentID int) (*Component, error) {
	return nil, nil
}
