package newrelic

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
)

func expandAlertPolicy(d *schema.ResourceData) *alerts.Policy {
	policy := alerts.Policy{
		Name: d.Get("name").(string),
	}

	if attr, ok := d.GetOk("incident_preference"); ok {
		policy.IncidentPreference = alerts.IncidentPreferenceType(attr.(string))
	}

	return &policy
}

func flattenAlertPolicyDataSource(policy *alerts.Policy, d *schema.ResourceData) error {
	d.SetId(strconv.Itoa(policy.ID))

	return flattenAlertPolicy(policy, d)
}

func flattenAlertPolicy(policy *alerts.Policy, d *schema.ResourceData) error {
	// New Relic provides created_at and updated_at as millisecond unix timestamps
	// https://www.terraform.io/docs/extend/schemas/schema-types.html#date-amp-time-data
	// "TypeString is also used for date/time data, the preferred format is RFC 3339."
	created := time.Time(*policy.CreatedAt).Format(time.RFC3339)
	updated := time.Time(*policy.UpdatedAt).Format(time.RFC3339)

	d.Set("name", policy.Name)
	d.Set("incident_preference", policy.IncidentPreference)
	d.Set("created_at", created)
	d.Set("updated_at", updated)

	return nil
}
