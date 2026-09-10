package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/changetracking"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
	"time"
)

func expandChangeTrackingEvent(d *schema.ResourceData) (changetracking.ChangeTrackingDataHandlingRules, changetracking.ChangeTrackingDeploymentInput, error) {
	dataHandlingRules := changetracking.ChangeTrackingDataHandlingRules{}

	deployment := changetracking.ChangeTrackingDeploymentInput{
		EntityGUID: common.EntityGUID(d.Get("entity_guid").(string)),
		Version:    d.Get("version").(string),
	}

	if v, ok := d.GetOk("changelog"); ok {
		deployment.Changelog = v.(string)
	}

	if v, ok := d.GetOk("commit"); ok {
		deployment.Commit = v.(string)
	}

	if v, ok := d.GetOk("deep_link"); ok {
		deployment.DeepLink = v.(string)
	}

	if v, ok := d.GetOk("deployment_type"); ok {
		deployment.DeploymentType = changetracking.ChangeTrackingDeploymentType(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		deployment.Description = v.(string)
	}

	if v, ok := d.GetOk("group_id"); ok {
		deployment.GroupId = v.(string)
	}

	if v, ok := d.GetOk("timestamp"); ok {
		deployment.Timestamp = nrtime.EpochMilliseconds(time.UnixMilli(int64(v.(int))))
	}

	if v, ok := d.GetOk("user"); ok {
		deployment.User = v.(string)
	}

	return dataHandlingRules, deployment, nil
}

// No flatten functions: this API is write-only.
