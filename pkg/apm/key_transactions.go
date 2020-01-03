package apm

import (
	"fmt"

	"github.com/newrelic/newrelic-client-go/internal/utils"
)

// ListKeyTransactionsParams represents a set of filters to be
// used when querying New Relic key transactions.
type ListKeyTransactionsParams struct {
	Name string
	IDs  []int
}

// ListKeyTransactions returns all key transactions for an account.
func (apm *APM) ListKeyTransactions(params *ListKeyTransactionsParams) ([]*KeyTransaction, error) {
	response := keyTransactionsResponse{}
	results := []*KeyTransaction{}
	paramsMap := buildListKeyTransactionsParamsMap(params)
	nextURL := "/key_transactions.json"

	for nextURL != "" {
		resp, err := apm.client.Get(nextURL, &paramsMap, &response)

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

func buildListKeyTransactionsParamsMap(params *ListKeyTransactionsParams) map[string]string {
	paramsMap := map[string]string{}

	if params == nil {
		return paramsMap
	}

	if params.Name != "" {
		paramsMap["filter[name]"] = params.Name
	}

	if params.IDs != nil && len(params.IDs) > 0 {
		paramsMap["filter[ids]"] = utils.IntArrayToString(params.IDs)
	}

	return paramsMap
}

type keyTransactionsResponse struct {
	KeyTransactions []*KeyTransaction `json:"key_transactions,omitempty"`
}

type keyTransactionResponse struct {
	KeyTransaction KeyTransaction `json:"key_transaction,omitempty"`
}
