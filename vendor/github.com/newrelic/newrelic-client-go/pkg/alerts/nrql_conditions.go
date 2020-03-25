package alerts

import (
	"fmt"

	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

// ThresholdOccurence specifies the threshold occurrence for NRQL alert condition terms.
type ThresholdOccurence string

var (
	// ThresholdOccurrences enumerates the possible threshold occurrence values for NRQL alert condition terms.
	ThresholdOccurrences = struct {
		All         ThresholdOccurence
		AtLeastOnce ThresholdOccurence
	}{
		All:         "ALL",
		AtLeastOnce: "AT_LEAST_ONCE",
	}
)

// NrqlConditionType specifies the type of NRQL alert condition.
type NrqlConditionType string

var (
	// NrqlConditionTypes enumerates the possible NRQL condition type values for NRQL alert conditions.
	NrqlConditionTypes = struct {
		Baseline NrqlConditionType
		Static   NrqlConditionType
	}{
		Baseline: "BASELINE",
		Static:   "STATIC",
	}
)

// NrqlConditionValueFunction specifies the value function of NRQL alert condition.
type NrqlConditionValueFunction string

var (
	// NrqlConditionValueFunctions enumerates the possible NRQL condition value function values for NRQL alert conditions.
	NrqlConditionValueFunctions = struct {
		SingleValue NrqlConditionValueFunction
		Sum         NrqlConditionValueFunction
	}{
		SingleValue: "SINGLE_VALUE",
		Sum:         "SUM",
	}
)

// NrqlConditionValueFunction specifies the value function of NRQL alert condition.
type NrqlConditionViolationTimeLimit string

var (
	// NrqlConditionViolationTimeLimits enumerates the possible NRQL condition violation time limit values for NRQL alert conditions.
	NrqlConditionViolationTimeLimits = struct {
		OneHour         NrqlConditionViolationTimeLimit
		TwoHours        NrqlConditionViolationTimeLimit
		FourHours       NrqlConditionViolationTimeLimit
		EightHours      NrqlConditionViolationTimeLimit
		TwelveHours     NrqlConditionViolationTimeLimit
		TwentyFourHours NrqlConditionViolationTimeLimit
	}{
		OneHour:         "ONE_HOUR",
		TwoHours:        "TWO_HOURS",
		FourHours:       "FOUR_HOURS",
		EightHours:      "EIGHT_HOURS",
		TwelveHours:     "TWELVE_HOURS",
		TwentyFourHours: "TWENTY_FOUR_HOURS",
	}
)

// NrqlConditionOperator specifies the operator for alert condition terms.
type NrqlConditionOperator string

var (
	// NrqlConditionOperators enumerates the possible operator values for alert condition terms.
	NrqlConditionOperators = struct {
		Above NrqlConditionOperator
		Below NrqlConditionOperator
		Equal NrqlConditionOperator
	}{
		Above: "ABOVE",
		Below: "BELOW",
		Equal: "EQUAL",
	}
)

// NrqlConditionPriority specifies the priority for alert condition terms.
type NrqlConditionPriority string

var (
	// NrqlConditionPriorities enumerates the possible priority values for alert condition terms.
	NrqlConditionPriorities = struct {
		Critical NrqlConditionPriority
		Warning  NrqlConditionPriority
	}{
		Critical: "CRITICAL",
		Warning:  "WARNING",
	}
)

type NrqlBaselineDirection string

var (
	// NrqlBaselineDirections enumerates the possible baseline direction values for a baseline NRQL alert condition.
	NrqlBaselineDirections = struct {
		LowerOnly     NrqlBaselineDirection
		UpperAndLower NrqlBaselineDirection
		UpperOnly     NrqlBaselineDirection
	}{
		LowerOnly:     "LOWER_ONLY",
		UpperAndLower: "UPPER_AND_LOWER",
		UpperOnly:     "UPPER_ONLY",
	}
)

// NrqlConditionTerms represents the terms of a New Relic alert condition.
type NrqlConditionTerms struct {
	Operator             NrqlConditionOperator `json:"operator,omitempty"`
	Priority             NrqlConditionPriority `json:"priority,omitempty"`
	Threshold            float64               `json:"threshold,omitempty"`
	ThresholdDuration    float64               `json:"thresholdDuration,omitempty"`
	ThresholdOccurrences ThresholdOccurence    `json:"thresholdOccurrences,omitempty"`
}

// NrqlConditionQuery represents the NRQL query object returned in a NerdGraph response object.
type NrqlConditionQuery struct {
	Query            string `json:"query,omitempty"`
	EvaluationOffset int    `json:"evaluationOffset,omitempty"`
}

// NrqlConditionBase represents the base fields for a New Relic NRQL Alert condition. These fields
// shared between the NrqlConditionMutationInput struct and NrqlConditionMutationResponse struct.
type NrqlConditionBase struct {
	Description        string                          `json:"description,omitempty"`
	Enabled            bool                            `json:"enabled"`
	Name               string                          `json:"name,omitempty"`
	Nrql               NrqlConditionQuery              `json:"nrql,omitempty"`
	RunbookURL         string                          `json:"runbookUrl,omitempty"`
	Terms              []NrqlConditionTerms            `json:"terms,omitempty"`
	ViolationTimeLimit NrqlConditionViolationTimeLimit `json:"violationTimeLimit,omitempty"`
}

// NrqlConditionInput represents the input options for creating or updating a Nrql Condition.
type NrqlConditionInput struct {
	NrqlConditionBase

	// BaselineDirection exists ONLY for NRQL conditions of type BASELINE.
	BaselineDirection *NrqlBaselineDirection `json:"baselineDirection,omitempty"`

	// ValueFunction is returned ONLY for NRQL conditions of type STATIC.
	ValueFunction *NrqlConditionValueFunction `json:"value_function,omitempty"`
}

type NrqlConditionsSearchCriteria struct {
	Name      string `json:"name,omitempty"`
	NameLike  string `json:"nameLike,omitempty"`
	PolicyID  string `json:"policyId,omitempty"`
	Query     string `json:"query,omitempty"`
	QueryLike string `json:"queryLike,omitempty"`
}

// NrqlAlertCondition represents a NerdGraph NRQL alert condition, which is type AlertsNrqlCondition in NerdGraph.
// NrqlAlertCondition could be a baseline condition OR a static condition.
type NrqlAlertCondition struct {
	NrqlConditionBase

	ID       string            `json:"id,omitempty"`
	PolicyID string            `json:"policyId,omitempty"`
	Type     NrqlConditionType `json:"type,omitempty"`

	// BaselineDirection exists ONLY for NRQL conditions of type BASELINE.
	BaselineDirection *NrqlBaselineDirection `json:"baselineDirection,omitempty"`

	// ValueFunction is returned ONLY for NRQL conditions of type STATIC.
	ValueFunction *NrqlConditionValueFunction `json:"value_function,omitempty"`
}

// NrqlCondition represents a New Relic NRQL Alert condition.
type NrqlCondition struct {
	Terms               []ConditionTerm   `json:"terms,omitempty"`
	Nrql                NrqlQuery         `json:"nrql,omitempty"`
	Type                string            `json:"type,omitempty"`
	Name                string            `json:"name,omitempty"`
	RunbookURL          string            `json:"runbook_url,omitempty"`
	ValueFunction       ValueFunctionType `json:"value_function,omitempty"`
	ID                  int               `json:"id,omitempty"`
	ViolationCloseTimer int               `json:"violation_time_limit_seconds,omitempty"`
	ExpectedGroups      int               `json:"expected_groups,omitempty"`
	IgnoreOverlap       bool              `json:"ignore_overlap,omitempty"`
	Enabled             bool              `json:"enabled"`
}

// NrqlQuery represents a NRQL query to use with a NRQL alert condition
type NrqlQuery struct {
	Query      string `json:"query,omitempty"`
	SinceValue string `json:"since_value,omitempty"`
}

// ListNrqlConditions returns NRQL alert conditions for a specified policy.
func (a *Alerts) ListNrqlConditions(policyID int) ([]*NrqlCondition, error) {
	conditions := []*NrqlCondition{}
	queryParams := listNrqlConditionsParams{
		PolicyID: policyID,
	}

	nextURL := "/alerts_nrql_conditions.json"

	for nextURL != "" {
		response := nrqlConditionsResponse{}
		resp, err := a.client.Get(nextURL, &queryParams, &response)

		if err != nil {
			return nil, err
		}

		conditions = append(conditions, response.NrqlConditions...)

		paging := a.pager.Parse(resp)
		nextURL = paging.Next
	}

	return conditions, nil
}

// GetNrqlCondition gets information about a NRQL alert condition
// for a specified policy ID and condition ID.
func (a *Alerts) GetNrqlCondition(policyID int, id int) (*NrqlCondition, error) {
	conditions, err := a.ListNrqlConditions(policyID)
	if err != nil {
		return nil, err
	}

	for _, condition := range conditions {
		if condition.ID == id {
			return condition, nil
		}
	}

	return nil, errors.NewNotFoundf("no condition found for policy %d and condition ID %d", policyID, id)
}

// CreateNrqlCondition creates a NRQL alert condition.
func (a *Alerts) CreateNrqlCondition(policyID int, condition NrqlCondition) (*NrqlCondition, error) {
	reqBody := nrqlConditionRequestBody{
		NrqlCondition: condition,
	}
	resp := nrqlConditionResponse{}

	u := fmt.Sprintf("/alerts_nrql_conditions/policies/%d.json", policyID)
	_, err := a.client.Post(u, nil, &reqBody, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.NrqlCondition, nil
}

// UpdateNrqlCondition updates a NRQL alert condition.
func (a *Alerts) UpdateNrqlCondition(condition NrqlCondition) (*NrqlCondition, error) {
	reqBody := nrqlConditionRequestBody{
		NrqlCondition: condition,
	}
	resp := nrqlConditionResponse{}

	u := fmt.Sprintf("/alerts_nrql_conditions/%d.json", condition.ID)
	_, err := a.client.Put(u, nil, &reqBody, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.NrqlCondition, nil
}

// DeleteNrqlCondition deletes a NRQL alert condition.
func (a *Alerts) DeleteNrqlCondition(id int) (*NrqlCondition, error) {
	resp := nrqlConditionResponse{}
	u := fmt.Sprintf("/alerts_nrql_conditions/%d.json", id)

	_, err := a.client.Delete(u, nil, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.NrqlCondition, nil
}

// GetNrqlConditionQuery fetches a NRQL alert condition via New Relic's NerdGraph API.
func (a *Alerts) GetNrqlConditionQuery(
	accountID int,
	conditionID string,
) (*NrqlAlertCondition, error) {
	resp := getNrqlConditionQueryResponse{}
	vars := map[string]interface{}{
		"accountId": accountID,
		"id":        conditionID,
	}

	if err := a.client.Query(getNrqlConditionQuery, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.Actor.Account.Alerts.NrqlCondition, nil
}

// SearchNrqlConditionsQuery fetches multiple NRQL alert conditions based on the provided search criteria via New Relic's NerdGraph API.
func (a *Alerts) SearchNrqlConditionsQuery(
	accountID int,
	searchCriteria NrqlConditionsSearchCriteria,
) ([]*NrqlAlertCondition, error) {
	conditions := []*NrqlAlertCondition{}
	var nextCursor *string

	for ok := true; ok; ok = nextCursor != nil {
		resp := searchNrqlConditionsResponse{}
		vars := map[string]interface{}{
			"accountId":      accountID,
			"searchCriteria": searchCriteria,
			"cursor":         nextCursor,
		}

		if err := a.client.Query(searchNrqlConditionsQuery, vars, &resp); err != nil {
			return nil, err
		}

		conditions = append(conditions, resp.Actor.Account.Alerts.NrqlConditionsSearch.NrqlConditions...)
		nextCursor = resp.Actor.Account.Alerts.NrqlConditionsSearch.NextCursor
	}

	return conditions, nil
}

// CreateNrqlConditionBaselineMutation creates a baseline NRQL alert condition via New Relic's NerdGraph API.
func (a *Alerts) CreateNrqlConditionBaselineMutation(
	accountID int,
	policyID int,
	nrqlCondition NrqlConditionInput,
) (*NrqlAlertCondition, error) {
	resp := nrqlConditionBaselineCreateResponse{}
	vars := map[string]interface{}{
		"accountId": accountID,
		"policyId":  policyID,
		"condition": nrqlCondition,
	}

	if err := a.client.Query(createNrqlConditionBaselineMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.AlertsNrqlConditionBaselineCreate, nil
}

// UpdateNrqlConditionBaselineMutation updates a baseline NRQL alert condition via New Relic's NerdGraph API.
func (a *Alerts) UpdateNrqlConditionBaselineMutation(
	accountID int,
	conditionID string, // GraphQL scalar type `ID` is a string in JSON
	nrqlCondition NrqlConditionInput,
) (*NrqlAlertCondition, error) {
	resp := nrqlConditionBaselineUpdateResponse{}
	vars := map[string]interface{}{
		"accountId": accountID,
		"id":        conditionID,
		"condition": nrqlCondition,
	}

	if err := a.client.Query(updateNrqlConditionBaselineMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.AlertsNrqlConditionBaselineUpdate, nil
}

// CreateNrqlConditionStaticMutation creates a static NRQL alert condition via New Relic's NerdGraph API.
func (a *Alerts) CreateNrqlConditionStaticMutation(
	accountID int,
	policyID int,
	nrqlCondition NrqlConditionInput,
) (*NrqlAlertCondition, error) {
	resp := nrqlConditionStaticCreateResponse{}
	vars := map[string]interface{}{
		"accountId": accountID,
		"policyId":  policyID,
		"condition": nrqlCondition,
	}

	if err := a.client.Query(createNrqlConditionStaticMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.AlertsNrqlConditionStaticCreate, nil
}

// UpdateNrqlConditionStaticMutation updates a static NRQL alert condition via New Relic's NerdGraph API.
func (a *Alerts) UpdateNrqlConditionStaticMutation(
	accountID int,
	conditionID string, // GraphQL scalar type `ID` is a string in JSON
	nrqlCondition NrqlConditionInput,
) (*NrqlAlertCondition, error) {
	resp := nrqlConditionStaticUpdateResponse{}
	vars := map[string]interface{}{
		"accountId": accountID,
		"id":        conditionID,
		"condition": nrqlCondition,
	}

	if err := a.client.Query(updateNrqlConditionStaticMutation, vars, &resp); err != nil {
		return nil, err
	}

	return &resp.AlertsNrqlConditionStaticUpdate, nil
}

type listNrqlConditionsParams struct {
	PolicyID int `url:"policy_id,omitempty"`
}

type nrqlConditionsResponse struct {
	NrqlConditions []*NrqlCondition `json:"nrql_conditions,omitempty"`
}

type nrqlConditionResponse struct {
	NrqlCondition NrqlCondition `json:"nrql_condition,omitempty"`
}

type nrqlConditionRequestBody struct {
	NrqlCondition NrqlCondition `json:"nrql_condition,omitempty"`
}

type nrqlConditionBaselineCreateResponse struct {
	AlertsNrqlConditionBaselineCreate NrqlAlertCondition `json:"alertsNrqlConditionBaselineCreate"`
}

type nrqlConditionBaselineUpdateResponse struct {
	AlertsNrqlConditionBaselineUpdate NrqlAlertCondition `json:"alertsNrqlConditionBaselineUpdate"`
}

type nrqlConditionStaticCreateResponse struct {
	AlertsNrqlConditionStaticCreate NrqlAlertCondition `json:"alertsNrqlConditionStaticCreate"`
}

type nrqlConditionStaticUpdateResponse struct {
	AlertsNrqlConditionStaticUpdate NrqlAlertCondition `json:"alertsNrqlConditionStaticUpdate"`
}

type searchNrqlConditionsResponse struct {
	Actor struct {
		Account struct {
			Alerts struct {
				NrqlConditionsSearch struct {
					NextCursor     *string
					NrqlConditions []*NrqlAlertCondition `json:"nrqlConditions"`
				} `json:"nrqlConditionsSearch"`
			} `json:"alerts"`
		} `json:"account"`
	} `json:"actor"`
}

type getNrqlConditionQueryResponse struct {
	Actor struct {
		Account struct {
			Alerts struct {
				NrqlCondition NrqlAlertCondition `json:"nrqlCondition"`
			} `json:"alerts"`
		} `json:"account"`
	} `json:"actor"`
}

const (
	graphqlNrqlConditionStructFields = `
		id
		name
		nrql {
			evaluationOffset
			query
		}
		enabled
		description
		policyId
		runbookUrl
		terms {
			operator
			priority
			threshold
			thresholdDuration
			thresholdOccurrences
		}
		type
		violationTimeLimit
	`

	graphqlFragmentNrqlBaselineConditionFields = `
		... on AlertsNrqlBaselineCondition {
			baselineDirection
		}
	`

	graphqlFragmentNrqlStaticConditionFields = `
		... on AlertsNrqlStaticCondition {
			valueFunction
		}
	`
	searchNrqlConditionsQuery = `
		query($accountId: Int!, $searchCriteria: AlertsNrqlConditionsSearchCriteriaInput, $cursor: String) {
			actor {
				account(id: $accountId) {
					alerts {
						nrqlConditionsSearch(searchCriteria: $searchCriteria, cursor: $cursor) {
							nextCursor
							totalCount
							nrqlConditions {` +
		graphqlNrqlConditionStructFields +
		graphqlFragmentNrqlBaselineConditionFields +
		graphqlFragmentNrqlStaticConditionFields +
		`} } } } } }`

	getNrqlConditionQuery = `
		query ($accountId: Int!, $id: ID!) {
			actor {
				account(id: $accountId) {
					alerts {
						nrqlCondition(id: $id) {` +
		graphqlNrqlConditionStructFields +
		graphqlFragmentNrqlBaselineConditionFields +
		graphqlFragmentNrqlStaticConditionFields +
		`} } } } }`

	// Baseline
	createNrqlConditionBaselineMutation = `
		mutation($accountId: Int!, $policyId: ID!, $condition: AlertsNrqlConditionBaselineInput!) {
			alertsNrqlConditionBaselineCreate(accountId: $accountId, policyId: $policyId, condition: $condition) {
				baselineDirection` +
		graphqlNrqlConditionStructFields +
		` } }`

	// Baseline
	updateNrqlConditionBaselineMutation = `
		mutation($accountId: Int!, $id: ID!, $condition: AlertsNrqlConditionUpdateBaselineInput!) {
			alertsNrqlConditionBaselineUpdate(accountId: $accountId, id: $id, condition: $condition) {
				baselineDirection` +
		graphqlNrqlConditionStructFields +
		` } }`

	// Static
	createNrqlConditionStaticMutation = `
		mutation($accountId: Int!, $policyId: ID!, $condition: AlertsNrqlConditionStaticInput!) {
			alertsNrqlConditionStaticCreate(accountId: $accountId, policyId: $policyId, condition: $condition) {
				valueFunction` +
		graphqlNrqlConditionStructFields +
		` } }`

	// Static
	updateNrqlConditionStaticMutation = `
		mutation($accountId: Int!, $id: ID!, $condition: AlertsNrqlConditionUpdateStaticInput!) {
			alertsNrqlConditionStaticUpdate(accountId: $accountId, id: $id, condition: $condition) {
				valueFunction` +
		graphqlNrqlConditionStructFields +
		` } }`
)
