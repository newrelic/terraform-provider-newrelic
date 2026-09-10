package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/changetracking"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
	"time"
)

func expandChangeTrackingEvent(d *schema.ResourceData) (changetracking.ChangeTrackingCreateEventInput, changetracking.ChangeTrackingDataHandlingRules, error) {
	eventInput := changetracking.ChangeTrackingCreateEventInput{}

	deploymentFields := &changetracking.ChangeTrackingDeploymentFieldsInput{
		Version: d.Get("version").(string),
	}

	if v, ok := d.GetOk("changelog"); ok {
		deploymentFields.Changelog = v.(string)
	}

	if v, ok := d.GetOk("commit"); ok {
		deploymentFields.Commit = v.(string)
	}

	if v, ok := d.GetOk("deep_link"); ok {
		deploymentFields.DeepLink = v.(string)
	}

	categoryFields := &changetracking.ChangeTrackingCategoryFieldsInput{
		Deployment: deploymentFields,
	}

	kind := &changetracking.ChangeTrackingCategoryAndTypeInput{
		Category: "DEPLOYMENT",
		Type:     "BASIC",
	}

	if v, ok := d.GetOk("deployment_type"); ok {
		kind.Type = v.(string)
	}

	eventInput.CategoryAndTypeData = &changetracking.ChangeTrackingCategoryRelatedInput{
		CategoryFields: categoryFields,
		Kind:           kind,
	}

	if v, ok := d.GetOk("description"); ok {
		eventInput.Description = v.(string)
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

	entityGUID := d.Get("entity_guid").(string)
	eventInput.EntitySearch = changetracking.ChangeTrackingEntitySearchInput{
		Query: "id = '" + string(common.EntityGUID(entityGUID)) + "'",
	}

	dataHandlingRules := changetracking.ChangeTrackingDataHandlingRules{
		ValidationFlags: []changetracking.ChangeTrackingValidationFlag{},
	}

	return eventInput, dataHandlingRules, nil
}

// No flatten functions: this API is write-only.
