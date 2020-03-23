package newrelic

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
)

func flattenAlertPolicyDataSource(policy *alerts.AlertsPolicy, d *schema.ResourceData, accountID int) error {
	d.SetId(strconv.Itoa(policy.ID))

	var err error

	err = d.Set("name", policy.Name)
	if err != nil {
		return err
	}

	err = d.Set("incident_preference", policy.IncidentPreference)
	if err != nil {
		return err
	}

	err = d.Set("account_id", accountID)
	if err != nil {
		return err
	}

	return nil
}

func flattenAlertPolicy(policy *alerts.AlertsPolicy, d *schema.ResourceData, accountID int) error {
	var err error

	err = d.Set("name", policy.Name)
	if err != nil {
		return err
	}

	err = d.Set("incident_preference", policy.IncidentPreference)
	if err != nil {
		return err
	}

	err = d.Set("account_id", accountID)
	if err != nil {
		return err
	}

	err = d.Set("policy_id", policy.ID)
	if err != nil {
		return err
	}

	return nil
}
