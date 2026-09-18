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

	if v, ok := d.GetOk("short_description"); ok {
		eventInput.ShortDescription = v.(string)
	}

	if v, ok := d.GetOk("group_id"); ok {
		eventInput.GroupId = v.(string)
	}

	if v, ok := d.GetOk("user"); ok {
		eventInput.User = v.(string)
	}

	if v, ok := d.GetOk("timestamp"); ok {
		eventInput.Timestamp = nrtime.EpochMilliseconds(time.UnixMilli(int64(v.(int))))
	}

	// Build entitySearch
	entitySearch := changetracking.ChangeTrackingEntitySearchInput{}
	hasEntitySearch := false

	if v, ok := d.GetOk("entity_guid"); ok {
		entitySearch.Query = "id = '" + v.(string) + "'"
		hasEntitySearch = true
	}

	if v, ok := d.GetOk("entity_search_query"); ok {
		entitySearch.Query = v.(string)
		hasEntitySearch = true
	}

	if hasEntitySearch {
		eventInput.EntitySearch = entitySearch
	}

	// Build categoryAndTypeData
	categoryAndTypeData := &changetracking.ChangeTrackingCategoryRelatedInput{}
	hasCategoryData := false

	if v, ok := d.GetOk("category_type"); ok {
		kind := &changetracking.ChangeTrackingCategoryAndTypeInput{}
		categoryType := changetracking.ChangeTrackingCategoryType(v.(string))
		// Parse category and type from the combined string
		// The format is CATEGORY__TYPE or CATEGORY__SUBTYPE
		// We need to split on the first double underscore
		catStr := string(categoryType)
		splitIdx := -1
		for i := 0; i < len(catStr)-1; i++ {
			if catStr[i] == '_' && catStr[i+1] == '_' {
				splitIdx = i
				break
			}
		}
		if splitIdx >= 0 {
			kind.Category = catStr[:splitIdx]
			kind.Type = catStr[splitIdx+2:]
		} else {
			kind.Category = catStr
			kind.Type = catStr
		}
		categoryAndTypeData.Kind = kind
		hasCategoryData = true
	} else if v, ok := d.GetOk("category"); ok {
		kind := &changetracking.ChangeTrackingCategoryAndTypeInput{}
		kind.Category = v.(string)
		if t, ok := d.GetOk("deployment_type"); ok {
			kind.Type = t.(string)
		}
		categoryAndTypeData.Kind = kind
		hasCategoryData = true
	}

	// Build categoryFields
	categoryFields := &changetracking.ChangeTrackingCategoryFieldsInput{}
	hasCategoryFields := false

	// Deployment fields
	deploymentFields := &changetracking.ChangeTrackingDeploymentFieldsInput{}
	hasDeploymentFields := false

	if v, ok := d.GetOk("version"); ok {
		deploymentFields.Version = v.(string)
		hasDeploymentFields = true
	}

	if v, ok := d.GetOk("changelog"); ok {
		deploymentFields.Changelog = v.(string)
		hasDeploymentFields = true
	}

	if v, ok := d.GetOk("commit"); ok {
		deploymentFields.Commit = v.(string)
		hasDeploymentFields = true
	}

	if v, ok := d.GetOk("deep_link"); ok {
		deploymentFields.DeepLink = v.(string)
		hasDeploymentFields = true
	}

	if hasDeploymentFields {
		categoryFields.Deployment = deploymentFields
		hasCategoryFields = true
	}

	// Feature flag fields
	if v, ok := d.GetOk("feature_flag_id"); ok {
		categoryFields.FeatureFlag = &changetracking.ChangeTrackingFeatureFlagFieldsInput{
			FeatureFlagId: v.(string),
		}
		hasCategoryFields = true
	}

	if hasCategoryFields {
		categoryAndTypeData.CategoryFields = categoryFields
		hasCategoryData = true
	}

	if hasCategoryData {
		eventInput.CategoryAndTypeData = categoryAndTypeData
	}

	// Build dataHandlingRules
	dataHandlingRules := changetracking.ChangeTrackingDataHandlingRules{}

	if v, ok := d.GetOk("validation_flags"); ok {
		flags := v.([]interface{})
		validationFlags := make([]changetracking.ChangeTrackingValidationFlag, len(flags))
		for i, f := range flags {
			validationFlags[i] = changetracking.ChangeTrackingValidationFlag(f.(string))
		}
		dataHandlingRules.ValidationFlags = validationFlags
	}

	return eventInput, dataHandlingRules, nil
}

var _ = common.EntityGUID("")

// No flatten functions: this API is write-only.
