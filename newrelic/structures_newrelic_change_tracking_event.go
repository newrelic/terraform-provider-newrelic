package newrelic

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/changetracking"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
)

func expandChangeTrackingEvent(d *schema.ResourceData) (changetracking.ChangeTrackingCreateEventInput, changetracking.ChangeTrackingDataHandlingRules, error) {
	eventInput := changetracking.ChangeTrackingCreateEventInput{}

	if v, ok := d.GetOk("description"); ok {
		eventInput.Description = v.(string)
	}

	if v, ok := d.GetOk("group_id"); ok {
		eventInput.GroupId = v.(string)
	}

	if v, ok := d.GetOk("short_description"); ok {
		eventInput.ShortDescription = v.(string)
	}

	if v, ok := d.GetOk("user"); ok {
		eventInput.User = v.(string)
	}

	if v, ok := d.GetOk("timestamp"); ok {
		eventInput.Timestamp = nrtime.EpochMilliseconds(time.UnixMilli(int64(v.(int))))
	}

	if entitySearch, ok := expandChangeTrackingEventEntitySearch(d); ok {
		eventInput.EntitySearch = entitySearch
	}

	if categoryData, hasCategory := expandChangeTrackingEventCategoryAndTypeData(d); hasCategory {
		eventInput.CategoryAndTypeData = categoryData
	}

	dataHandlingRules := expandChangeTrackingEventDataHandlingRules(d)

	return eventInput, dataHandlingRules, nil
}

func expandChangeTrackingEventEntitySearch(d *schema.ResourceData) (changetracking.ChangeTrackingEntitySearchInput, bool) {
	v, ok := d.GetOk("entity_search")
	if !ok {
		return changetracking.ChangeTrackingEntitySearchInput{}, false
	}

	items := v.([]interface{})
	if len(items) == 0 {
		return changetracking.ChangeTrackingEntitySearchInput{}, false
	}

	cfg := items[0].(map[string]interface{})
	input := changetracking.ChangeTrackingEntitySearchInput{}

	if q, ok := cfg["query"]; ok {
		input.Query = q.(string)
	}

	return input, true
}

func expandChangeTrackingEventCategoryAndTypeData(d *schema.ResourceData) (*changetracking.ChangeTrackingCategoryRelatedInput, bool) {
	v, ok := d.GetOk("category_and_type_data")
	if !ok {
		return nil, false
	}

	items := v.([]interface{})
	if len(items) == 0 {
		return nil, false
	}

	cfg := items[0].(map[string]interface{})
	input := &changetracking.ChangeTrackingCategoryRelatedInput{}

	if kindRaw, ok := cfg["kind"]; ok {
		kindList := kindRaw.([]interface{})
		if len(kindList) > 0 {
			kindCfg := kindList[0].(map[string]interface{})
			kind := &changetracking.ChangeTrackingCategoryAndTypeInput{}
			if cat, ok := kindCfg["category"]; ok {
				kind.Category = cat.(string)
			}
			if typ, ok := kindCfg["type"]; ok {
				kind.Type = typ.(string)
			}
			input.Kind = kind
		}
	}

	if cfRaw, ok := cfg["category_fields"]; ok {
		cfList := cfRaw.([]interface{})
		if len(cfList) > 0 {
			input.CategoryFields = expandChangeTrackingEventCategoryFields(cfList[0].(map[string]interface{}))
		}
	}

	return input, true
}

func expandChangeTrackingEventCategoryFields(cfg map[string]interface{}) *changetracking.ChangeTrackingCategoryFieldsInput {
	fields := &changetracking.ChangeTrackingCategoryFieldsInput{}

	if depRaw, ok := cfg["deployment"]; ok {
		depList := depRaw.([]interface{})
		if len(depList) > 0 {
			fields.Deployment = expandChangeTrackingEventDeploymentFields(depList[0].(map[string]interface{}))
		}
	}

	if ffRaw, ok := cfg["feature_flag"]; ok {
		ffList := ffRaw.([]interface{})
		if len(ffList) > 0 {
			fields.FeatureFlag = expandChangeTrackingEventFeatureFlagFields(ffList[0].(map[string]interface{}))
		}
	}

	return fields
}

func expandChangeTrackingEventDeploymentFields(cfg map[string]interface{}) *changetracking.ChangeTrackingDeploymentFieldsInput {
	input := &changetracking.ChangeTrackingDeploymentFieldsInput{}

	if v, ok := cfg["changelog"]; ok {
		input.Changelog = v.(string)
	}

	if v, ok := cfg["commit"]; ok {
		input.Commit = v.(string)
	}

	if v, ok := cfg["deep_link"]; ok {
		input.DeepLink = v.(string)
	}

	if v, ok := cfg["version"]; ok {
		input.Version = v.(string)
	}

	return input
}

func expandChangeTrackingEventFeatureFlagFields(cfg map[string]interface{}) *changetracking.ChangeTrackingFeatureFlagFieldsInput {
	input := &changetracking.ChangeTrackingFeatureFlagFieldsInput{}

	if v, ok := cfg["feature_flag_id"]; ok {
		input.FeatureFlagId = v.(string)
	}

	return input
}

func expandChangeTrackingEventDataHandlingRules(d *schema.ResourceData) changetracking.ChangeTrackingDataHandlingRules {
	rules := changetracking.ChangeTrackingDataHandlingRules{}

	v, ok := d.GetOk("data_handling_rules")
	if !ok {
		return rules
	}

	items := v.([]interface{})
	if len(items) == 0 {
		return rules
	}

	cfg := items[0].(map[string]interface{})

	if flagsRaw, ok := cfg["validation_flags"]; ok {
		flagsList := flagsRaw.([]interface{})
		flags := make([]changetracking.ChangeTrackingValidationFlag, 0, len(flagsList))
		for _, f := range flagsList {
			flags = append(flags, changetracking.ChangeTrackingValidationFlag(f.(string)))
		}
		rules.ValidationFlags = flags
	}

	return rules
}

// Ensure imported packages are used.
var _ = common.EntityGUID("")

// No flatten functions: this API is write-only.
