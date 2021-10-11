package newrelic

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
)

func expandInfraAlertCondition(d *schema.ResourceData) (*alerts.InfrastructureCondition, error) {
	condition := alerts.InfrastructureCondition{
		Name:        d.Get("name").(string),
		Enabled:     d.Get("enabled").(bool),
		PolicyID:    d.Get("policy_id").(int),
		Event:       d.Get("event").(string),
		Comparison:  strings.ToLower(d.Get("comparison").(string)),
		Select:      d.Get("select").(string),
		Type:        strings.ToLower(d.Get("type").(string)),
		Critical:    expandInfraAlertThreshold(d.Get("critical")),
		Description: d.Get("description").(string),
	}

	if attr, ok := d.GetOk("runbook_url"); ok {
		condition.RunbookURL = attr.(string)
	}

	if attr, ok := d.GetOk("warning"); ok {
		condition.Warning = expandInfraAlertThreshold(attr)
	}

	if attr, ok := d.GetOk("where"); ok {
		condition.Where = attr.(string)
	}

	if attr, ok := d.GetOk("process_where"); ok {
		condition.ProcessWhere = attr.(string)
	}

	if attr, ok := d.GetOk("integration_provider"); ok {
		condition.IntegrationProvider = attr.(string)
	}

	if attr, ok := d.GetOkExists("violation_close_timer"); ok {
		t := attr.(int)
		condition.ViolationCloseTimer = &t
	}

	err := validateAttributesForType(&condition)

	if err != nil {
		return nil, err
	}

	return &condition, nil
}

func expandInfraAlertThreshold(v interface{}) *alerts.InfrastructureConditionThreshold {
	if len(v.([]interface{})) < 1 {
		return nil
	}

	rah := v.([]interface{})[0].(map[string]interface{})

	alertInfraThreshold := &alerts.InfrastructureConditionThreshold{
		Duration: rah["duration"].(int),
	}

	if val, ok := rah["value"]; ok {
		value := val.(float64)
		alertInfraThreshold.Value = &value
	}

	if val, ok := rah["time_function"]; ok {
		alertInfraThreshold.Function = strings.ToLower(val.(string))
	}

	return alertInfraThreshold
}

func flattenInfraAlertCondition(condition *alerts.InfrastructureCondition, d *schema.ResourceData) error {
	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]

	_ = d.Set("policy_id", policyID)
	_ = d.Set("name", condition.Name)
	_ = d.Set("runbook_url", condition.RunbookURL)
	_ = d.Set("enabled", condition.Enabled)
	_ = d.Set("comparison", strings.ToLower(condition.Comparison))
	_ = d.Set("event", condition.Event)
	_ = d.Set("select", condition.Select)
	_ = d.Set("type", strings.ToLower(condition.Type))
	_ = d.Set("description", condition.Description)
	_ = d.Set("created_at", time.Time(*condition.CreatedAt).Unix())
	_ = d.Set("updated_at", time.Time(*condition.UpdatedAt).Unix())

	if condition.Where != "" {
		_ = d.Set("where", condition.Where)
	}

	if condition.ProcessWhere != "" {
		_ = d.Set("process_where", condition.ProcessWhere)
	}

	if condition.IntegrationProvider != "" {
		_ = d.Set("integration_provider", condition.IntegrationProvider)
	}

	if condition.ViolationCloseTimer != nil {
		_ = d.Set("violation_close_timer", condition.ViolationCloseTimer)
	}

	if err := d.Set("critical", flattenAlertThreshold(condition.Critical)); err != nil {
		return err
	}

	if condition.Warning != nil {
		if err := d.Set("warning", flattenAlertThreshold(condition.Warning)); err != nil {
			return err
		}
	}

	return nil
}

func flattenAlertThreshold(v *alerts.InfrastructureConditionThreshold) []interface{} {
	alertInfraThreshold := map[string]interface{}{
		"duration":      v.Duration,
		"time_function": strings.ToLower(v.Function),
	}

	if v.Value != nil {
		alertInfraThreshold["value"] = *v.Value
	}

	return []interface{}{alertInfraThreshold}
}

func validateAttributesForType(c *alerts.InfrastructureCondition) error {
	switch c.Type {
	case "infra_process_running":
		if c.Event != "" {
			return fmt.Errorf("event is not supported by condition type %s", c.Type)
		}
		if c.IntegrationProvider != "" {
			return fmt.Errorf("integration_provider is not supported by condition type %s", c.Type)
		}
		if c.Select != "" {
			return fmt.Errorf("select is not supported by condition type %s", c.Type)
		}
		if c.Critical.Function != "" {
			return fmt.Errorf("time_function is not supported by condition type %s", c.Type)
		}
	case "infra_metric":
		if c.ProcessWhere != "" {
			return fmt.Errorf("process_where is not supported by condition type %s", c.Type)
		}
	case "infra_host_not_reporting":
		if c.Event != "" {
			return fmt.Errorf("event is not supported by condition type %s", c.Type)
		}
		if c.IntegrationProvider != "" {
			return fmt.Errorf("integration_provider is not supported by condition type %s", c.Type)
		}
		if c.Select != "" {
			return fmt.Errorf("select is not supported by condition type %s", c.Type)
		}
		if c.ProcessWhere != "" {
			return fmt.Errorf("process_where is not supported by condition type %s", c.Type)
		}
		if c.Comparison != "" {
			return fmt.Errorf("comparison is not supported by condition type %s", c.Type)
		}
		if c.Critical.Function != "" {
			return fmt.Errorf("time_function is not supported by condition type %s", c.Type)
		}
		if c.Critical.Value != nil && *c.Critical.Value != 0.0 {
			return fmt.Errorf("value is not supported by condition type %s", c.Type)
		}
	}

	return nil
}
