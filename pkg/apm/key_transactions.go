package apm

func (apm *APM) ListKeyTransactions() ([]*KeyTransaction, error) {
	response := keyTransactionsResponse{}
	results := []*KeyTransaction{}
	nextURL := "/key_transactions.json"

	for nextURL != "" {
		resp, err := apm.client.Get(nextURL, nil, &response)

		if err != nil {
			return nil, err
		}

		results = append(results, response.KeyTransactions...)

		paging := apm.pager.Parse(resp)
		nextURL = paging.Next
	}

	return results, nil
}

type keyTransactionsResponse struct {
	KeyTransactions []*KeyTransaction `json:"key_transactions,omitempty"`
}
