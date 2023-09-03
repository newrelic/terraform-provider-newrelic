package newrelic

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/newrelic/newrelic-client-go/v2/pkg/dashboards"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
)

// Assemble the *dashboards.DashboardInput struct.
// Used by the newrelic_one_dashboard Create function.
func expandDashboardJSONInput(d *schema.ResourceData, meta interface{}) (*dashboards.DashboardInput, error) {
	dash := dashboards.DashboardInput{}
	err := json.Unmarshal([]byte(d.Get("json").(string)), &dash)
	if err != nil {
		return nil, err
	}

	if dataAccounts := d.Get("data_accounts"); dataAccounts != nil {
		if err := setDataAccounts(&dash, dataAccounts); err != nil {
			return nil, fmt.Errorf("failed to set data_accounts for dashboard: %w", err)
		}
	}

	return &dash, nil
}

func setDataAccounts(dashboard *dashboards.DashboardInput, dataAccountsRaw any) error {
	dataAccountsAsserted, ok := dataAccountsRaw.([]any)
	if !ok {
		return fmt.Errorf("expected dataAccounts to be []any, was %T", dataAccountsRaw)
	}
	if len(dataAccountsAsserted) == 0 {
		return nil
	}
	dataAccounts, err := anySliceToIntSlice(dataAccountsAsserted)
	if err != nil {
		return fmt.Errorf("failed to convert dataAccounts to []int: %w", err)
	}
	for pageIdx, page := range dashboard.Pages {
		for widgetIdx, widget := range page.Widgets {
			var rawConfig map[string]any
			if err := json.Unmarshal([]byte(widget.RawConfiguration), &rawConfig); err != nil {
				return fmt.Errorf("failed to unmarshal pages.%d.widgets.%d.rawConfiguration: %w", pageIdx, widgetIdx, err)
			}
			queriesRaw, hasQueries := rawConfig["nrqlQueries"]
			if !hasQueries {
				continue
			}
			queriesSlice, ok := queriesRaw.([]any)
			if !ok {
				return fmt.Errorf("expected pages.%d.widgets.%d.rawConfiguration to be []interface{}, was %T", pageIdx, widgetIdx, queriesRaw)
			}
			for queryIdx, queryRaw := range queriesSlice {
				query, ok := queryRaw.(map[string]any)
				if !ok {
					return fmt.Errorf("expected pages.%d.widgets.%d.rawConfiguration.nrqlQueries.%d to be of type map[string]any, was %T", pageIdx, widgetIdx, queryIdx, queryRaw)
				}
				query["accountIds"] = dataAccounts
				queriesSlice[queryIdx] = query
			}
			serialized, err := json.Marshal(rawConfig)
			if err != nil {
				return fmt.Errorf("failed to serialize pages.%d.widgets.%d.rawConfiguration: %w", pageIdx, widgetIdx, err)
			}
			dashboard.Pages[pageIdx].Widgets[widgetIdx].RawConfiguration = entities.DashboardWidgetRawConfiguration(serialized)
		}
	}

	for variableIdx := range dashboard.Variables {
		if dashboard.Variables[variableIdx].NRQLQuery != nil {
			dashboard.Variables[variableIdx].NRQLQuery.AccountIDs = dataAccounts
		}
	}

	return nil
}

func anySliceToIntSlice(i []any) ([]int, error) {
	result := make([]int, 0, len(i))
	for _, e := range i {
		asserted, ok := e.(int)
		if !ok {
			return nil, fmt.Errorf("expected member to be of type int, was %T", e)
		}
		result = append(result, asserted)
	}

	return result, nil
}
