package newrelic

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/changetracking"
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

	if categoryData, ok := expandChangeTrackingEventCategoryAndTypeData(d); ok {
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
	input := changetracking.ChangeTrackingEntitySearchInput{
		Query: cfg["query"].(string),
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

	if kindList, ok := cfg["kind"].([]interface{}); ok && len(kindList) > 0 {
		kindCfg := kindList[0].(map[string]interface{})
		input.Kind = &changetracking.ChangeTrackingCategoryAndTypeInput{
			Category: kindCfg["category"].(string),
			Type:     kindCfg["type"].(string),
		}
	}

	if cfList, ok := cfg["category_fields"].([]interface{}); ok && len(cfList) > 0 {
		cfCfg := cfList[0].(map[string]interface{})
		input.CategoryFields = expandChangeTrackingEventCategoryFields(cfCfg)
	}

	return input, true
}

func expandChangeTrackingEventCategoryFields(cfg map[string]interface{}) *changetracking.ChangeTrackingCategoryFieldsInput {
	fields := &changetracking.ChangeTrackingCategoryFieldsInput{}

	if depList, ok := cfg["deployment"].([]interface{}); ok && len(depList) > 0 {
		depCfg := depList[0].(map[string]interface{})
		fields.Deployment = expandChangeTrackingEventDeploymentFields(depCfg)
	}

	if ffList, ok := cfg["feature_flag"].([]interface{}); ok && len(ffList) > 0 {
		ffCfg := ffList[0].(map[string]interface{})
		fields.FeatureFlag = &changetracking.ChangeTrackingFeatureFlagFieldsInput{
			FeatureFlagId: ffCfg["feature_flag_id"].(string),
		}
	}

	return fields
}

func expandChangeTrackingEventDeploymentFields(cfg map[string]interface{}) *changetracking.ChangeTrackingDeploymentFieldsInput {
	dep := &changetracking.ChangeTrackingDeploymentFieldsInput{
		Version: cfg["version"].(string),
	}

	if v, ok := cfg["changelog"].(string); ok {
		dep.Changelog = v
	}

	if v, ok := cfg["commit"].(string); ok {
		dep.Commit = v
	}

	if v, ok := cfg["deep_link"].(string); ok {
		dep.DeepLink = v
	}

	return dep
}

func expandChangeTrackingEventDataHandlingRules(d *schema.ResourceData) changetracking.ChangeTrackingDataHandlingRules {
	rules := changetracking.ChangeTrackingDataHandlingRules{}

	if v, ok := d.GetOk("validation_flags"); ok {
		flagList := v.([]interface{})
		flags := make([]changetracking.ChangeTrackingValidationFlag, len(flagList))
		for i, f := range flagList {
			flags[i] = changetracking.ChangeTrackingValidationFlag(f.(string))
		}
		rules.ValidationFlags = flags
	}

	return rules
}

// No flatten functions: this API is write-only.
