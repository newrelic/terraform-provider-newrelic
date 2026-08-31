package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/changetracking"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
)

func expandChangeTrackingDeployment(d *schema.ResourceData) (*changetracking.ChangeTrackingDeploymentInput, error) {
	input := changetracking.ChangeTrackingDeploymentInput{}

	input.Version = d.Get("version").(string)
	input.EntityGUID = common.EntityGUID(d.Get("entity_guid").(string))

	if v, ok := d.GetOk("changelog"); ok {
		input.Changelog = v.(string)
	}

	if v, ok := d.GetOk("commit"); ok {
		input.Commit = v.(string)
	}

	if v, ok := d.GetOk("deep_link"); ok {
		input.DeepLink = v.(string)
	}

	if v, ok := d.GetOk("deployment_type"); ok {
		input.DeploymentType = changetracking.ChangeTrackingDeploymentType(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = v.(string)
	}

	if v, ok := d.GetOk("group_id"); ok {
		input.GroupId = v.(string)
	}

	if v, ok := d.GetOk("timestamp"); ok {
		input.Timestamp = nrtime.EpochMilliseconds(v.(int))
	}

	if v, ok := d.GetOk("user"); ok {
		input.User = v.(string)
	}

	return &input, nil
}

func expandChangeTrackingDataHandlingRules(d *schema.ResourceData) changetracking.ChangeTrackingDataHandlingRules {
	return changetracking.ChangeTrackingDataHandlingRules{
		ValidationFlags: []changetracking.ChangeTrackingValidationFlag{},
	}
}
func flattenChangeTrackingDeployment(result *changetracking.ChangeTrackingDeployment, d *schema.ResourceData) error {
	_ = d.Set("version", result.Version)
	_ = d.Set("entity_guid", string(result.EntityGUID))
	_ = d.Set("changelog", result.Changelog)
	_ = d.Set("commit", result.Commit)
	_ = d.Set("deep_link", result.DeepLink)
	_ = d.Set("deployment_type", string(result.DeploymentType))
	_ = d.Set("description", result.Description)
	_ = d.Set("group_id", result.GroupId)
	_ = d.Set("timestamp", int(result.Timestamp))
	_ = d.Set("user", result.User)
	_ = d.Set("deployment_id", result.DeploymentId)
	return nil
}