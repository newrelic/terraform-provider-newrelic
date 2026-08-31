package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logconfigurations"
)

func expandLogParsingRuleInput(d *schema.ResourceData) logconfigurations.LogConfigurationsParsingRuleConfiguration {
	inp := logconfigurations.LogConfigurationsParsingRuleConfiguration{}

	if e, ok := d.GetOk("attribute"); ok {
		inp.Attribute = e.(string)
	}

	if e, ok := d.GetOk("enabled"); ok {
		inp.Enabled = e.(bool)
	}

	if e, ok := d.GetOk("name"); ok {
		inp.Description = e.(string)
	}

	if e, ok := d.GetOk("grok"); ok {
		inp.Grok = e.(string)
	}

	if e, ok := d.GetOk("lucene"); ok {
		inp.Lucene = e.(string)
	}

	if e, ok := d.GetOk("nrql"); ok {
		inp.NRQL = logconfigurations.NRQL(e.(string))
	}

	if e, ok := d.GetOk("source"); ok {
		inp.Source = logconfigurations.LogConfigurationsParsingRuleSource(e.(string))
	}

	return inp
}

func expandLogParsingRuleUpdateInput(d *schema.ResourceData) logconfigurations.LogConfigurationsParsingRuleConfiguration {
	updateInp := logconfigurations.LogConfigurationsParsingRuleConfiguration{}

	if e, ok := d.GetOk("attribute"); ok {
		updateInp.Attribute = e.(string)
	}

	if e, ok := d.GetOk("enabled"); ok {
		updateInp.Enabled = e.(bool)
	}

	if e, ok := d.GetOk("name"); ok {
		updateInp.Description = e.(string)
	}

	if e, ok := d.GetOk("grok"); ok {
		updateInp.Grok = e.(string)
	}

	if e, ok := d.GetOk("lucene"); ok {
		updateInp.Lucene = e.(string)
	}

	if e, ok := d.GetOk("nrql"); ok {
		updateInp.NRQL = logconfigurations.NRQL(e.(string))
	}

	if e, ok := d.GetOk("source"); ok {
		updateInp.Source = logconfigurations.LogConfigurationsParsingRuleSource(e.(string))
	}

	return updateInp
}
func flattenLogParsingRule(rule *logconfigurations.LogConfigurationsParsingRule, d *schema.ResourceData) error {
	_ = d.Set("account_id", rule.AccountID)
	_ = d.Set("attribute", rule.Attribute)
	_ = d.Set("name", rule.Description)
	_ = d.Set("enabled", rule.Enabled)
	_ = d.Set("grok", rule.Grok)
	_ = d.Set("lucene", rule.Lucene)
	_ = d.Set("nrql", string(rule.NRQL))
	_ = d.Set("deleted", rule.Deleted)
	_ = d.Set("source", string(rule.Source))
	return nil
}