package newrelic

import (
	"fmt"
	"strings"
	"encoding/base64"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/alerts"
)

func expandMultiLocationSyntheticsCondition(d *schema.ResourceData) (*alerts.MultiLocationSyntheticsCondition, error) {
	condition := alerts.MultiLocationSyntheticsCondition{
		Name:                      d.Get("name").(string),
		Enabled:                   d.Get("enabled").(bool),
		ViolationTimeLimitSeconds: d.Get("violation_time_limit_seconds").(int),
	}

	if attr, ok := d.GetOk("runbook_url"); ok {
		condition.RunbookURL = attr.(string)
	}

	terms, err := expandMultiLocationSyntheticsConditionTerms(d)
	if err != nil {
		return nil, err
	}

	condition.Terms = terms

	var entities []string
	for _, x := range d.Get("entities").([]interface{}) {
		entities = append(entities, getMonitorID(x.(string)))
	}

	condition.Entities = entities

	return &condition, nil
}

func expandMultiLocationSyntheticsConditionTerms(d *schema.ResourceData) ([]alerts.MultiLocationSyntheticsConditionTerm, error) {
	var expandedTerms []alerts.MultiLocationSyntheticsConditionTerm

	if critical, ok := d.GetOk("critical"); ok {
		x := critical.([]interface{})
		// A critical attribute is a list, but is limited to a single item in the shema.
		if len(x) > 0 {
			single := x[0].(map[string]interface{})

			criticalTerm, err := expandMultiLocationSyntheticsConditionTerm(single, "critical")
			if err != nil {
				return nil, err
			}
			if criticalTerm != nil {
				expandedTerms = append(expandedTerms, *criticalTerm)
			}
		}
	}

	if warning, ok := d.GetOk("warning"); ok {
		x := warning.([]interface{})
		// A warning attribute is a list, but is limited to a single item in the shema.
		if len(x) > 0 {
			single := x[0].(map[string]interface{})

			warningTerm, err := expandMultiLocationSyntheticsConditionTerm(single, "warning")
			if err != nil {
				return nil, err
			}

			if warningTerm != nil {
				expandedTerms = append(expandedTerms, *warningTerm)
			}
		}
	}

	return expandedTerms, nil
}

func expandMultiLocationSyntheticsConditionTerm(term map[string]interface{}, priority string) (*alerts.MultiLocationSyntheticsConditionTerm, error) {
	// required
	threshold := term["threshold"].(int)

	return &alerts.MultiLocationSyntheticsConditionTerm{
		Priority:  strings.ToLower(priority),
		Threshold: threshold,
	}, nil
}

func getMultiLocationSyntheticsConditionEntityGUID(condition *alerts.MultiLocationSyntheticsCondition, accountID int) string {
	rawGUID := fmt.Sprintf("%d|AIOPS|CONDITION|%d", accountID, condition.ID)
	return base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(rawGUID))
}

func flattenMultiLocationSyntheticsCondition(condition *alerts.MultiLocationSyntheticsCondition, accountID int, d *schema.ResourceData) error {
	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	entityGUID := getMultiLocationSyntheticsConditionEntityGUID(condition, accountID)

	_ = d.Set("policy_id", policyID)
	_ = d.Set("name", condition.Name)
	_ = d.Set("runbook_url", condition.RunbookURL)
	_ = d.Set("enabled", condition.Enabled)
	_ = d.Set("violation_time_limit_seconds", condition.ViolationTimeLimitSeconds)
	_ = d.Set("policy_id", policyID)
	_ = d.Set("entity_guid", entityGUID)

	for _, term := range condition.Terms {
		switch term.Priority {
		case "critical":
			terms := []map[string]interface{}{
				{
					"threshold": term.Threshold,
				},
			}
			if err := d.Set("critical", terms); err != nil {
				return fmt.Errorf("[DEBUG] Error setting synthetics multi-location alert condition `critical`: %v", err)
			}
		case "warning":
			terms := []map[string]interface{}{
				{
					"threshold": term.Threshold,
				},
			}
			if err := d.Set("warning", terms); err != nil {
				return fmt.Errorf("[DEBUG] Error setting synthetics multi-location alert condition `warning`: %v", err)
			}
		}
	}

	return nil
}
