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

	if v, ok := d.GetOk("timestamp"); ok {
		eventInput.Timestamp = nrtime.EpochMilliseconds(time.UnixMilli(int64(v.(int))))
	}

	if v, ok := d.GetOk("user"); ok {
		eventInput.User = v.(string)
	}

	if v, ok := d.GetOk("entity_search"); ok {
		items := v.([]interface{})
		if len(items) > 0 {
			eventInput.EntitySearch = expandChangeTrackingEventEntitySearch(items[0].(map[string]interface{}))
		}
	}

	if v, ok := d.GetOk("category_and_type_data"); ok {
		items := v.([]interface{})
		if len(items) > 0 {
			catData := expandChangeTrackingEventCategoryRelated(items[0].(map[string]interface{}))
			eventInput.CategoryAndTypeData = catData
		}
	}

	dataHandlingRules := expandChangeTrackingEventDataHandlingRules(d)

	return eventInput, dataHandlingRules, nil
}

func expandChangeTrackingEventEntitySearch(cfg map[string]interface{}) changetracking.ChangeTrackingEntitySearchInput {
	input := changetracking.ChangeTrackingEntitySearchInput{}
	if v, ok := cfg["query"]; ok {
		input.Query = v.(string)
	}
	return input
}

func expandChangeTrackingEventCategoryRelated(cfg map[string]interface{}) *changetracking.ChangeTrackingCategoryRelatedInput {
	input := &changetracking.ChangeTrackingCategoryRelatedInput{}

	if v, ok := cfg["kind"]; ok {
		items := v.([]interface{})
		if len(items) > 0 {
			kind := expandChangeTrackingEventCategoryAndType(items[0].(map[string]interface{}))
			input.Kind = kind
		}
	}

	if v, ok := cfg["category_fields"]; ok {
		items := v.([]interface{})
		if len(items) > 0 {
			fields := expandChangeTrackingEventCategoryFields(items[0].(map[string]interface{}))
			input.CategoryFields = fields
		}
	}

	return input
}

func expandChangeTrackingEventCategoryAndType(cfg map[string]interface{}) *changetracking.ChangeTrackingCategoryAndTypeInput {
	input := &changetracking.ChangeTrackingCategoryAndTypeInput{}
	if v, ok := cfg["category"]; ok {
		input.Category = v.(string)
	}
	if v, ok := cfg["type"]; ok {
		input.Type = v.(string)
	}
	return input
}

func expandChangeTrackingEventCategoryFields(cfg map[string]interface{}) *changetracking.ChangeTrackingCategoryFieldsInput {
	input := &changetracking.ChangeTrackingCategoryFieldsInput{}

	if v, ok := cfg["deployment"]; ok {
		items := v.([]interface{})
		if len(items) > 0 {
			input.Deployment = expandChangeTrackingEventDeploymentFields(items[0].(map[string]interface{}))
		}
	}

	if v, ok := cfg["feature_flag"]; ok {
		items := v.([]interface{})
		if len(items) > 0 {
			input.FeatureFlag = expandChangeTrackingEventFeatureFlagFields(items[0].(map[string]interface{}))
		}
	}

	return input
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

	if v, ok := d.GetOk("data_handling_rules"); ok {
		items := v.([]interface{})
		if len(items) > 0 {
			cfg := items[0].(map[string]interface{})
			if flags, ok := cfg["validation_flags"]; ok {
				flagList := flags.([]interface{})
				validationFlags := make([]changetracking.ChangeTrackingValidationFlag, len(flagList))
				for i, f := range flagList {
					validationFlags[i] = changetracking.ChangeTrackingValidationFlag(f.(string))
				}
				rules.ValidationFlags = validationFlags
			}
		}
	}

	return rules
}

// No flatten functions: this API is write-only.
