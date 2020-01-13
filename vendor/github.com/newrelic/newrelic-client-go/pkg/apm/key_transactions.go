package apm

import (
	"fmt"
)

// ListKeyTransactionsParams represents a set of filters to be
// used when querying New Relic key transactions.
type ListKeyTransactionsParams struct {
	Name string `url:"filter[name],omitempty"`
	IDs  []int  `url:"filter[ids],omitempty,comma"`
}

// ListKeyTransactions returns all key transactions for an account.
func (apm *APM) ListKeyTransactions(params *ListKeyTransactionsParams) ([]*KeyTransaction, error) {
	response := keyTransactionsResponse{}
	results := []*KeyTransaction{}
	nextURL := "/key_transactions.json"

	for nextURL != "" {
		resp, err := apm.client.Get(nextURL, &params, &response)

		if err != nil {
			return nil, err
		}

		results = append(results, response.KeyTransactions...)

		paging := apm.pager.Parse(resp)
		nextURL = paging.Next
	}

	return results, nil
}

// GetKeyTransaction returns a specific key transaction by ID.
func (apm *APM) GetKeyTransaction(id int) (*KeyTransaction, error) {
	response := keyTransactionResponse{}
	u := fmt.Sprintf("/key_transactions/%d.json", id)

	_, err := apm.client.Get(u, nil, &response)

	if err != nil {
		return nil, err
	}

	return &response.KeyTransaction, nil
}

type keyTransactionsResponse struct {
	KeyTransactions []*KeyTransaction `json:"key_transactions,omitempty"`
}

type keyTransactionResponse struct {
	KeyTransaction KeyTransaction `json:"key_transaction,omitempty"`
}
