package apm

import "fmt"

// ListLabels returns the labels for the account.
func (apm *APM) ListLabels() ([]*Label, error) {
	response := labelsResponse{}
	labels := []*Label{}
	nextURL := fmt.Sprintf("/labels.json")

	for nextURL != "" {
		resp, err := apm.client.Get(nextURL, nil, &response)

		if err != nil {
			return nil, err
		}

		labels = append(labels, response.Labels...)

		paging := apm.pager.Parse(resp)
		nextURL = paging.Next
	}

	return labels, nil
}

// GetLabel gets the label for the specified key.
func (apm *APM) GetLabel(key string) (*Label, error) {
	labels, err := apm.ListLabels()

	if err != nil {
		return nil, err
	}

	for _, label := range labels {
		if label.Key == key {
			return label, nil
		}
	}

	return nil, fmt.Errorf("no label found with key %s", key)
}

// CreateLabel creates a new label for an account.
func (apm *APM) CreateLabel(label Label) (*Label, error) {
	reqBody := labelRequestBody{
		Label: label,
	}
	resp := labelResponse{}

	u := fmt.Sprintf("/labels.json")

	// The API currently uses a PUT request for label creation
	_, err := apm.client.Put(u, nil, &reqBody, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Label, nil
}

type labelsResponse struct {
	Labels []*Label `json:"labels,omitempty"`
}

type labelResponse struct {
	Label Label `json:"label,omitempty"`
}

type labelRequestBody struct {
	Label Label `json:"label,omitempty"`
}
