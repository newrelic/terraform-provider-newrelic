package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
)

func flattenAlertMutingRule(mutingRule *alerts.Alerts, d *schema.ResourceData) error {
	d.Set("account_id", mutingRule.Account.ID)
	d.Set("rule", mutingRule.Rule)
	d.Set("condition", mutingRule.Condition) // Need to fix
	d.Set("createdAt", mutingRule.CreatedAt)
	d.Set("createdBy", mutingRule.CreatedBy)
	d.Set("description", mutingRule.Description)

	d.Set("id", mutingRule.ID)
	d.Set("name", mutingRule.Name)
	d.Set("updatedAt", mutingRule.UpdatedAt)
	d.Set("updatedBy", mutingRule.UpdatedBy)

	return nil
}
