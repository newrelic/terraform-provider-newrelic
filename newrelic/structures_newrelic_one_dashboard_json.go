package newrelic

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/newrelic/newrelic-client-go/v2/pkg/dashboards"
)

// Assemble the *dashboards.DashboardInput struct.
// Used by the newrelic_one_dashboard Create function.
func expandDashboardJSONInput(d *schema.ResourceData, meta interface{}) (*dashboards.DashboardInput, error) {
	dash := dashboards.DashboardInput{}
	err := json.Unmarshal([]byte(d.Get("json").(string)), &dash)
	if err != nil {
		return nil, err
	}

	return &dash, nil
}
