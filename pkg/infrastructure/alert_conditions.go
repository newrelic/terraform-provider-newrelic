package infrastructure

type listAlertConditionsResponse struct {
	AlertConditions []AlertCondition `json:"data,omitempty"`
}

// ListAlertConditions is used to retrieve New Relic Infrastructure alert conditions.
func (s *Infrastructure) ListAlertConditions() ([]AlertCondition, error) {
	return nil, nil
}
