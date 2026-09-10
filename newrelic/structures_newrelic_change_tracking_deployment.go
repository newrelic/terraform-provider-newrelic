package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/changetracking"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
	"time"
)

func expandChangeTrackingDeployment(d *schema.ResourceData) (changetracking.ChangeTrackingDataHandlingRules, changetracking.ChangeTrackingDeploymentInput, error) {
	dataHandlingRules := changetracking.ChangeTrackingDataHandlingRules{}

	input := changetracking.ChangeTrackingDeploymentInput{
		EntityGUID: common.EntityGUID(d.Get("entity_guid").(string)),
		Version:    d.Get("version").(string),
	}

	if v, ok := d.GetOk("deployment_type"); ok {
		input.DeploymentType = changetracking.ChangeTrackingDeploymentType(v.(string))
	}

	if v, ok := d.GetOk("changelog"); ok {
		input.Changelog = v.(string)
	}

	if v, ok := d.GetOk("commit"); ok {
		input.Commit = v.(string)
	}

	if v, ok := d.GetOk("deep_link"); ok {
		input.DeepLink = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = v.(string)
	}

	if v, ok := d.GetOk("group_id"); ok {
		input.GroupId = v.(string)
	}

	if v, ok := d.GetOk("user"); ok {
		input.User = v.(string)
	}

	if v, ok := d.GetOk("timestamp"); ok {
		input.Timestamp = nrtime.EpochMilliseconds(time.UnixMilli(int64(v.(int))))
	}

	return dataHandlingRules, input, nil
}

// No flatten functions: this API is write-only.
