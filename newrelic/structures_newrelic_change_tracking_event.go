package newrelic

import (
	"encoding/json"
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

	if v, ok := d.GetOk("custom_attributes"); ok {
		var attrs changetracking.ChangeTrackingRawCustomAttributesMap
		if err := json.Unmarshal([]byte(v.(string)), &attrs); err != nil {
			return eventInput, changetracking.ChangeTrackingDataHandlingRules{}, err
		}
		eventInput.CustomAttributes = attrs
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

	if kindList, ok := cfg["kind"].([]interface{}); ok && len(kindList) > 0 {
		kindCfg := kindList[0].(map[string]interface{})
		kind := &changetracking.ChangeTrackingCategoryAndTypeInput{}
		if c, ok := kindCfg["category"]; ok {
			kind.Category = c.(string)
		}
		if t, ok := kindCfg["type"]; ok {
			kind.Type = t.(string)
		}
		input.Kind = kind
	}

	if cfList, ok := cfg["category_fields"].([]interface{}); ok && len(cfList) > 0 {
		cfCfg := cfList[0].(map[string]interface{})
		categoryFields := &changetracking.ChangeTrackingCategoryFieldsInput{}

		if depList, ok := cfCfg["deployment"].([]interface{}); ok && len(depList) > 0 {
			depCfg := depList[0].(map[string]interface{})
			dep := &changetracking.ChangeTrackingDeploymentFieldsInput{}
			if v, ok := depCfg["changelog"]; ok {
				dep.Changelog = v.(string)
			}
			if v, ok := depCfg["commit"]; ok {
				dep.Commit = v.(string)
			}
			if v, ok := depCfg["deep_link"]; ok {
				dep.DeepLink = v.(string)
			}
			if v, ok := depCfg["version"]; ok {
				dep.Version = v.(string)
			}
			categoryFields.Deployment = dep
		}

		if ffList, ok := cfCfg["feature_flag"].([]interface{}); ok && len(ffList) > 0 {
			ffCfg := ffList[0].(map[string]interface{})
			ff := &changetracking.ChangeTrackingFeatureFlagFieldsInput{}
			if v, ok := ffCfg["feature_flag_id"]; ok {
				ff.FeatureFlagId = v.(string)
			}
			categoryFields.FeatureFlag = ff
		}

		input.CategoryFields = categoryFields
	}

	return input, true
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

// Ensure imported packages are used.
var _ = common.EntityGUID("")

// No flatten functions: this API is write-only.
